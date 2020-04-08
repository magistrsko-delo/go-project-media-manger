package router

import (
	"github.com/gorilla/mux"
	"go-project-media-manger/controllers"
	"go-project-media-manger/grpc_client"
	"go-project-media-manger/services"
)

type ProjectMediaRouter struct {
	Router *mux.Router
}

func (projectMediaRouter *ProjectMediaRouter) RegisterHandlers()  {
	projectMediaController := &controllers.ProjectMediaController{
		ProjectMediaService: &services.ProjectMediaService{
			MediaMetadataGrpcClient: grpc_client.InitMediaMetadataClient(),
		},
		MediaCutService: &services.MediaCutService{
			TimeShiftClient:grpc_client.InitTimeShiftClient(),
		},
	}
	(*projectMediaRouter).Router.HandleFunc("/project/{projectId}/media", projectMediaController.GetProjectMedia).Methods("GET")
	(*projectMediaRouter).Router.HandleFunc("/project/{projectId}/media/{mediaId}", projectMediaController.GetOneProjectMedia).Methods("GET")

	(*projectMediaRouter).Router.HandleFunc("/project/{projectId}/media/{mediaId}/cut", projectMediaController.CutMedia).Methods("POST")
}
