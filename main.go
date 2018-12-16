package main

import (
	"os"

	"github.com/owtotwo/Service-Computing-On-Cloud-Course/server"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	serverInstance := server.NewServer()
	serverInstance.Run(":" + port)
}
