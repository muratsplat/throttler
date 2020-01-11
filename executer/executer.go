package executer

import (
	"errors"
	"sync"

	"github.com/muratsplat/throttler/task"
)

var (
	ErrNothingToCall = errors.New("nothing to call")
)

type Result []error

type Throttler interface {
	Execute(number int) (error, *Result)
}

type Executer struct {
	mux   sync.RWMutex
	tasks *TaskListInSafe
}

type TaskListInSafe struct {
	list map[int]task.Task
	mux  sync.Mutex
}

func NewExecuter(list []task.Task) Throttler {
	mapList := make(map[int]task.Task)
	for k, task := range list {
		mapList[k] = task
	}
	return &Executer{
		tasks: &TaskListInSafe{
			list: mapList,
		},
	}
}

func (e *Executer) pop() (error, task.Task) {
	e.tasks.mux.Lock()
	defer e.tasks.mux.Unlock()
	if len(e.tasks.list) == 0 {
		return ErrNothingToCall, nil
	}

	for k, one := range e.tasks.list {
		delete(e.tasks.list, k)
		return nil, one
	}

	return ErrNothingToCall, nil

}

func (e *Executer) Execute(number int) (error, *Result) {
	result := make(chan int)
	var errList Result
	last := 0
	for i := 0; i < number; i++ {
		err, t := e.pop()
		if err == ErrNothingToCall {
			break
		}
		last = i + 1
		go func(result chan int, idx int, t task.Task) {
			err := t.Do()
			errList = append(errList, err)
			result <- idx
		}(result, i, t)

	}
	for i := 0; i < last; i++ {
		<-result
	}

	return nil, &errList
}
