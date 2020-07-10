package main

import (
	"log"
	logger "log"
	"os"

	"github.com/totemcaf/quongo/main/delivery"
	"github.com/totemcaf/quongo/main/infrastructure/memory"
	"github.com/totemcaf/quongo/main/model"
	"github.com/totemcaf/quongo/main/usecase"
	"github.com/totemcaf/quongo/main/utils"
)

func main() {
	var log *logger.Logger
	log = logger.New(os.Stderr, "", logger.LstdFlags)

	//	queueRep, msgRep := createRepos(log)

	// Interactors
	queueInt := usecase.NewQueueInteractor()
	msgInt := usecase.NewMessageInteractor()
	statusInt := usecase.NewStatusInteractor()

	// Views
	queueView := delivery.NewQueueView(queueInt)
	msgView := delivery.NewMessageView(msgInt, queueInt)
	monitorView := delivery.NewMonitorView(statusInt)
	systemView := delivery.NewSystemView(Version)

	port := getPort()

	server := delivery.NewServer(port)

	server.AddSystem(monitorView)
	server.AddSystem(systemView)
	server.Add(queueView)
	server.Add(msgView)

	log.Printf("Quongo %v running in port %v", Version, port)

	log.Fatal(server.Start())
}

func getPort() int {
	port, err := utils.GetIntEnv("QUONGO_PORT", 7070)

	if err != nil {
		log.Fatal("Invalid port number in QUONGO_PORT ", err)
		panic(err)
	}

	return port
}

func createRepos(log *logger.Logger) (model.QueueRepository, model.MessageRepositoryProvider) {
	switch repoType := utils.GetStrEnv("REPO_TYPE", "memory"); repoType {
	case "memory":
		return memory.NewQueueRepository(), memory.NewQueueRepository()

	default:
		log.Fatalf("Unknown repository type '%s'", repoType)
		return nil, nil
	}
}
