package delivery

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/totemcaf/quongo/main/model"
	"github.com/totemcaf/quongo/main/model/message"
	"github.com/totemcaf/quongo/main/utils"
)

// MessageInteractor ...
type MessageInteractor interface {
	FindAll(queueID string, offset int, limit int) ([]model.Message, error)
	FindByID(queueID string, msgID string) (*model.Message, error)
	Add(queueID string, message *model.Message) (*model.Message, error)

	Pop(queueID string, qunatity int) ([]*model.Message, error)
}

// MessageView todo
type MessageView struct {
	messageInt MessageInteractor
	queueInt   QueueInteractor
	routes     []*rest.Route
}

// NewMessageView TODO
func NewMessageView(messageInteractor MessageInteractor, queueInteractor QueueInteractor) *MessageView {
	ctrl := MessageView{messageInteractor, queueInteractor, nil}

	routes := []*rest.Route{
		rest.Put("/v1/queues/#queueId/messages/#mid", ctrl.PushWithID),

		rest.Get("/v1/queues/#queueId/messages", ctrl.All),
		rest.Get("/v1/queues/#queueId/pop", ctrl.Pop),
		rest.Get("/v1/queues/#queueId/messages/#mid", ctrl.Get),
		rest.Get("/v1/queues/#queueId/pop-many", ctrl.PopMany),
		rest.Post("/v1/queues/#queueId/messages", ctrl.Push),
		rest.Put("/v1/queues/#queueId/messages/#mid/ack/#ack", ctrl.KeepAlive),
		rest.Delete("/v1/queues/#queueId/messages/#mid/ack/#ack", ctrl.Ack),
		rest.Delete("/v1/queues/#queueId/messages/pop", ctrl.PopAck),
		rest.Delete("/v1/queues/#queueId/messages/pop-many", ctrl.PopAckMany),
	}

	ctrl.routes = routes

	return &ctrl
}

// Routes ...
func (v *MessageView) Routes() []*rest.Route {
	return v.routes
}

// All GET     /v1/queue/:queueId/message ?size= & page =       @controllers.MessageView.findAll(queueId: String)
func (v *MessageView) All(w rest.ResponseWriter, r *rest.Request) {
	queueID := r.PathParam("queueId")
	page, e1 := utils.IntParam(r, "page", 0)
	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusBadRequest)
		return
	}
	size, e2 := utils.IntParam(r, "size", 10)
	if e2 != nil {
		rest.Error(w, e2.Error(), http.StatusBadRequest)
		return
	}

	messages, e3 := v.messageInt.FindAll(queueID, size*page, size)

	if e3 != nil {
		rest.Error(w, e3.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else {
		w.WriteJson(messages)
	}
}

// Pop GET     /v1/queue/:queueId/pop                          @controllers.MessageView.pop(queueId: String, holder: Option[String] ?= None, window: Option[Int] ?= None)
func (v *MessageView) Pop(w rest.ResponseWriter, r *rest.Request) {
	queueID := r.PathParam("queueId")

	msg, e2 := v.messageInt.Pop(queueID, 1)

	if e2 != nil {
		rest.Error(w, e2.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else if msg == nil {
		rest.Error(w, "No pending messages", http.StatusNotFound)
	} else {
		w.WriteJson(msg)
	}
}

// Get     /v1/queue/:queueId/message/:mid                  @controllers.MessageView.findMsg(queueId: String, mid: String)
func (v *MessageView) Get(w rest.ResponseWriter, r *rest.Request) {
	queueID := r.PathParam("queueId")
	msgID := r.PathParam("mid")

	msg, e1 := v.messageInt.FindByID(queueID, msgID)

	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusInternalServerError) // TODO (caf) Discriminar si el error es interno o de datos
	} else if msg != nil {
		rest.Error(w, "Message not found", http.StatusNotFound)
	} else {
		w.WriteJson(msg)
	}
}

// PushWithID     /v1/queue/:queueId/message/:mid                  @controllers.MessageView.push(queueId: String, mid: String, cid: Option[String] ?= None, gid: Option[String] ?= None)
func (v *MessageView) PushWithID(w rest.ResponseWriter, r *rest.Request) {
	mid, err := message.ParseID(r.PathParam("mid"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	v.pushOne(w, r, mid)
}

// Push    /v1/queue/:queueId/message                       @controllers.MessageView.pushAnonymous(queueId: String,     cid: Option[String] ?= None, gid: Option[String] ?= None)
func (v *MessageView) Push(w rest.ResponseWriter, r *rest.Request) {
	v.pushOne(w, r, message.NewID())
}

func (v *MessageView) pushOne(w rest.ResponseWriter, r *rest.Request, id message.MID) {
	now := time.Now()
	var visible time.Time
	timeStr := r.FormValue("time")

	if timeStr == "" {
		visible = now
	} else {
		var e1 error
		visible, e1 = time.Parse(utils.TimeLayout, timeStr)
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
		ID:         id,
		Visible:    visible,
		Created:    now,
		Ack:        nil,
		Cid:        "",
		Gid:        "",
		Holder:     "",
		Payload:    string(payload),
		Programmed: visible,
		Retries:    0,
	}

	queueID := r.PathParam("queueId")
	newMsg, e3 := v.messageInt.Add(queueID, &msg)

	if e3 != nil {
		rest.Error(w, e2.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(newMsg)
}

// KeepAlive     /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageView.keepAlive(queueId: String, mid: String, ack: String)
func (v *MessageView) KeepAlive(w rest.ResponseWriter, r *rest.Request) {

}

// Ack  /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageView.acknowledge(queueId: String, mid: String, ack: String)
func (v *MessageView) Ack(w rest.ResponseWriter, r *rest.Request) {

}

// PopMany ...
func (v *MessageView) PopMany(w rest.ResponseWriter, r *rest.Request) {

}

// PopAck ...
func (v *MessageView) PopAck(w rest.ResponseWriter, r *rest.Request) {

}

// PopAckMany ...
func (v *MessageView) PopAckMany(w rest.ResponseWriter, r *rest.Request) {

}
