package main

import (
  "time"
  "net/http"
)

const ISO_DATE_PATTERN = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"

type Message struct {
  Id      string        `json:"_id"`
  Payload string        `json:"payload"`
  Created time.Time     `json:"created"`
  Visible time.Time     `json:"visible"`
  Cid     string        `json:"cid"`
  Gid     string        `json:"gif"`
  Holder  string        `json:"holder"`
  Ack     string        `json:"Ack"`
}

// GET     /v1/queue/:queueId/message                       @controllers.MessageController.findAll(queueId: String)
func MsgAll(w http.ResponseWriter, r *http.Request) {

}

// GET     /v1/queue/:queueId/next                          @controllers.MessageController.pop(queueId: String, holder: Option[String] ?= None, window: Option[Int] ?= None)
func MsgPop(w http.ResponseWriter, r *http.Request) {

}

// GET     /v1/queue/:queueId/message/:mid                  @controllers.MessageController.findMsg(queueId: String, mid: String)
func MsgGet(w http.ResponseWriter, r *http.Request) {

}

// PUT     /v1/queue/:queueId/message/:mid                  @controllers.MessageController.push(queueId: String, mid: String, cid: Option[String] ?= None, gid: Option[String] ?= None)
func MsgPushId(w http.ResponseWriter, r *http.Request) {

}

// POST    /v1/queue/:queueId/message                       @controllers.MessageController.pushAnonymous(queueId: String,     cid: Option[String] ?= None, gid: Option[String] ?= None)
func MsgPush(w http.ResponseWriter, r *http.Request) {

}

// PUT     /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageController.keepAlive(queueId: String, mid: String, ack: String)
func MsgKeepAlive(w http.ResponseWriter, r *http.Request) {

}

// DELETE  /v1/queue/:queueId/message/:mid/ack/:ack         @controllers.MessageController.acknowledge(queueId: String, mid: String, ack: String)
func MsgAck(w http.ResponseWriter, r *http.Request) {

}
