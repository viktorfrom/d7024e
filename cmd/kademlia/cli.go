package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

var str = []string{""}

func info() {
  app.Name = "Kademlia CLI"
  app.Usage = "An example prototype CLI for Kademlia instructions"
  app.Author = "viktorfrom, markhakansson, 97gushan" 
  app.Version = "1.0.0"
}

func commands() {
  app.Commands = []cli.Command{
    {
      Name:    "exit",
      Aliases: []string{"e"},
      Usage:   "Terminates specified node",
      Action: func(c *cli.Context) { 
        value := "node value" // TODO: retrieve actual value
        content := append(str, value)
        m := strings.Join(content, " ")
        fmt.Println(m)
      },
    },
    {
      Name:    "get",
      Aliases: []string{"g"},
      Usage:   "Retrieves content of specified node",
      Action: func(c *cli.Context) { 
        value := "node value" // TODO: retrieve actual value
        content := append(str, value)
        m := strings.Join(content, " ")
        fmt.Println(m)
      },
    },
    {
      Name:    "put",
      Aliases: []string{"p"},
      Usage:   "Appends node and content to network",
      Action: func(c *cli.Context) { 
        value := "node value" // TODO: retrieve actual value
        content := append(str, value)
        m := strings.Join(content, " ")
        fmt.Println(m)
      },
    },
  }
}

func cli() {
  info()
  commands()

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}