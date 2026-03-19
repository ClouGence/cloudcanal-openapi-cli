package testsupport

import (
	"io"
	"strings"
)

type TestConsole struct {
	inputs []string
	output strings.Builder
	index  int
}

func NewTestConsole(inputs ...string) *TestConsole {
	return &TestConsole{inputs: inputs}
}

func (t *TestConsole) ReadLine(prompt string) (string, error) {
	t.output.WriteString(prompt)
	if t.index >= len(t.inputs) {
		return "", io.EOF
	}
	value := t.inputs[t.index]
	t.index++
	return value, nil
}

func (t *TestConsole) ReadSecret(prompt string) (string, error) {
	t.output.WriteString(prompt)
	if t.index >= len(t.inputs) {
		return "", io.EOF
	}
	value := t.inputs[t.index]
	t.index++
	t.output.WriteString("\n")
	return value, nil
}

func (t *TestConsole) Println(text string) {
	t.output.WriteString(text)
	t.output.WriteString("\n")
}

func (t *TestConsole) Output() string {
	return t.output.String()
}
