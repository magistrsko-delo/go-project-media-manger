package grpc_client

import (
	"fmt"
	"go-project-media-manger/Models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbMediaChunks "go-project-media-manger/proto/media-chunks-metadata"
	"log"
)

type MediaChunksClient struct {
	Conn *grpc.ClientConn
	client pbMediaChunks.MediaMetadataClient
}

func (mediaChunksClient *MediaChunksClient) LinkMediaWithChunks(mediaId int32, position int32, resolution string, chunkId int32) (*pbMediaChunks.LinkMediaChunkResponse, error)  {

	response, err := mediaChunksClient.client.LinkMediaWithChunk(context.Background(), &pbMediaChunks.LinkMediaWithChunkRequest{
		MediaId:              mediaId,
		Position:             position,
		Resolution:           resolution,
		ChunkId:              chunkId,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func InitMediaChunksClient() *MediaChunksClient {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING media chunks client")
	conn, err := grpc.Dial(env.MediaChunkMetadataServer + ":" + env.MediaChunkMetadataPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION media chunks client")

	client := pbMediaChunks.NewMediaMetadataClient(conn)
	return &MediaChunksClient{
		Conn:    conn,
		client:  client,
	}
}
