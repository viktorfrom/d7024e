package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	errNoArg       string = "No argument!"
	errWrongArg    string = "Wrong amount of arguments!"
	errInvalidCmd  string = "Invalid command!"
	errNoFileFound string = "Could not find or open file: "
)

var (
	osExit   = os.Exit
	logFatal = log.Fatal
	helpFile = Prompt()
)

type Body struct {
	Value string `json:"value"`
}

type Response struct {
	Location string `json:"location"`
	Value    string `json:"value"`
}

var in *os.File = os.Stdin
var out io.Writer = os.Stdout

func main() {
	fmt.Fprintln(out, "Starting Client CLI...")
	reader := bufio.NewReader(in)

	for {
		input, _ := reader.ReadString('\n')
		trimInput := strings.TrimSpace(input)

		if trimInput == "\n" || trimInput == "" {
			continue

		} else {
			commands := strings.Fields(trimInput)

			Commands(out, commands)
		}
	}
}

// Commands handles the commands of the CLI. `output` is the io.Writer to output data to.
// `node` is the Kademlia node this CLI runs for. `commands` a list of program commands.
func Commands(output io.Writer, commands []string) {
	switch commands[0] {
	case "put":
		if len(commands) == 3 {
			Put(GetAPIUrl(commands[1]), commands[2])
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "p":
		if len(commands) == 3 {
			Put(GetAPIUrl(commands[1]), commands[2])
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "get":
		if len(commands) == 3 {
			Get(GetAPIUrl(commands[1]), commands[2])
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "g":
		if len(commands) == 3 {
			Get(GetAPIUrl(commands[1]), commands[2])
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "exit":
		Exit()
	case "e":
		Exit()
	case "help":
		Help(output)
	case "h":
		Help(output)
	default:
		fmt.Fprintln(output, errInvalidCmd)
	}
}

func Put(url, value string) (string, string, string, error) {
	body, err := json.Marshal(Body{value})
	if err != nil {
		return "500", "", "", err
	}
	resp, err := http.Post(url+"/objects", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "500", "", "", err
	}

	body, err = ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	data := Response{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
	} else {
		fmt.Fprintln(out, "Status:", resp.Status)
		fmt.Fprintln(out, "Location:", data.Location)
		fmt.Fprintln(out, "Value:", data.Value)
	}
	return resp.Status, data.Location, data.Value, err
}

func Get(url, hash string) (string, string, string, error) {

	resp, err := http.Get(url + "/objects/" + hash)
	if err != nil {
		return "500", "", "", err
	}

	if resp.StatusCode != 200 {
		return strconv.Itoa(resp.StatusCode), "", "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll((*resp).Body)
	if err != nil {
		return "500", "", "", err
	}
	defer resp.Body.Close()

	data := Response{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
	} else {
		fmt.Fprintln(out, "Status:", resp.Status)
		fmt.Fprintln(out, "Location:", data.Location)
		fmt.Fprintln(out, "Value:", data.Value)
	}
	return resp.Status, data.Location, data.Value, err
}

func GetAPIUrl(ip string) string {
	return "http://" + ip + ":3000"
}

func Exit() {
	osExit(3)
}

func Help(output io.Writer) {
	text := Prompt()
	fmt.Fprintln(output, text)
}
