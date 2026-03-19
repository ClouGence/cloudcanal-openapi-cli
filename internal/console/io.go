package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/peterh/liner"
	"golang.org/x/term"
)

type IO interface {
	ReadLine(prompt string) (string, error)
	ReadSecret(prompt string) (string, error)
	Println(text string)
	ClearScreen()
}

type Completer func(line string) []string

type TabCompletable interface {
	SetCompleter(Completer)
}

type StdIO struct {
	reader     *bufio.Reader
	writer     io.Writer
	inputFile  *os.File
	outputFile *os.File
	liner      *liner.State
}

func NewStdIO(reader io.Reader, writer io.Writer) *StdIO {
	inputFile, _ := reader.(*os.File)
	outputFile, _ := writer.(*os.File)
	var lineEditor *liner.State
	if inputFile != nil && outputFile != nil && term.IsTerminal(int(inputFile.Fd())) && term.IsTerminal(int(outputFile.Fd())) {
		lineEditor = liner.NewLiner()
		lineEditor.SetCtrlCAborts(true)
	}
	return &StdIO{
		reader:     bufio.NewReader(reader),
		writer:     writer,
		inputFile:  inputFile,
		outputFile: outputFile,
		liner:      lineEditor,
	}
}

func (s *StdIO) ReadLine(prompt string) (string, error) {
	if s.liner != nil {
		line, err := s.liner.Prompt(prompt)
		if err != nil {
			return "", err
		}
		if trimmed := trimLine(line); trimmed != "" {
			s.liner.AppendHistory(trimmed)
		}
		return trimLine(line), nil
	}

	if _, err := fmt.Fprint(s.writer, prompt); err != nil {
		return "", err
	}
	line, err := s.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(line) > 0 {
			return trimLine(line), nil
		}
		return "", err
	}
	return trimLine(line), nil
}

func (s *StdIO) ReadSecret(prompt string) (string, error) {
	if s.liner != nil {
		line, err := s.liner.PasswordPrompt(prompt)
		if err != nil {
			return "", err
		}
		return trimLine(line), nil
	}

	if s.inputFile == nil || !term.IsTerminal(int(s.inputFile.Fd())) {
		return s.ReadLine(prompt)
	}

	if _, err := fmt.Fprint(s.writer, prompt); err != nil {
		return "", err
	}
	line, err := term.ReadPassword(int(s.inputFile.Fd()))
	if _, printErr := fmt.Fprintln(s.writer); err == nil && printErr != nil {
		return "", printErr
	}
	if err != nil {
		return "", err
	}
	return trimLine(string(line)), nil
}

func (s *StdIO) Println(text string) {
	_, _ = fmt.Fprintln(s.writer, text)
}

func (s *StdIO) ClearScreen() {
	_, _ = fmt.Fprint(s.writer, "\033[H\033[2J")
}

func (s *StdIO) SetCompleter(completer Completer) {
	if s.liner == nil {
		return
	}
	if completer == nil {
		s.liner.SetCompleter(nil)
		return
	}
	s.liner.SetCompleter(func(line string) []string {
		return completer(line)
	})
}

func (s *StdIO) Close() error {
	if s.liner == nil {
		return nil
	}
	return s.liner.Close()
}

func IsPromptAborted(err error) bool {
	return errors.Is(err, liner.ErrPromptAborted)
}

func trimLine(line string) string {
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line
}
