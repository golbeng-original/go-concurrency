package workerpool

type Job struct {
	task func(args ...interface{})
	arg  []interface{}
}
