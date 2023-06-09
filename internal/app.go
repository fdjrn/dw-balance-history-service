package internal

import (
	"errors"
	"fmt"
	"github.com/fdjrn/dw-balance-history-service/configs"
	"github.com/fdjrn/dw-balance-history-service/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type MongoDB struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Application struct {
	Database MongoDB
}

func Initialize() error {
	// Initialize
	err := configs.Initialize()
	if err != nil {
		return errors.New(fmt.Sprintf("error on config initialization: %s", err.Error()))
	}

	// DB Connection
	if err = db.Mongo.Connect(); err != nil {
		return errors.New(fmt.Sprintf("error on mongodb connection: %s", err.Error()))
	}

	SetupCloseHandler()
	return nil
}

func ExitGracefully() {
	// close mongodb connection
	if err := db.Mongo.Disconnect(); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("[CONN] db >> connection successfully disconnected")

	// close kafka connection
	//_ = kafka.Producer.Close()
	//log.Println("[CONN] kafka >> producer successfully closed")
}

// SetupCloseHandler :
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal,... Good Bye...")
		ExitGracefully()
		os.Exit(0)
	}()
}
