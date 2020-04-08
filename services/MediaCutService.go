package services

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"go-project-media-manger/Models"
	"go-project-media-manger/Models/DTO"
	"go-project-media-manger/Models/Input"
	"go-project-media-manger/grpc_client"
	"math/rand"
	"strconv"
)

type MediaCutService struct {
	TimeShiftClient *grpc_client.TimeShiftClient
	MediaMetadataClient *grpc_client.MediaMetadataClient
	MediaChunksClient *grpc_client.MediaChunksClient
	RabbitMQ *RabbitMQ
}

func (mediaCutService *MediaCutService) CutMedia(mediaId int32, projectId int32, inputCut *Input.InputCut) (*DTO.ResponseDTO, error)  {

	timeShiftChunkMediaRsp, err := mediaCutService.TimeShiftClient.GetMediaChunkInfo(mediaId)
	if err != nil {
		return nil, err
	}

	if len(timeShiftChunkMediaRsp.GetData()[0].GetChunks()) == 0 {
		return nil,  errors.New("no media chunks")
	}

	newMediaRsp, err := mediaCutService.MediaMetadataClient.CreateNewMedia(
		timeShiftChunkMediaRsp.GetName() + "_PROJECT_" + strconv.Itoa(int(projectId)) + "_" + strconv.Itoa(rand.Intn(1000000000000)),
		projectId)

	if err != nil  {
		return nil, err
	}
	mediaChunks1080p := timeShiftChunkMediaRsp.GetData()[0].GetChunks()
	resolution := timeShiftChunkMediaRsp.GetData()[0].GetResolution()
	mediaLength := 0.0
	startMediaSearch := true
	endMediaSearch := false
	rabbitMQMediaCutQueueData := [] *Models.MediaCutDataQueue {}
	position := 0

	for i := 0; i < len(mediaChunks1080p); i++ {
		currMediaLength := mediaLength + mediaChunks1080p[i].GetLength()
		if !startMediaSearch && endMediaSearch && currMediaLength < inputCut.To {
			// take chunk and index..)
			_, err := mediaCutService.MediaChunksClient.LinkMediaWithChunks(newMediaRsp.GetMediaId(), int32(position), resolution, mediaChunks1080p[i].GetChunkId())
			if err != nil  {
				return nil, err
			}
			position++
		}

		if startMediaSearch && currMediaLength > inputCut.From && currMediaLength > inputCut.To {
			if mediaLength == inputCut.From && currMediaLength == inputCut.To {
				// take whole chunk for media
				_, err := mediaCutService.MediaChunksClient.LinkMediaWithChunks(newMediaRsp.GetMediaId(), int32(position), resolution, mediaChunks1080p[i].GetChunkId())
				if err != nil  {
					return nil, err
				}
			} else if mediaLength == inputCut.From  {
				// take from chunk 0.0 til to
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					0.0,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position),
					resolution,
					newMediaRsp.MediaId,)

			} else if currMediaLength == inputCut.To {
				// take from to the end
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength(),
					mediaChunks1080p[i].GetChunkId(),
					int32(position),
					resolution,
					newMediaRsp.MediaId,)
			} else {
				// take from to.
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position),
					resolution,
					newMediaRsp.MediaId,
					)
			}
			position++
			break
		} else if startMediaSearch && currMediaLength > inputCut.From {
			if mediaLength == inputCut.From {
				// take whole chunk for indexing
				_, err := mediaCutService.MediaChunksClient.LinkMediaWithChunks(newMediaRsp.GetMediaId(), int32(position), resolution, mediaChunks1080p[i].GetChunkId())
				if err != nil  {
					return nil, err
				}
			} else {
				// take chunk (from - mediaLenght) to end
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength(),
					mediaChunks1080p[i].GetChunkId(),
					int32(position),
					resolution,
					newMediaRsp.MediaId,
					)
			}

			position++
			startMediaSearch = false
			endMediaSearch = true

		} else if endMediaSearch && currMediaLength >= inputCut.To {
			// take chunks from beggining to   ( chunkLenght - ( currLenght - toLenght) )

			if currMediaLength - inputCut.To == 0 {
				_, err := mediaCutService.MediaChunksClient.LinkMediaWithChunks(newMediaRsp.GetMediaId(), int32(position), resolution, mediaChunks1080p[i].GetChunkId())
				if err != nil  {
					return nil, err
				}
			} else {
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					0.0,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position),
					resolution,
					newMediaRsp.MediaId,
					)
			}
			position++
			break
		}
		mediaLength = currMediaLength
	}

	err = mediaCutService.publishMessageOnQueue(rabbitMQMediaCutQueueData)

	if err != nil {
		return nil, err
	}

	newtimeShiftChunkMediaRsp, err := mediaCutService.TimeShiftClient.GetMediaChunkInfo(newMediaRsp.GetMediaId())
	if err != nil {
		return nil, err
	}

	return &DTO.ResponseDTO{
		Status:  200,
		Message: "New project media created",
		Data:    newtimeShiftChunkMediaRsp,
	}, nil
}

func (mediaCutService *MediaCutService) addRabbitMqDataWorker(
	dataArray [] *Models.MediaCutDataQueue,
	from float64,
	to float64,
	chunkId int32,
	position int32,
	resolution string,
	mediaId int32) [] *Models.MediaCutDataQueue {
	mediaCutDataQueue := &Models.MediaCutDataQueue{
		ChunkId:  	chunkId,
		From:     	from,
		To:       	to,
		Position: 	position,
		Resolution: resolution,
		MediaId:	mediaId,
	}
	return append(dataArray, mediaCutDataQueue)
}

func (mediaCutService *MediaCutService) publishMessageOnQueue(rabbitMQMediaCutQueueData [] *Models.MediaCutDataQueue) error {
	dataForQueue, err := json.Marshal(rabbitMQMediaCutQueueData)
	if err != nil {
		return err
	}
	err = mediaCutService.RabbitMQ.Ch.Publish(
		"",
		mediaCutService.RabbitMQ.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			Body:            dataForQueue,
		})
	if err != nil {
		return err
	}
	return nil
}