package task

import (
	"bytes"
	"io"
	"testing"
)

var (
	_ Task = &Person{}

	saySomethingInTest = []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Ut tincidunt quam venenatis mi pharetra blandit.",
		"Nullam aliquet quam et ipsum maximus, nec auctor risus ornare.",
		"Sed at dui non velit molestie malesuada.",
	}
)

func createTestTasks(to io.Writer) (tasklist []Task) {
	for _, speach := range saySomethingInTest {
		tasklist = append(tasklist, NewPerson(&speach, to))
	}
	return tasklist
}

func TestPersonSimple(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	if buffer.Len() != 0 {
		t.Fatalf("Expected: zero but got %d", buffer.Len())
	}

	tasks := createTestTasks(buffer)
	for _, task := range tasks {
		err := task.Do()
		if err != nil {
			t.Fatal(err)
		}
	}
	if buffer.Len() < 1 {
		t.Fatalf("Expected: more than zero but got %d", buffer.Len())
	}

}
