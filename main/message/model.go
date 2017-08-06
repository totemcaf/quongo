package message

import "time"

type Message struct {
  Id      string        `json:"_id"`
  Payload string        `json:"payload"`
  Created time.Time     `json:"created"`
  Visible time.Time     `json:"visible"`
  Cid     string        `json:"cid"`
  Gid     string        `json:"gif"`
  Holder  string        `json:"holder"`
  Retries int16         `json:"retries"`
  Ack     string        `json:"Ack"`
}
