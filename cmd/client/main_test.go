package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var respValue string
var respKey string
var respErr error

var (
	mux    *http.ServeMux
	server *httptest.Server
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

func TestGetAPIUrl(t *testing.T) {
	assert.Equal(t, "http://10.0.8.4:3000", GetAPIUrl("10.0.8.4"))
}

func TestPut201(t *testing.T) {
	server, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		body := Body{}
		json.Unmarshal(b, &body)

		sha1 := sha1.Sum([]byte(body.Value))
		key := hex.EncodeToString(sha1[:])
		res := Response{"/objects/" + key, body.Value}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	})

	status, loc, value, err := Put(server.URL, "test")
	assert.Equal(t, "201 Created", status)
	assert.Equal(t, "/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", loc)
	assert.Equal(t, "test", value)
	assert.Nil(t, err)
}

func TestPut500(t *testing.T) {
	server, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	status, loc, value, err := Put(server.URL, "test")
	assert.Equal(t, "500 Internal Server Error", status)
	assert.Equal(t, "", loc)
	assert.Equal(t, "", value)
	assert.Error(t, err)
}

func TestGet200(t *testing.T) {
	server, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", func(w http.ResponseWriter, r *http.Request) {
		hash := strings.Split(r.URL.Path, "/")[2]
		res := Response{"/objects/" + hash, "test"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})

	status, loc, value, err := Get(server.URL, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")
	assert.Equal(t, "200 OK", status)
	assert.Equal(t, "/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", loc)
	assert.Equal(t, "test", value)
	assert.Nil(t, err)
}

func TestGet400(t *testing.T) {
	server, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	status, loc, value, err := Get(server.URL, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd")
	assert.Equal(t, "400", status)
	assert.Equal(t, "", loc)
	assert.Equal(t, "", value)
	assert.Error(t, err)
}
func TestGet404(t *testing.T) {
	server, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/objects/a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	status, loc, value, err := Get(server.URL, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")
	assert.Equal(t, "404", status)
	assert.Equal(t, "", loc)
	assert.Equal(t, "", value)
	assert.Error(t, err)
}

func setupTestServer() (*httptest.Server, func()) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	teardown := func() {
		server.Close()
	}
	return server, teardown
}
