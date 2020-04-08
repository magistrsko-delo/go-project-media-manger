package grpc_client

import (
	"fmt"
	"go-project-media-manger/Models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbMediaMediaMetadata "go-project-media-manger/proto/media-metadata"
	"log"
)

type MediaMetadataClient struct {
	Conn *grpc.ClientConn
	client pbMediaMediaMetadata.MediaMetadataClient
}

func (mediaMetadataClient *MediaMetadataClient) GetProjectMedias(projectId int32) (*pbMediaMediaMetadata.MediaMetadataResponseRepeated, error)  {

	response, err := mediaMetadataClient.client.GetProjectMediasMetadata(context.Background(), &pbMediaMediaMetadata.GetProjectMediasRequest{
		ProjectId:            projectId,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mediaMetadataClient *MediaMetadataClient) GetOneProjectMedia(projectId int32, mediaId int32)  (*pbMediaMediaMetadata.MediaMetadataResponse, error) {
	response, err := mediaMetadataClient.client.GetOneProjectMediasMetadata(context.Background(), &pbMediaMediaMetadata.GetOneProjectMedia{
		MediaId:              mediaId,
		ProjectId:            projectId,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mediaMetadataClient *MediaMetadataClient) CreateNewMedia(name string, projectId int32) (*pbMediaMediaMetadata.MediaMetadataResponse, error) {
	response, err := mediaMetadataClient.client.NewMediaMetadata(context.Background(),&pbMediaMediaMetadata.CreateNewMediaMetadataRequest{
		Name:                     name,
		SiteName:                 "PROJECT",
		Length:                   -1,
		Status:                   0,
		Thumbnail:                "",
		ProjectId:                projectId,
		AwsBucketWholeMedia:      "",
		AwsStorageNameWholeMedia: "",
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitMediaMetadataClient() *MediaMetadataClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING mediaMetadata client")
	conn, err := grpc.Dial(env.MediaMetadataGrpcServer + ":" + env.MediaMetadataGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION mediaMetadata client")

	client := pbMediaMediaMetadata.NewMediaMetadataClient(conn)
	return &MediaMetadataClient{
		Conn:    conn,
		client:  client,
	}
}