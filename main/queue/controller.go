package queue

import (
  "github.com/ant0ine/go-json-rest/rest"
  "net/http"
)

type Controller struct {
  repository *Repository
  Routes     [] *rest.Route
}

func NewQueueController(queueRep *Repository) *Controller {

  ctrl := Controller{repository: queueRep}

  routes := [] *rest.Route{
    rest.Get("/v1/queue", ctrl.QueueAll),
    rest.Get("/v1/queue/#queueId", ctrl.QueueGet),
    // Queue services),
    rest.Post("/v1/queue", ctrl.QueueAdd, ),
    rest.Put("/v1/queue/#queueId", ctrl.QueueUpd),
  }

  ctrl.Routes = routes

  return &ctrl
}

// GET  	/v1/queue                                        @controllers.Controller.findAll()
func (c *Controller) QueueAll(w rest.ResponseWriter, r *rest.Request) {
  q, e := c.repository.findAll(0, 100)

  if e == nil {
    w.WriteJson(q)
  } else {
    rest.Error(w, e.Error(), http.StatusInternalServerError)
  }
}

// GET  	/v1/queue/:queueId                               @controllers.Controller.find(queueId: String)
func (c *Controller) QueueGet(w rest.ResponseWriter, r *rest.Request) {
  qId := r.PathParam("queueId")

  q, e := c.repository.findById(qId)

  if e == nil {
    w.WriteJson(q)
  } else {
    rest.Error(w, e.Error(), http.StatusInternalServerError)
  }
}

// POST 	/v1/queue                                        @controllers.Controller.create(queueName: String)
func (c *Controller) QueueAdd(w rest.ResponseWriter, r *rest.Request) {
  qId := r.PathParam("queueId")

  pQ, e := NewQueue(qId)

  if e == nil {
    added, e2 := c.repository.add(pQ)

    if e2 == nil {
      w.WriteJson(added)
    } else {
      rest.Error(w, e.Error(), http.StatusBadRequest)
    }

  } else {
    rest.Error(w, e.Error(), http.StatusBadRequest)
  }
}

// #PUT  /v1/queue/:queueId                                 @controllers.Controller.Update(queueId: String)
func (c *Controller) QueueUpd(w rest.ResponseWriter, r *rest.Request) {
  queue := Queue{}
  err := r.DecodeJsonPayload(&queue)

  if err != nil {
    rest.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  queue.Name = r.PathParam("queueId")

  updated, errU := c.repository.Update(&queue)

  if errU == nil {
    w.WriteJson(updated)
  } else {
    rest.Error(w, errU.Error(), http.StatusBadRequest)
  }
}

// #DELETE /v1/queue/:queueId                               @controllers.Controller.delete(queueId: String)
func (c *Controller) Delete(w rest.ResponseWriter, r *rest.Request) {

}

