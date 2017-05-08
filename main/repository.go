package main

import (
  "container/list"
  "time"
)

type QueueRepository struct {

}

func NewQueueRepository() *QueueRepository {
  return &QueueRepository{}
}

func (r *QueueRepository) findAll() *list.List {
  return list.New()
}

func (r *QueueRepository) findById(id string) (*QueueWithStats, error) {
  wnd, _ := time.ParseDuration("30s")

  q := QueueWithStats{
    Queue: &Queue{ id, time.Now(), wnd },
    Stats: &QueueStats{10, 2, 3},
  }

  return &q, nil
}

