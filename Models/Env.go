package Models

import (
	"fmt"
	"os"
)

var envStruct *Env

type Env struct {
	OriginAllowed string
	Env string
	Port string
	TimeShiftGrpcServer string
	TimeShiftGrpcPort string
	MediaMetadataGrpcServer string
	MediaMetadataGrpcPort string
	MediaChunkMetadataServer string
	MediaChunkMetadataPort string
	RabbitUser string
	RabbitPassword string
	RabbitQueue string
	RabbitHost string
	RabbitPort string
	AawsStorageUrl string
}

func InitEnv()  {
	envStruct = &Env{
		OriginAllowed:  			os.Getenv("ORIGIN_ALLOWED"),
		Env: 			  			os.Getenv("ENV"),
		Port: 						os.Getenv("PORT"),
		TimeShiftGrpcServer:		os.Getenv("TIMESHIFT_GRPC_SERVER"),
		TimeShiftGrpcPort:			os.Getenv("TIMESHIFT_GRPC_PORT"),
		MediaMetadataGrpcServer:	os.Getenv("MEDIA_METADATA_GRPC_SERVER"),
		MediaMetadataGrpcPort:		os.Getenv("MEDIA_METADATA_GRPC_PORT"),
		MediaChunkMetadataServer:	os.Getenv("MEDIA_CHUNKS_METADATA_GRPC_SERVER"),
		MediaChunkMetadataPort:		os.Getenv("MEDIA_CHUNKS_METADATA_GRPC_PORT"),
		RabbitUser:       			os.Getenv("RABBIT_USER"),
		RabbitPassword:   			os.Getenv("RABBIT_PASSWORD"),
		RabbitQueue:      			os.Getenv("RABBIT_QUEUE"),
		RabbitHost:       			os.Getenv("RABBIT_HOST"),
		RabbitPort: 				os.Getenv("RABBIT_PORT"),
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}