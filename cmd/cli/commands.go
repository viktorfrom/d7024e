package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

const (
	errNoArg       string = "No argument!"
	errInvalidCmd  string = "Invalid command!"
	errNoFileFound string = "Could not find or open file: "
)

var (
	osExit   = os.Exit
	logFatal = log.Fatal
	helpFile = "prompt.txt"
)

// Commands handles the commands of the CLI. `output` is the io.Writer to output data to.
// `node` is the Kademlia node this CLI runs for. `commands` a list of program commands.
func Commands(output io.Writer, node *kademlia.Node, commands []string) {

	switch commands[0] {
	case "put":
		if len(commands) == 2 {
			Put(*node, commands[1])
		} else {
			fmt.Fprintln(output, errNoArg)
		}
	case "p":
		if len(commands) == 2 {
			Put(*node, commands[1])
		} else {
			fmt.Fprintln(output, errNoArg)
		}
	case "ping":
		if len(commands) == 2 {
			Ping(*node, commands[1])
		} else {
			fmt.Fprintln(output, errNoArg)
		}
	case "get":
		if len(commands) == 2 {
			Get(*node, commands[1])
		} else {
			fmt.Fprintln(output, errNoArg)
		}
	case "g":
		if len(commands) == 2 {
			Get(*node, commands[1])
		} else {
			fmt.Fprintln(output, errNoArg)
		}
	case "t":
		//c := kademlia.NewContact(kademlia.NewRandomNodeID(), "10.0.8.9")
		c := kademlia.NewContact(kademlia.NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.3:8080")
		c.CalcDistance(node.RT.GetMeID())
		fmt.Fprintln(output, node.NodeLookup(c.ID))
	case "info":
		fmt.Println("ID: ", node.RT.GetMeID())
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

func Put(node kademlia.Node, input string) {
	hash := node.StoreValue(input)
	println("Hash = ", hash)
}

func Ping(node kademlia.Node, input string) {
	// node.Ping()
}

func Get(node kademlia.Node, hash string) {
	value, err := node.FindValue(hash)

	if err != nil {
		println(err.Error())
	} else {
		println("value = ", value)
	}
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
