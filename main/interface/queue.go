package _interface

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/totemcaf/quongo/main/model"
	"net/http"
	"time"
)

type QueueInteractor interface {
	FindAll(offset int, limit int) ([]*model.Queue, error)
	FindById(queueId string) (*model.Queue, error)
	Complete(queue *model.Queue) *model.QueueWithStats
	IsQueueNameValid(queueId string) bool
	Add(queue model.Queue) (*model.Queue, error)
	Update(queue model.Queue) (*model.Queue, error)
}

type QueueView struct {
	interactor QueueInteractor
	Routes     []*rest.Route
}

func NewQueueView(interactor QueueInteractor) *QueueView {

	ctrl := QueueView{interactor: interactor}

	routes := []*rest.Route{
		rest.Get("/v1/queue", ctrl.QueueAll),
		rest.Get("/v1/queue/#queueId", ctrl.QueueGet),
		// Queue services),
		rest.Post("/v1/queue", ctrl.QueueAdd),
		rest.Put("/v1/queue/#queueId", ctrl.QueueUpd),
	}

	ctrl.Routes = routes

	return &ctrl
}

// GET  	/v1/queue                                        @controllers.Controller.findAll()
func (c *QueueView) QueueAll(w rest.ResponseWriter, r *rest.Request) {

	qs, e := c.interactor.FindAll(0, 100)

	if e == nil {
		w.WriteJson(qs)
	} else {
		rest.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

// GET  	/v1/queue/:queueId                               @controllers.Controller.find(queueId: String)
func (c *QueueView) QueueGet(w rest.ResponseWriter, r *rest.Request) {
	qId := r.PathParam("queueId")

	q, e := c.interactor.FindById(qId)

	if e == nil {
		w.WriteJson(c.interactor.Complete(q))
	} else {
		rest.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func payloadAsQueue(w rest.ResponseWriter, r *rest.Request) (*model.Queue, bool) {
	var q model.Queue

	e1 := r.DecodeJsonPayload(&q)

	if e1 != nil {
		rest.Error(w, e1.Error(), http.StatusBadRequest)
		return nil, false
	}

	return &q, true
}

// POST 	/v1/queue                                        @controllers.Controller.create(queueName: String)
func (c *QueueView) QueueAdd(w rest.ResponseWriter, r *rest.Request) {
	if queue, ok := payloadAsQueue(w, r); ok {
		added, e2 := c.interactor.Add(*queue)

		if e2 == nil {
			w.WriteJson(added)
		} else {
			rest.Error(w, e2.Error(), http.StatusBadRequest)
		}
	}
}

// #PUT  /v1/queue/:queueId                                 @controllers.Controller.Update(queueId: String)
func (c *QueueView) QueueUpd(w rest.ResponseWriter, r *rest.Request) {
	if queue, ok := payloadAsQueue(w, r); ok {
		queue.Name = r.PathParam("queueId")
		queue.Created = time.Now()

		updated, errU := c.interactor.Update(*queue)

		if errU == nil {
			w.WriteJson(updated)
		} else {
			rest.Error(w, errU.Error(), http.StatusBadRequest)
		}
	}
}

// #DELETE /v1/queue/:queueId                               @controllers.Controller.delete(queueId: String)
func (c *QueueView) Delete(w rest.ResponseWriter, r *rest.Request) {

}
