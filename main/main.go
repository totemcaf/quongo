package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"
  "strings"
  "encoding/json"
)

const VERSION = "1.0.0"

func main() {
  port, err := getIntEnv("QUONGO_PORT", 7070)

  if (err != nil) {
    log.Fatal("Invalid port number in QUONGO_PORT ", err)
    panic(err)
  }

  fmt.Println("Quongo ", VERSION, " running in port ", port)

  queueRep := NewQueueRepository()

  qCtrl := NewQueueController(queueRep)

  router := NewRouter(&qCtrl.routes)

  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), router))
}

func getIntEnv(name string, defValue int) (int, error) {
  v := strings.TrimSpace(os.Getenv(name))

  if (v == "") {
    return defValue, nil
  } else {
    return strconv.Atoi(v)
  }
}

func Ok(w http.ResponseWriter, v interface{}) {
  if err := json.NewEncoder(w).Encode(v); err != nil {
    panic(err)
  }
}