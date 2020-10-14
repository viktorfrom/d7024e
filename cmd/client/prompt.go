package main

var help = `
NAME:
   Kademlia CLI - An example prototype CLI for Kademlia instructions
USAGE:
   cli [global options] command [command options] [arguments...]
VERSION:
   1.0.0
AUTHOR:
   viktorfrom, markhakansson, 97gushan
COMMANDS:
   exit, e      Terminates specified node
   get, g       Retrieves content of specified node
   put, p       Appends node and content to network
   help, h      Show help
   version, v   Print the version
`

func Prompt() string {
	return help
}
