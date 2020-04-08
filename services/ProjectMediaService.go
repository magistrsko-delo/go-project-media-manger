package services

import (
	"go-project-media-manger/Models/DTO"
	"go-project-media-manger/grpc_client"
)

type ProjectMediaService struct {
	MediaMetadataGrpcClient *grpc_client.MediaMetadataClient
}

func (projectMediaService *ProjectMediaService) GetProjectMetadata(projectId int32) (*DTO.ResponseDTO, error)  {
	rsp, err := projectMediaService.MediaMetadataGrpcClient.GetProjectMedias(projectId)

	if err != nil {
		return nil, err
	}

	return &DTO.ResponseDTO{
		Status:  200,
		Message: "",
		Data:    rsp.GetData(),
	}, nil
}

func (projectMediaService *ProjectMediaService) GetOneProjectMedia(projectId int32, mediaId int32) (*DTO.ResponseDTO, error)  {
	rsp, err := projectMediaService.MediaMetadataGrpcClient.GetOneProjectMedia(projectId, mediaId);
	if err != nil {
		return nil, err
	}

	return &DTO.ResponseDTO{
		Status:  200,
		Message: "",
		Data:    rsp,
	}, nil
}
