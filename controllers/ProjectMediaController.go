package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-project-media-manger/Models/Input"
	"go-project-media-manger/services"
	"net/http"
	"strconv"
)

type ProjectMediaController struct {
	ProjectMediaService *services.ProjectMediaService
	MediaCutService *services.MediaCutService
}

func (projectMediaController *ProjectMediaController) GetProjectMedia(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	projectId, err := strconv.Atoi(params["projectId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsp, err := projectMediaController.ProjectMediaService.GetProjectMetadata(int32(projectId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rsp)
}

func (projectMediaController *ProjectMediaController) GetOneProjectMedia(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	projectId, err := strconv.Atoi(params["projectId"])
	mediaId, err := strconv.Atoi(params["mediaId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsp, err := projectMediaController.ProjectMediaService.GetOneProjectMedia(int32(projectId), int32(mediaId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rsp)
}

func (projectMediaController *ProjectMediaController) CutMedia(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	projectId, err := strconv.Atoi(params["projectId"])
	mediaId, err := strconv.Atoi(params["mediaId"])

	inputCut := &Input.InputCut{}
	err = json.NewDecoder(r.Body).Decode(inputCut)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsp, err := projectMediaController.MediaCutService.CutMedia(int32(mediaId), int32(projectId), inputCut)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rsp)
}