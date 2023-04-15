package main

import (
	"github.com/sirupsen/logrus"

	"github.com/sergripenko/chatRPC/internal/server"
)

func main() {

	app := server.NewApp()
	if err := app.Run(); err != nil {
		logrus.Fatal(err)
	}
}
