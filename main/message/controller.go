package message

import (
  "github.com/ant0ine/go-json-rest/rest"
)

const ISO_DATE_PATTERN = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"

type Controller struct {
  Repository    *Repository
  Routes        [] *rest.Route
 }

func NewMsgController(repository *Repository) *Controller {
  ctrl := Controller{repository, nil }

  routes := [] *rest.Route {
    rest.Get("/v1/queue/{queueId}/message", ctrl.All),
    rest.Get("/v1/queue/{queueId}/next", ctrl.Pop),
    rest.Get("/v1/queue/{queueId}/message/{mid}", ctrl.Get),
    rest.Put("/v1/queue/{queueId}/message/{mid}", ctrl.PushWithId),
    rest.Post("/v1/queue/{queueId}/message", ctrl.Push),
    rest.Put("/v1/queue/{queueId}/message/{mid}/ack/{ack}", ctrl.KeepAlive),
    rest.Delete("/v1/queue/{queueId}/message/{mid}/ack/{ack}", ctrl.Ack),
  }

  ctrl.Routes = routes

  return &ctrl
}

// GET     /v1/queue/:queueId/message                       @controllers.MessageController.findAll(queueId: String)
func (c *Controller) All(w rest.ResponseWriter, r *rest.Request) {

}

// GET     /v1/queue/:queueId/next                          @controllers.MessageController.pop(queueId: String, holder: Option[String] ?= None, window: Option[Int] ?= None)
func (c *Controller) Pop(w rest.ResponseWriter, r *rest.Request) {
  w.WriteJson("Siguiente mensaje para:" + r.PathParam("queueId"))
}

// GET     /v1/queue/:queueId/message/:mid                  @controllers.MessageController.findMsg(queueId: String, mid: String)
func (c *Controller) Get(w rest.ResponseWriter, r *rest.Request) {

}

// PUT     /v1/queue/:queueId/message/:mid                  @controllers.MessageController.push(queueId: String, mid: String, cid: Option[String] ?= None, gid: Option[String] ?= None)
func (c *Controller) PushWithId(w rest.ResponseWriter, r *rest.Request) {

}

// POST    /v1/queue/:queueId/message                       @controllers.MessageController.pushAnonymous(queueId: String,     cid: Option[String] ?= None, gid: Option[String] ?= None)
func (c *Controller) Push(w rest.ResponseWriter, r *rest.Request) {

}

// PUT     /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageController.keepAlive(queueId: String, mid: String, ack: String)
func (c *Controller) KeepAlive(w rest.ResponseWriter, r *rest.Request) {

}

// DELETE  /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageController.acknowledge(queueId: String, mid: String, ack: String)
func (c *Controller) Ack(w rest.ResponseWriter, r *rest.Request) {

}

