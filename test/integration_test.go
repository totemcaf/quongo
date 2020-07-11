package test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/totemcaf/quongo/app"
	"github.com/totemcaf/quongo/app/utils"
)

// This tests check the App as a component. It used the memory implementation of the repositories
var logger = log.New(os.Stderr, "", log.LstdFlags)

// Check Ping endpoint
func TestPing(t *testing.T) {
	newTest().
		Get("/ping").
		Expect(t).
		Status(http.StatusOK).
		Body(`"ok"`). // expectations
		End()
}

// Check info endpoint
func TestInfo(t *testing.T) {
	newTest().
		Get("/info").
		Expect(t).
		Status(http.StatusOK).
		Body(`{
			"name": "quongo",
			"version": "0.0.1",
			"goVersion": "go1.13.4"
		}`). // expectations
		End()
}
func TestSendMessage(t *testing.T) {
	newTest().
		Put("/api/v1/queues/queue-1/messages/1010101").
		Header("Content-type", "application/json").
		Body(`{
				"value": 12
			}`).
		Expect(t). // expectations
		Status(http.StatusCreated).
		Header("X-Mid", "1010101").
		Body("").
		End()
}

func TestGetMessageFromEmptyQueue(t *testing.T) {
	newTest().
		Get("/api/v1/queues/queue-1/pop").
		Expect(t). // expectations
		Status(http.StatusNotFound).
		End()
}

func TestSendMessageAndGetIt(t *testing.T) {
	handler := makeHandler()

	givenMessageIsSent(t, handler)

	// THEN message is received
	newTestWith(handler).
		Get("/api/v1/queues/queue-1/pop").
		Expect(t). // expectations
		Status(http.StatusOK).
		Header("X-Mid", "1010101").
		Header("X-Retries", "0").
		Header("X-Created", "1963-11-29T04:27:42.000Z").
		HeaderPresent("X-Ack").
		//	headers.Add("X-Ack", msg.Ack)
		Body(`{"value": 12}`).
		End()
}

func givenMessageIsSent(t *testing.T, handler http.Handler) {
	newTestWith(handler).
		Put("/api/v1/queues/queue-1/messages/1010101").
		Header("Content-type", "application/json").
		Body(`{
				"value": 12
			}`).
		Expect(t). // expectations
		Status(http.StatusCreated).
		End()
}

const aSpecialDatetime = "1963-11-29 04:27:42"

func makeHandler() http.Handler {
	return app.NewApp("memory", utils.FixedClockAt(aSpecialDatetime), logger).Handler()
}

func newTestWith(handler http.Handler) *apitest.APITest {
	return apitest.
		New().
		Handler(handler)
}

func newTest() *apitest.APITest {
	return newTestWith(makeHandler())
}
