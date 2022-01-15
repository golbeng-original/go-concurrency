package workerpool

import (
	"context"
	"sync"
)

type TaskPoolContext struct {
	maxWorker  int
	queuedTask chan Job
	wait       sync.WaitGroup
	dones      map[int]chan interface{}
	ctx        context.Context
	cancel     context.CancelFunc
}

func (p *TaskPoolContext) Run() {

	for i := 0; i < p.maxWorker; i++ {

		//done := make(chan interface{})
		//p.dones[i] = done

		p.wait.Add(1)
		go func() {
			defer p.wait.Done()

			//workerCtx, cancel := context.WithCancel(p.ctx)
			//defer cancel()

			for {
				select {
				case <-p.ctx.Done():
					return
				case job, ok := <-p.queuedTask:
					if !ok {
						return
					}

					job.task(job.arg...)
				}
			}
		}()
	}

	go func() {
		<-p.ctx.Done()
		close(p.queuedTask)
	}()
}

func (p *TaskPoolContext) AddTask(task func(args ...interface{}), arg ...interface{}) {

	if p.IsDone() {
		return
	}

	job := Job{
		task: task,
		arg:  make([]interface{}, 0),
	}

	job.arg = append(job.arg, arg...)

	select {
	case <-p.ctx.Done():
		return
	case p.queuedTask <- job:
	}

}

func (p *TaskPoolContext) Wait() {
	p.wait.Wait()
}

func (p *TaskPoolContext) Done() {
	p.cancel()
}

func (p *TaskPoolContext) IsDone() bool {
	// 취소 된 상태..
	if p.ctx.Err() != nil {
		//fmt.Println("addTask : WorkerPool Canceled")
		return true
	}

	return false
}

func NewPoolContext(maxWorker int, ctx context.Context) *TaskPoolContext {

	poolCtx, cancel := context.WithCancel(ctx)

	return &TaskPoolContext{
		maxWorker:  maxWorker,
		queuedTask: make(chan Job),
		wait:       sync.WaitGroup{},
		dones:      make(map[int]chan interface{}),
		ctx:        poolCtx,
		cancel:     cancel,
	}
}
