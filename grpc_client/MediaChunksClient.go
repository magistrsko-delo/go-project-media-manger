package grpc_client

import (
	"fmt"
	"go-project-media-manger/Models"
	"google.golang.org/grpc"
	pbMediaChunks "go-project-media-manger/proto/media-chunks-metadata"
	"log"

	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	ot "github.com/opentracing/opentracing-go"
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
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_opentracing.StreamClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer()))))
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_opentracing.UnaryClientInterceptor(
			grpc_opentracing.WithTracer(ot.GlobalTracer()))))
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())

	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING media chunks client")
	conn, err := grpc.DialContext(context.Background(), env.MediaChunkMetadataServer + ":" + env.MediaChunkMetadataPort, opts...) // grpc.Dial(env.MediaChunkMetadataServer + ":" + env.MediaChunkMetadataPort, grpc.WithInsecure(), grpc.WithBlock())
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
