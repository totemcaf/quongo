package app

import (
	"log"
	"net/http"

	"github.com/totemcaf/quongo/app/delivery"
	"github.com/totemcaf/quongo/app/infrastructure/memory"
	"github.com/totemcaf/quongo/app/model"
	"github.com/totemcaf/quongo/app/usecase"
	"github.com/totemcaf/quongo/app/utils"
)

type App struct {
	logger   *log.Logger
	repoType string
	clock    utils.Clock
	server   *delivery.Server
}

func NewApp(repoType string, clock utils.Clock, logger *log.Logger) *App {
	queueRep, _ := createRepos(repoType, clock, logger)

	// Interactors
	queueInt := usecase.NewQueueInteractor(queueRep, clock)
	msgInt := usecase.NewMessageInteractor(queueInt)
	statusInt := usecase.NewStatusInteractor()

	// Views
	queueView := delivery.NewQueueView(queueInt)
	msgView := delivery.NewMessageView(msgInt, queueInt, clock)
	monitorView := delivery.NewMonitorView(statusInt)
	systemView := delivery.NewSystemView(Version)

	server := delivery.NewServer()

	server.AddSystem(monitorView)
	server.AddSystem(systemView)
	server.Add(queueView)
	server.Add(msgView)

	return &App{logger: logger, repoType: repoType, server: server, clock: clock}
}

func (a *App) Handler() http.Handler {
	return a.server.MakeHandler()
}

func (a *App) Start(port int) {
	a.logger.Printf("Quongo %v running in port %v", Version, port)
	a.logger.Fatal(a.server.Start(port))
}

func createRepos(repoType string, clock utils.Clock, logger *log.Logger) (model.QueueRepository, model.MessageRepositoryProvider) {
	switch repoType {
	case "memory":
		repo := memory.NewQueueRepository(clock)
		return repo, repo

	default:
		logger.Fatalf("Unknown repository type '%s'", repoType)
		return nil, nil
	}
}
