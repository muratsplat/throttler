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
	mux   sync.Mutex
	tasks map[int]task.Task
}

func NewExecuter(list []task.Task) Throttler {
	mapList := make(map[int]task.Task)
	for k, task := range list {
		mapList[k] = task
	}
	return &Executer{
		tasks: mapList,
	}
}

func (e *Executer) Execute(number int) (error, *Result) {
	e.mux.Lock()
	defer e.mux.Unlock()
	result := make(chan int)
	var errList Result

	if len(e.tasks) == 0 {
		return ErrNothingToCall, nil
	}

	for i := 0; i < number; i++ {
		for k, t := range e.tasks {
			go func(result chan int, idx int, t task.Task) {
				err := t.Do()
				errList = append(errList, err)

				result <- idx
			}(result, k, t)
		}
	}

	for i := 0; i < number; i++ {
		doneIdx := <-result
		delete(e.tasks, doneIdx)
	}

	return nil, &errList
}
