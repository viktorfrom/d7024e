package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var respValue string
var respKey string
var respErr error

func TestHelp(t *testing.T) {
	content, _ := ioutil.ReadFile("prompt.txt")
	assert.Equal(t, cmdTester("help"), string(content))
}

func TestHelpShort(t *testing.T) {
	content, _ := ioutil.ReadFile("prompt.txt")
	assert.Equal(t, cmdTester("h"), string(content))
}

func TestHelpError(t *testing.T) {
	helpFile = "error.txt"
	out = bytes.NewBuffer(nil)

	// Save current function and restore at the end:
	oldLogFatal := logFatal
	defer func() { logFatal = oldLogFatal }()

	var gotV []interface{}
	myFatal := func(v ...interface{}) {
		gotV = v
	}

	logFatal = myFatal
	Help(out)
	expV := []interface{}{errNoFileFound + helpFile}

	assert.Equal(t, expV, gotV)
}

func cmdTester(cmd string) string {
	out = bytes.NewBuffer(nil)

	Commands(out, []string{cmd})
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
	Commands(out, []string{exit})
	return got
}

func TestDefault(t *testing.T) {
	assert.Equal(t, errInvalidCmd, cmdTester(""))
}

func TestPutCommand(t *testing.T) {
	assert.Equal(t, errWrongArg, cmdTester("put"))
}

func TestPutShortCommand(t *testing.T) {
	assert.Equal(t, errWrongArg, cmdTester("p"))
}

func TestGetCommand(t *testing.T) {
	assert.Equal(t, errWrongArg, cmdTester("get"))
}

func TestGetShortCommand(t *testing.T) {
	assert.Equal(t, errWrongArg, cmdTester("g"))
}

func TestPutCreated(t *testing.T) {
	status, loc, value, err := Put("localhost", "test", PostCreated)
	assert.Equal(t, "201 Created", status)
	assert.Equal(t, "/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", loc)
	assert.Equal(t, "test", value)
	assert.Nil(t, err)
}

// func TestGetBadRequest(t *testing.T) {
// 	status, loc, value, err := Get("localhost", "test")
// 	assert.Equal(t, "400 Bad Request", status)
// 	assert.Equal(t, "", loc)
// 	assert.Equal(t, "", value)
// 	assert.Error(t, err)
// }

func TestGet(t *testing.T) {
	respValue = "test"
	status, loc, value, err := Get("localhost", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", GetExisting)
	assert.Equal(t, "200 OK", status)
	assert.Equal(t, "/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", loc)
	assert.Equal(t, "test", value)
	assert.Nil(t, err)
}

func TestGetNonExisting(t *testing.T) {
	respValue = "test"
	status, loc, value, err := Get("localhost", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", GetNonexisting)
	assert.Equal(t, "404 NotFound", status)
	assert.Equal(t, "", loc)
	assert.Equal(t, "", value)
	assert.Nil(t, err)
}

func PostCreated(ip, contentType string, buffer io.Reader) (*http.Response, error) {
	b, _ := ioutil.ReadAll(buffer)

	body := Body{}
	_ = json.Unmarshal(b, &body)

	sha1 := sha1.Sum([]byte(body.Value))
	key := hex.EncodeToString(sha1[:])
	res := Response{"/objects/" + key, body.Value}

	s, _ := json.Marshal(res)
	r := ioutil.NopCloser(bytes.NewBuffer(s)) // r type is io.ReadCloser
	resp := http.Response{Status: "201 Created", StatusCode: 201, Proto: "", ProtoMajor: 1, ProtoMinor: 0, Header: nil, Body: r, ContentLength: 1024, TransferEncoding: []string{}, Uncompressed: true, Close: true, Trailer: nil, Request: nil, TLS: nil}

	return &resp, nil
}

func GetExisting(ip string) (*http.Response, error) {
	hash := strings.Split(ip, "/")[4]
	res := Response{"/objects/" + hash, respValue}

	s, _ := json.Marshal(res)
	r := ioutil.NopCloser(bytes.NewBuffer(s)) // r type is io.ReadCloser
	resp := http.Response{Status: "200 OK", StatusCode: 200, Proto: "", ProtoMajor: 1, ProtoMinor: 0, Header: nil, Body: r, ContentLength: 1024, TransferEncoding: []string{}, Uncompressed: true, Close: true, Trailer: nil, Request: nil, TLS: nil}

	return &resp, nil
}

func GetNonexisting(ip string) (*http.Response, error) {
	res := Response{}

	s, _ := json.Marshal(res)
	r := ioutil.NopCloser(bytes.NewBuffer(s)) // r type is io.ReadCloser
	resp := http.Response{Status: "404 NotFound", StatusCode: 404, Proto: "", ProtoMajor: 1, ProtoMinor: 0, Header: nil, Body: r, ContentLength: 1024, TransferEncoding: []string{}, Uncompressed: true, Close: true, Trailer: nil, Request: nil, TLS: nil}

	return &resp, nil
}
