package main

import (
	"github.com/fdjrn/dw-balance-history-service/configs"
	"github.com/fdjrn/dw-balance-history-service/internal"
	"github.com/fdjrn/dw-balance-history-service/internal/kafka"
	"github.com/fdjrn/dw-balance-history-service/internal/routes"
	"log"
	"strings"
	"sync"
)

func main() {
	var err error
	internal.SetupCloseHandler()

	defer internal.ExitGracefully()

	wg := &sync.WaitGroup{}

	// Config Initialization
	if err = internal.Initialize(); err != nil {
		log.Fatalln(err)
	}

	kafka.StartConsumer()

	// Initialize Rest API
	wg.Add(1)
	go func() {
		err = routes.Initialize()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("[INIT] routes >> initialized")
		wg.Done()
	}()
	log.Println("[INIT] Rest API start at port: ", configs.MainConfig.APIServer.Port)
	log.Println(strings.Repeat("_", 40))

	wg.Wait()

}
