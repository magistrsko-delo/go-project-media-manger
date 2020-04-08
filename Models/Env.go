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
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}