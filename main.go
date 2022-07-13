package main

import (
	"app/app/model"
	"app/server"
)

func main() {
	model.Setup()
	server.SetupAndListen()
}
