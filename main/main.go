package main

import (
  "github.com/totemcaf/quongo/main/queue"
  "github.com/totemcaf/quongo/main/message"
  "github.com/ant0ine/go-json-rest/rest"
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"
  "strings"
  "gopkg.in/mgo.v2"
)

const VERSION = "1.0.0"

func main() {
  port, err := getIntEnv("QUONGO_PORT", 7070)
  mongoUrls := getStrEnv("MONGO_URL", "localhost:27017")
  mongoDb := getStrEnv("MONGO_DB", "quongo")

  if err != nil {
    log.Fatal("Invalid port number in QUONGO_PORT ", err)
    panic(err)
  }

  fmt.Println("Quongo ", VERSION, " running in port ", port)

  // Mongo
  session, err := mgo.Dial(mongoUrls)

  if err != nil {
    panic(err)
  }

  var aLogger *log.Logger
  aLogger = log.New(os.Stderr, "", log.LstdFlags)
  mgo.SetLogger(aLogger)

  mgo.SetDebug(true)

  defer session.Close()

  // Repositories
  msgRep := message.NewMsgRepository(session, mongoDb)
  queueRep := queue.NewQueueRepository(session, mongoDb, msgRep)

  // Controllers
  queueCtr := queue.NewQueueController(queueRep)
  msgCtrl := message.NewMsgController(msgRep)

  routes := append(queueCtr.Routes, msgCtrl.Routes...)

  /*
  // System
  rest.Get("/system/status", SysStatus)
)
   */

  stack := [] rest.Middleware{
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
    }),       // TODO Joing both
    rest.Get("/system/status", SysStatus),
  )

  router, err := rest.MakeRouter(routes2...)

  if err != nil {
    log.Fatal(err)
  }
  api.SetApp(router)

  http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
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
