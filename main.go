package main

import (
	"github.com/deemanrip/promotions/controller"
	"github.com/deemanrip/promotions/repository"
	"github.com/deemanrip/promotions/service"
)

func main() {
	minioClient := service.CreateClient()
	repository.Init()
	go service.ListenNotifications(minioClient)
	controller.GinInit()
}
