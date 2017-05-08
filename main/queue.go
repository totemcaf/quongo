package main

import (
  "time"
  "net/http"
  "github.com/gorilla/mux"
)

type Queue struct {
  Name      string        `json:"_id"`
  Created   time.Time     `json:"created"`
  VisWnd    time.Duration `json:"visibilityWindow"`
}

type QueueWithStats struct {
  *Queue
  Stats     *QueueStats  `json:"stats"`
}

type QueueStats struct {
  Total     int16       `json:"total"`
  Hidden    int16       `json:"hidden"`
  InProcess int16       `json:"inProcess"`
}

type QueueController struct {
  queueRep    *QueueRepository
  routes      Routes
}

func NewQueueController(queueRep *QueueRepository) *QueueController {

  qc := QueueController{queueRep: queueRep}

  var routes = Routes{
    // Queue services
    Route{
      "List Queues",
      "GET",
      "/v1/queue",
      qc.QueueAll,
    },
    Route{
      "Find Queue",
      "GET",
      "/v1/queue/{queueId}",
      qc.QueueGet,
    },
    Route{
      "Find Queue",
      "POST",
      "/v1/queue",
      qc.QueueAdd,
    },
    Route{
      "Update Queue",
      "PUT",
      "/v1/queue/{queueId}",
      qc.QueueUpd,
    },
  }

  qc.routes = routes

  return &qc
}

// GET  	/v1/queue                                        @controllers.QueueController.findAll()
func (c *QueueController) QueueAll(w http.ResponseWriter, r *http.Request) {

}
// GET  	/v1/queue/:queueId                               @controllers.QueueController.find(queueId: String)
func (c *QueueController) QueueGet(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  qId := vars["queueId"]

  q, _ := c.queueRep.findById(qId)

  Ok(w, q)
}

// POST 	/v1/queue                                        @controllers.QueueController.create(queueName: String)
func (c *QueueController) QueueAdd(w http.ResponseWriter, r *http.Request) {

}
// #PUT  /v1/queue/:queueId                                 @controllers.QueueController.update(queueId: String)
func (c *QueueController) QueueUpd(w http.ResponseWriter, r *http.Request) {

}

// #DELETE /v1/queue/:queueId                               @controllers.QueueController.delete(queueId: String)
func (c *QueueController) QueueDelete(w http.ResponseWriter, r *http.Request) {

}

