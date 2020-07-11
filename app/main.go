package app

import (
	"log"
	"os"

	"github.com/totemcaf/quongo/app/utils"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	a := NewApp(utils.GetStrEnv("REPO_TYPE", "memory"), utils.ProductionClock(), logger)

	a.Start(getPort())
}

func getPort() int {
	port, err := utils.GetIntEnv("QUONGO_PORT", 7070)

	if err != nil {
		log.Fatal("Invalid port number in QUONGO_PORT ", err)
		panic(err)
	}

	return port
}
