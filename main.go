package main

import (
	"bufio"

	"github.com/muratsplat/throttler/executer"
	"github.com/muratsplat/throttler/task"

	"os"
)

func main() {
	var tasklist []task.Task
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l := scanner.Text()
		tasklist = append(tasklist, task.NewPerson(&l, os.Stdout))
	}

	exe := executer.NewExecuter(tasklist)
	times := len(tasklist)
	done := make(chan error)

	for i := 0; i < times; i++ {
		go func() {
			err, _ := exe.Execute(2)
			//  Maybe error hanfling
			if err != nil {
				if err == executer.ErrNothingToCall {
					done <- err
					return
				}
				// an external error
				panic(err)
			}
			done <- err

		}()
	}

	for i := 0; i < times; i++ {
		<-done
	}
}
