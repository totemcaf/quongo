package queue

import (
"time"
)

type Queue struct {
  Name      string        `json:"_id"`
  Created   time.Time     `json:"created"`
  VisWnd    time.Duration `json:"visibilityWindow"`
}

func NewQueue(name string) (*Queue, error) {
  wnd, _ := time.ParseDuration("30s")

  return &Queue{ Name: name, Created: time.Now(), VisWnd: wnd }, nil
}

type QueueWithStats struct {
  *Queue
  Stats     *QueueStats  `json:"stats"`
}

type QueuesWithStats []QueueWithStats

type QueueStats struct {
  Total     int16       `json:"total"`
  Hidden    int16       `json:"hidden"`
  InProcess int16       `json:"inProcess"`
}

