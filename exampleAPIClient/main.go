package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	helpFile = "prompt.txt"
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
	fmt.Fprintln(out, "Starting CLI...")
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
			status, location, value, err := Put(commands[1], commands[2], http.Post)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(status)
				fmt.Println(location)
				fmt.Println(value)
			}
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "p":
		if len(commands) == 3 {
			status, location, value, err := Put(commands[1], commands[2], http.Post)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(status)
				fmt.Println(location)
				fmt.Println(value)
			}
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "get":
		if len(commands) == 3 {
			status, location, value, err := Get(commands[1], commands[2], http.Get)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(status)
				fmt.Println(location)
				fmt.Println(value)
			}
		} else {
			fmt.Fprintln(output, errWrongArg)
		}
	case "g":
		if len(commands) == 3 {
			status, location, value, err := Get(commands[1], commands[2], http.Get)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(status)
				fmt.Println(location)
				fmt.Println(value)
			}
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

func Put(ip, value string, post func(ip, contentType string, buffer io.Reader) (*http.Response, error)) (string, string, string, error) {
	body, err := json.Marshal(Body{value})
	if err != nil {
		return "500", "", "", err
	}
	resp, err := post("http://"+ip+":3000/objects", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "500", "", "", err
	}

	body, err = ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	data := Response{}
	err = json.Unmarshal(body, &data)
	return resp.Status, data.Location, data.Value, err
}

func Get(ip, hash string, get func(url string) (*http.Response, error)) (string, string, string, error) {
	fmt.Println(ip)
	resp, err := get("http://" + ip + ":3000/objects/" + hash)
	if err != nil {
		return "500", "", "", err
	}

	body, err := ioutil.ReadAll((*resp).Body)
	if err != nil {
		return "500", "", "", err
	}
	defer resp.Body.Close()

	data := Response{}

	err = json.Unmarshal(body, &data)
	return resp.Status, data.Location, data.Value, err
}

func Exit() {
	osExit(3)
}

func Help(output io.Writer) {
	content, err := ioutil.ReadFile(helpFile)
	if err != nil {
		logFatal(errNoFileFound + helpFile)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	fmt.Fprintln(output, text)
}
