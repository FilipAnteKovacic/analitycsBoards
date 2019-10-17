package main

import (
	"log"
	"os"
)

var quitUpdate chan bool

func main() {

	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	restAddress := ":" + os.Getenv("REST")
	grpcAddress := ":" + os.Getenv("GRPC")
	uiAddress := ":" + os.Getenv("UI")
	certFile := "cert/server.crt"
	keyFile := "cert/server.pem"

	if os.Getenv("GRPC") != "" {

		// fire the gRPC server in a goroutine
		go func() {
			err := startGRPCServer(grpcAddress, certFile, keyFile)
			if err != nil {
				log.Fatalf("failed to start gRPC server: %s", err)
			}
		}()

	}

	if os.Getenv("REST") != "" && os.Getenv("GRPC") != "" {

		// fire the REST server in a goroutine
		go func() {
			err := startRESTServer(restAddress, grpcAddress, certFile)
			if err != nil {
				log.Fatalf("failed to start rest server: %s", err)
			}
		}()

	}
	// connect wsConnections

	if os.Getenv("UI") != "" {

		go func() {
			wsConection()
		}()

		// fire the UI server in a goroutine
		go func() {
			err := startUIServer(uiAddress)
			if err != nil {
				log.Fatalf("failed to start ws server: %s", err)
			}
		}()

	}

	quitUpdate = make(chan bool)

	go func() {
		chartsDataUpdate()
	}()

	// infinite loop
	log.Printf("Entering infinite loop")
	select {}

}

func dataUpdateRestart() {

	go func() {

		// Quit goroutine
		quitUpdate <- true

		chartsDataUpdate()
		return

	}()

	return
}
