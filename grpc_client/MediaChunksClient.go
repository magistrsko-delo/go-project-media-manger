package grpc_client

import (
	"fmt"
	"go-project-media-manger/Models"
	"google.golang.org/grpc"
	pbMediaChunks "go-project-media-manger/proto/media-chunks-metadata"
	"log"
)

type MediaChunksClient struct {
	Conn *grpc.ClientConn
	client pbMediaChunks.MediaMetadataClient
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
