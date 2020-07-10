package delivery

import (
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/totemcaf/quongo/main/model"
)

// QueueInteractor ...
type QueueInteractor interface {
	FindAll(offset int, limit int) ([]*model.Queue, error)
	FindByID(queueID string) (*model.Queue, error)
	Complete(queue *model.Queue) *model.QueueWithStats
	IsQueueNameValid(queueID string) bool
	Add(queue model.Queue) (*model.Queue, error)
	Update(queue model.Queue) (*model.Queue, error)
}

// QueueView ...
type QueueView struct {
	interactor QueueInteractor
	routes     []*rest.Route
}

// NewQueueView ...
func NewQueueView(interactor QueueInteractor) *QueueView {

	ctrl := QueueView{interactor: interactor}

	routes := []*rest.Route{
		rest.Get("/v1/queue", ctrl.QueueAll),
		rest.Get("/v1/queue/#queueId", ctrl.QueueGet),
		// Queue services),
		rest.Post("/v1/queue", ctrl.QueueAdd),
		rest.Put("/v1/queue/#queueId", ctrl.QueueUpd),
	}

	ctrl.routes = routes

	return &ctrl
}

// Routes ...
func (v *QueueView) Routes() []*rest.Route {
	return v.routes
}

// QueueAll GET  	/v1/queue                                        @controllers.Controller.findAll()
func (v *QueueView) QueueAll(w rest.ResponseWriter, r *rest.Request) {

	qs, e := v.interactor.FindAll(0, 100)

	if e == nil {
		w.WriteJson(qs)
	} else {
		rest.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

// QueueGet GET  	/v1/queue/:queueId                               @controllers.Controller.find(queueId: String)
func (v *QueueView) QueueGet(w rest.ResponseWriter, r *rest.Request) {
	qID := r.PathParam("queueId")

	q, e := v.interactor.FindByID(qID)

	if e == nil {
		w.WriteJson(v.interactor.Complete(q))
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

// QueueAdd POST 	/v1/queue                                        @controllers.Controller.create(queueName: String)
func (v *QueueView) QueueAdd(w rest.ResponseWriter, r *rest.Request) {
	if queue, ok := payloadAsQueue(w, r); ok {
		added, e2 := v.interactor.Add(*queue)

		if e2 == nil {
			w.WriteJson(added)
		} else {
			rest.Error(w, e2.Error(), http.StatusBadRequest)
		}
	}
}

// QueueUpd #PUT  /v1/queue/:queueId                                 @controllers.Controller.Update(queueId: String)
func (v *QueueView) QueueUpd(w rest.ResponseWriter, r *rest.Request) {
	if queue, ok := payloadAsQueue(w, r); ok {
		queue.Name = r.PathParam("queueId")
		queue.Created = time.Now()

		updated, errU := v.interactor.Update(*queue)

		if errU == nil {
			w.WriteJson(updated)
		} else {
			rest.Error(w, errU.Error(), http.StatusBadRequest)
		}
	}
}

// Delete #DELETE /v1/queue/:queueId                               @controllers.Controller.delete(queueId: String)
func (v *QueueView) Delete(w rest.ResponseWriter, r *rest.Request) {

}
