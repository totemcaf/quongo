package queue

type Repository struct {

}

func NewQueueRepository() *Repository {
  return &Repository{}
}


func NewQueueStats(name string) *QueueWithStats {
  pQ, _ := NewQueue(name)

  return &QueueWithStats{
    Queue: pQ,
    Stats: &QueueStats{10, 2, 3},
  }
}

func (r *Repository) findAll() (*[] *QueueWithStats, error) {
  l := [] *QueueWithStats{
    NewQueueStats("one"),
    NewQueueStats("two"),
  }
  return &l, nil
}

func (r *Repository) findById(id string) (*QueueWithStats, error) {
  return NewQueueStats(id), nil
}

func (r *Repository) add(pQueue *Queue) (*Queue, error) {
  return pQueue, nil
}

func (r *Repository) Update(pQueue *Queue) (*Queue, error) {
  return pQueue, nil
}
