package DTO

type ResponseDTO struct {
	Status int32 `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
} 
