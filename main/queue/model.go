package queue

import (
"time"
)

type Queue struct {
  Name      string        `json:"name" bson:"_id"`
  Created   time.Time     `json:"created"`
  VisWnd    time.Duration `json:"visibilityWindow" bson:"visibilityWindow"`
}

func NewQueue(name string) (*Queue, error) {
  wnd, _ := time.ParseDuration("30s")

  return &Queue{ Name: name, Created: time.Now(), VisWnd: wnd }, nil
}

type QueueWithStats struct {
  *Queue
  Stats     QueueStats  `json:"stats"`
}

type QueueStats struct {
  Total     int       `json:"total"`
  Hidden    int       `json:"hidden"`
  InProcess int       `json:"inProcess"`
}

func NewQueueWithStat(queue *Queue, total int, hidden int, inProcess int) *QueueWithStats {
  return &QueueWithStats{queue, QueueStats{total, hidden, inProcess}}
}
