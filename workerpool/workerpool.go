package workerpool

import "sync"

type TaskPool struct {
	maxWorker  int
	queuedTask chan Job
	wait       sync.WaitGroup
	dones      map[int]chan interface{}
}

func (p *TaskPool) Run() {

	for i := 0; i < p.maxWorker; i++ {

		done := make(chan interface{})
		p.dones[i] = done

		p.wait.Add(1)
		// worker
		go func(done <-chan interface{}) {
			defer p.wait.Done()

			for {

				select {
				case <-done:
					return
				case job, ok := <-p.queuedTask:
					if !ok {
						return
					}

					job.task(job.arg...)
				}

			}
		}(done)
	}
}

func (p *TaskPool) AddTask(task func(args ...interface{}), arg ...interface{}) {
	job := Job{
		task: task,
		arg:  make([]interface{}, 0),
	}

	job.arg = append(job.arg, arg...)

	p.queuedTask <- job
}

func (p *TaskPool) Wait() {
	p.wait.Wait()
}

func (p *TaskPool) Done() {

	for _, done := range p.dones {
		done <- true
	}
}

func NewPool(maxWorker int) *TaskPool {
	return &TaskPool{
		maxWorker:  maxWorker,
		queuedTask: make(chan Job),
		wait:       sync.WaitGroup{},
		dones:      make(map[int]chan interface{}),
	}
}
