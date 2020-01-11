package task

// go:generate mockgen -destination=../mocks/mock_task.go -package=mocks -source=task.go

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrNoSpeachToSay = errors.New("not word to say")
)

type Person struct {
	MySpeach *string
	To       io.Writer
}

func NewPerson(speach *string, to io.Writer) *Person {
	if to == nil {
		to = os.Stdout
	}
	return &Person{
		MySpeach: speach,
		To:       to,
	}
}

func (p *Person) Do() error {
	if p.MySpeach == nil {
		return ErrNoSpeachToSay
	}

	speach := fmt.Sprintf("%s\n", *p.MySpeach)
	_, err := p.To.Write([]byte(speach))
	if err != nil {
		return err
	}
	return nil
}

type Task interface {
	Do() error
}

type TaskList struct {
	Queue *[]Task
}
