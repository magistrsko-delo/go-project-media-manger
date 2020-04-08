package services

import (
	"errors"
	"fmt"
	"go-project-media-manger/Models"
	"go-project-media-manger/Models/DTO"
	"go-project-media-manger/Models/Input"
	"go-project-media-manger/grpc_client"
)

type MediaCutService struct {
	TimeShiftClient *grpc_client.TimeShiftClient
}

func (mediaCutService *MediaCutService) CutMedia(mediaId int32, projectId int32, inputCut *Input.InputCut) (*DTO.ResponseDTO, error)  {

	timeShiftChunkMediaRsp, err := mediaCutService.TimeShiftClient.GetMediaChunkInfo(mediaId)
	if err != nil {
		return nil, err
	}

	if ( len(timeShiftChunkMediaRsp.GetData()[0].GetChunks()) == 0 ) {
		return nil,  errors.New("no media chunks")
	}

	mediaChunks1080p := timeShiftChunkMediaRsp.GetData()[0].GetChunks()

	mediaLength := 0.0
	startMediaSearch := true
	endMediaSearch := false
	rabbitMQMediaCutQueueData := [] *Models.MediaCutDataQueue {}
	position := 0
	indexedChunks := []int32{}

	for i := 0; i < len(mediaChunks1080p); i++ {
		currMediaLength := mediaLength + mediaChunks1080p[i].GetLength()
		if !startMediaSearch && endMediaSearch && currMediaLength < inputCut.To {
			// take chunk and index..)
			indexedChunks = append(indexedChunks, mediaChunks1080p[i].GetChunkId())
			position++
		}

		if startMediaSearch && currMediaLength > inputCut.From && currMediaLength > inputCut.To {
			if mediaLength == inputCut.From && currMediaLength == inputCut.To {
				// take whole chunk for media
				indexedChunks = append(indexedChunks, mediaChunks1080p[i].GetChunkId())
			} else if mediaLength == inputCut.From  {
				// take from chunk 0.0 til to
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					0.0,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position))

			} else if currMediaLength == inputCut.To {
				// take from to the end
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength(),
					mediaChunks1080p[i].GetChunkId(),
					int32(position))
			} else {
				// take from to.
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position))
			}
			position++
			break
		} else if startMediaSearch && currMediaLength > inputCut.From {
			if mediaLength == inputCut.From {
				// take whole chunk for indexing
				indexedChunks = append(indexedChunks, mediaChunks1080p[i].GetChunkId())
			} else {
				// take chunk (from - mediaLenght) to end
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					inputCut.From - mediaLength,
					mediaChunks1080p[i].GetLength(),
					mediaChunks1080p[i].GetChunkId(),
					int32(position))
			}

			position++
			startMediaSearch = false
			endMediaSearch = true

		} else if endMediaSearch && currMediaLength >= inputCut.To {
			// take chunks from beggining to   ( chunkLenght - ( currLenght - toLenght) )

			if currMediaLength - inputCut.To == 0 {
				indexedChunks = append(indexedChunks, mediaChunks1080p[i].GetChunkId())
			} else {
				rabbitMQMediaCutQueueData = mediaCutService.addRabbitMqDataWorker(
					rabbitMQMediaCutQueueData,
					0.0,
					mediaChunks1080p[i].GetLength() - ( currMediaLength - inputCut.To),
					mediaChunks1080p[i].GetChunkId(),
					int32(position))
			}
			position++
			break
		}
		mediaLength = currMediaLength
	}

	fmt.Println("MEDIA CUT CHUNKS: ", rabbitMQMediaCutQueueData)
	for i := 0; i < len(rabbitMQMediaCutQueueData); i++ {
		fmt.Println(rabbitMQMediaCutQueueData[i])
	}

	fmt.Println("Indexed chunks: ", indexedChunks)

	return &DTO.ResponseDTO{
		Status:  200,
		Message: "New project media created",
		Data:    nil,
	}, nil
}

func (mediaCutService *MediaCutService) addRabbitMqDataWorker(dataArray [] *Models.MediaCutDataQueue, from float64, to float64, chunkId int32, position int32) [] *Models.MediaCutDataQueue  {
	mediaCutDataQueue := &Models.MediaCutDataQueue{
		ChunkId:  chunkId,
		From:     from,
		To:       to,
		Position: position,
	}
	return append(dataArray, mediaCutDataQueue)
}