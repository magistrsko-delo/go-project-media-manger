package grpc_client

import (
	"fmt"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	ot "github.com/opentracing/opentracing-go"
	"go-project-media-manger/Models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbTimeshift "go-project-media-manger/proto/timeshift-service"
	"log"
)

type TimeShiftClient struct {
	Conn *grpc.ClientConn
	client pbTimeshift.TimeshiftClient
}

func (timeShiftClient *TimeShiftClient) GetMediaChunkInfo(mediaId int32) (*pbTimeshift.TimeShitResponse, error)  {
	response, err := timeShiftClient.client.GetMediaChunkInformation(context.Background(), &pbTimeshift.TimeShiftRequest{
		MediaId: mediaId,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitTimeShiftClient() *TimeShiftClient  {

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
	fmt.Println("CONNECTING timeshift client")
	conn, err := grpc.DialContext(context.Background(), env.TimeShiftGrpcServer + ":" + env.TimeShiftGrpcPort, opts...) //  grpc.Dial(env.TimeShiftGrpcServer + ":" + env.TimeShiftGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION timeshift client")

	client := pbTimeshift.NewTimeshiftClient(conn)
	return &TimeShiftClient{
		Conn:    conn,
		client:  client,
	}
}