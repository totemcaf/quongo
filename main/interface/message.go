package _interface

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/totemcaf/quongo/main/model"
	util "github.com/totemcaf/quongo/main/utils"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

type MessageInteractor interface {
	FindAll(queueId string, offset int, limit int) ([]model.Message, error)
	FindById(queueId string, msgId string) (*model.Message, error)
	Add(queueId string, message *model.Message) (*model.Message, error)

	Pop(queueId string, delay time.Duration) (*model.Message, error)
}

type MessageView struct {
	messageInt MessageInteractor
	queueInt   QueueInteractor
	Routes     []*rest.Route
}

func NewMessageView(messageInteractor MessageInteractor, queueInteractor QueueInteractor) *MessageView {
	ctrl := MessageView{messageInteractor, queueInteractor, nil}

	routes := []*rest.Route{
		rest.Get("/v1/queue/#queueId/message", ctrl.All),
		rest.Get("/v1/queue/#queueId/next", ctrl.Pop),
		rest.Get("/v1/queue/#queueId/message/#mid", ctrl.Get),
		rest.Put("/v1/queue/#queueId/message/#mid", ctrl.PushWithId),
		rest.Post("/v1/queue/#queueId/message", ctrl.Push),
		rest.Put("/v1/queue/#queueId/message/#mid/ack/#ack", ctrl.KeepAlive),
		rest.Delete("/v1/queue/#queueId/message/#mid/ack/#ack", ctrl.Ack),
	}

	ctrl.Routes = routes

	return &ctrl
}

// GET     /v1/queue/:queueId/message ?size= & page =       @controllers.MessageView.findAll(queueId: String)
func (c *MessageView) All(w rest.ResponseWriter, r *rest.Request) {
	queueId := r.PathParam("queueId")
	page, e1 := util.IntParam(r, "page", 0)
	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusBadRequest)
		return
	}
	size, e2 := util.IntParam(r, "size", 10)
	if e2 != nil {
		rest.Error(w, e2.Error(), http.StatusBadRequest)
		return
	}

	messages, e3 := c.messageInt.FindAll(queueId, size*page, size)

	if e3 != nil {
		rest.Error(w, e3.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else {
		w.WriteJson(messages)
	}
}

// GET     /v1/queue/:queueId/next                          @controllers.MessageView.pop(queueId: String, holder: Option[String] ?= None, window: Option[Int] ?= None)
func (c *MessageView) Pop(w rest.ResponseWriter, r *rest.Request) {
	queueId := r.PathParam("queueId")

	q, e1 := c.queueInt.FindById(queueId)

	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusNotFound)
	}

	msg, e2 := c.messageInt.Pop(queueId, q.VisWnd)

	if e2 != nil {
		rest.Error(w, e2.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else if msg != nil {
		rest.Error(w, "No pending messages", http.StatusNotFound)
	} else {
		w.WriteJson(msg)
	}
}

// GET     /v1/queue/:queueId/message/:mid                  @controllers.MessageView.findMsg(queueId: String, mid: String)
func (c *MessageView) Get(w rest.ResponseWriter, r *rest.Request) {
	queueId := r.PathParam("queueId")
	msgId := r.PathParam("mid")

	msg, e1 := c.messageInt.FindById(queueId, msgId)

	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else if msg != nil {
		rest.Error(w, "Message not found", http.StatusNotFound)
	} else {
		w.WriteJson(msg)
	}
}

// PUT     /v1/queue/:queueId/message/:mid                  @controllers.MessageView.push(queueId: String, mid: String, cid: Option[String] ?= None, gid: Option[String] ?= None)
func (c *MessageView) PushWithId(w rest.ResponseWriter, r *rest.Request) {
	mid := r.PathParam("mid")
	if !bson.IsObjectIdHex(mid) {
		rest.Error(w, "Invalid message id", http.StatusBadRequest)
	}
	c.pushOne(w, r, bson.ObjectIdHex(mid))
}

// POST    /v1/queue/:queueId/message                       @controllers.MessageView.pushAnonymous(queueId: String,     cid: Option[String] ?= None, gid: Option[String] ?= None)
func (c *MessageView) Push(w rest.ResponseWriter, r *rest.Request) {
	c.pushOne(w, r, bson.NewObjectId())
}

func (c *MessageView) pushOne(w rest.ResponseWriter, r *rest.Request, id bson.ObjectId) {
	now := time.Now()
	var visible time.Time
	timeStr := r.FormValue("time")

	if timeStr == "" {
		visible = now
	} else {
		var e1 error
		visible, e1 = time.Parse(util.TimeLayout, timeStr)
		if e1 != nil {
			rest.Error(w, "Invalid time", http.StatusBadRequest)
			return
		}
	}

	payload, e2 := ioutil.ReadAll(r.Body)

	if e2 != nil {
		rest.Error(w, e2.Error(), http.StatusBadRequest)
		return
	}

	msg := model.Message{
		Id:         id,
		Visible:    visible,
		Created:    now,
		Ack:        "",
		Cid:        "",
		Gid:        "",
		Holder:     "",
		Payload:    string(payload),
		Programmed: visible,
		Retries:    0,
	}

	queueId := r.PathParam("queueId")
	newMsg, e3 := c.messageInt.Add(queueId, &msg)

	if e3 != nil {
		rest.Error(w, e2.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(newMsg)
}

// PUT     /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageView.keepAlive(queueId: String, mid: String, ack: String)
func (c *MessageView) KeepAlive(w rest.ResponseWriter, r *rest.Request) {

}

// DELETE  /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageView.acknowledge(queueId: String, mid: String, ack: String)
func (c *MessageView) Ack(w rest.ResponseWriter, r *rest.Request) {

}
