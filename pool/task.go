package pool

type Task struct {
	Handler func(para ...interface{})
	Parameters []interface{}
}
