package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/totemcaf/quongo/main/infraestructure/mongodb"
	"github.com/totemcaf/quongo/main/interface"
	"github.com/totemcaf/quongo/main/usecase"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const VERSION = "1.0.0"

func main() {
	port, err := getIntEnv("QUONGO_PORT", 7070)
	user := getStrEnv("MONGO_USER", "quongo")
	pwd := getStrEnv("MONGO_PWD", "quongosecret")
	mongoUrls := getStrEnv("MONGO_URL", "localhost:27017")
	mongoDb := getStrEnv("MONGO_DB", "quongo")

	if err != nil {
		log.Fatal("Invalid port number in MONGO_PORT ", err)
		panic(err)
	}

	fmt.Println("Quongo ", VERSION, " running in port ", port)

	var aLogger *log.Logger
	aLogger = log.New(os.Stderr, "", log.LstdFlags)

	// Database
	db, err := mongodb.NewDatabase(mongoUrls, user, pwd, mongoDb, aLogger)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Repositories
	queueRep := mongodb.NewQueueRepository(db)
	msgRep := mongodb.NewMongodbRepositoryProvider(db)

	// Interactors
	queueInt := usecase.NewQueueInteractor(queueRep)
	msgInt := usecase.NewMessageInteractor(msgRep)

	// Views
	queueView := _interface.NewQueueView(queueInt)
	msgView := _interface.NewMessageView(msgInt, queueInt)

	routes := append(queueView.Routes, msgView.Routes...)

	stack := []rest.Middleware{
		&rest.AccessLogApacheMiddleware{},
		//      Format: rest.CombinedLogFormat,
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
		&rest.GzipMiddleware{},
	}

	api := rest.NewApi()

	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)

	api.Use(stack...)

	// Add system endpoints
	routes2 := append(routes,
		rest.Get("/system/stats", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}), // TODO Joining both
		rest.Get("/system/status", SysStatus),
	)

	router, err := rest.MakeRouter(routes2...)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getIntEnv(name string, defValue int) (int, error) {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue, nil
	} else {
		return strconv.Atoi(v)
	}
}

func getStrEnv(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		return v
	}
}
