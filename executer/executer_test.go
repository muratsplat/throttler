package executer

import (
	"github.com/golang/mock/gomock"
	"github.com/muratsplat/throttler/mocks"
	"github.com/muratsplat/throttler/task"

	"testing"
)

var (
	_ Throttler = &Executer{}
)

func TestExecuterSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	taks1 := mocks.NewMockTask(ctrl)
	taks1.EXPECT().Do().Times(1).Return(nil)

	taks2 := mocks.NewMockTask(ctrl)
	taks2.EXPECT().Do().Times(1).Return(nil)

	taks3 := mocks.NewMockTask(ctrl)
	taks3.EXPECT().Do().Times(1).Return(nil)

	taks4 := mocks.NewMockTask(ctrl)
	taks4.EXPECT().Do().Times(1).Return(nil)

	taks5 := mocks.NewMockTask(ctrl)
	taks5.EXPECT().Do().Times(1).Return(nil)

	throttler := NewExecuter([]task.Task{
		taks1,
		taks2,
		taks3,
		taks4,
		taks5,
	})

	err, results := throttler.Execute(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(*results) != 1 {
		t.Fatalf("Expexted: one, but got: %d", len(*results))
	}

	err, results = throttler.Execute(4)
	if err != nil {
		t.Fatal(err)
	}
	if len(*results) != 4 {
		t.Fatalf("Expexted: 4, but got: %d", len(*results))
	}

	err, results = throttler.Execute(1)
	if err != nil && err != ErrNothingToCall {
		t.Fatal(err)
	}

}

func TestInSafeExecuterSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	taks1 := mocks.NewMockTask(ctrl)
	taks1.EXPECT().Do().Times(1).Return(nil)

	taks2 := mocks.NewMockTask(ctrl)
	taks2.EXPECT().Do().Times(1).Return(nil)

	taks3 := mocks.NewMockTask(ctrl)
	taks3.EXPECT().Do().Times(1).Return(nil)

	taks4 := mocks.NewMockTask(ctrl)
	taks4.EXPECT().Do().Times(1).Return(nil)

	taks5 := mocks.NewMockTask(ctrl)
	taks5.EXPECT().Do().Times(1).Return(nil)

	taks6 := mocks.NewMockTask(ctrl)
	taks6.EXPECT().Do().Times(1).Return(nil)

	throttler := NewExecuter([]task.Task{
		taks1,
		taks2,
		taks3,
		taks4,
		taks5,
		taks6,
	})

	done := make(chan int)
	go func() {
		err, results := throttler.Execute(3)
		if err != nil {
			t.Fatal(err)
		}
		if len(*results) != 3 {
			t.Fatalf("Expexted: one, but got: %d", len(*results))
		}
		done <- len(*results)

	}()

	go func() {
		err, results := throttler.Execute(3)
		if err != nil {
			t.Fatal(err)
		}
		if len(*results) != 3 {
			t.Fatalf("Expexted: 4, but got: %d", len(*results))
		}
		done <- len(*results)
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
}
