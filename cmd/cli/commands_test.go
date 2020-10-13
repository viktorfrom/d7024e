package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelp(t *testing.T) {
	content := Prompt()
	assert.Equal(t, content, cmdTester("help"))
}

func TestHelpShort(t *testing.T) {
	content := Prompt()
	assert.Equal(t, content, cmdTester("h"))
}

func cmdTester(cmd string) string {
	out = bytes.NewBuffer(nil)

	Commands(out, nil, []string{cmd})
	return trimWriterOutput(out)
}

func trimWriterOutput(out io.Writer) string {
	str := out.(*bytes.Buffer).String()
	return strings.TrimSuffix(str, "\n")
}

func TestExit(t *testing.T) {
	assert.Equal(t, 3, cmdExit("exit"))
}

func TestExitShort(t *testing.T) {
	assert.Equal(t, 3, cmdExit("e"))
}

// This is too crazy to come up with by oneself...
// Source: https://stackoverflow.com/questions/40615641/testing-os-exit-scenarios-in-go-with-coverage-information-coveralls-io-goverall/40801733#40801733
func cmdExit(exit string) int {
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()

	var got int
	myExit := func(code int) {
		got = code
	}

	osExit = myExit
	Commands(out, nil, []string{exit})
	return got
}

func TestDefault(t *testing.T) {
	assert.Equal(t, errInvalidCmd, cmdTester(""))
}

func TestPut(t *testing.T) {
	assert.Equal(t, errNoArg, cmdTester("put"))
}

func TestPutShort(t *testing.T) {
	assert.Equal(t, errNoArg, cmdTester("p"))
}

func TestGet(t *testing.T) {
	assert.Equal(t, errNoArg, cmdTester("get"))
}

func TestGetShort(t *testing.T) {
	assert.Equal(t, errNoArg, cmdTester("g"))
}
