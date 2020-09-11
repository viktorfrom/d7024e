package main

import (
	"os"
	"time"

	"github.com/mvisonneau/gitlab-ci-pipelines-exporter/cmd"
	"github.com/urfave/cli"
)

// Run handles the instanciation of the CLI application
func main() {
	Run("ads")

}

func Run(version string) {
	NewApp(version, time.Now()).Run(os.Args)
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "gitlab-ci-pipelines-exporter"
	app.Version = version
	app.Usage = "Export metrics about GitLab CI pipelines statuses"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "enable-pprof",
			EnvVar: "GCPE_ENABLE_PPROF",
			Usage:  "Enable profiling endpoints at /debug/pprof",
		},
		cli.StringFlag{
			Name:   "gitlab-token",
			EnvVar: "GCPE_GITLAB_TOKEN",
			Usage:  "GitLab access `token`. Can be use to override the gitlab token in config file",
		},
		cli.StringFlag{
			Name:   "listen-address, l",
			EnvVar: "GCPE_LISTEN_ADDRESS",
			Usage:  "listen-address `address:port`",
			Value:  ":8080",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "GCPE_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "GCPE_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
	}

	app.Action = cmd.ExecWrapper(cmd.Run)

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}


// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/urfave/cli"
// )

// var app = cli.NewApp()

// var str = []string{""}

// func info() {
// 	app.Name = "Kademlia CLI"
// 	app.Usage = "An example prototype CLI for Kademlia instructions"
// 	app.Author = "viktorfrom, markhakansson, 97gushan"
// 	app.Version = "1.0.0"
// }

// func commands() {
// 	app.Commands = []cli.Command{
// 		{
// 			Name:    "exit",
// 			Aliases: []string{"e"},
// 			Usage:   "Terminates specified node",
// 			Action: func(c *cli.Context) {
// 				value := "node value" // TODO: retrieve actual value
// 				content := append(str, value)
// 				m := strings.Join(content, " ")
// 				fmt.Println(m)
// 			},
// 		},
// 		{
// 			Name:    "get",
// 			Aliases: []string{"g"},
// 			Usage:   "Retrieves content of specified node",
// 			Action: func(c *cli.Context) {
// 				value := "node value" // TODO: retrieve actual value
// 				content := append(str, value)
// 				m := strings.Join(content, " ")
// 				fmt.Println(m)
// 			},
// 		},
// 		{
// 			Name:    "put",
// 			Aliases: []string{"p"},
// 			Usage:   "Appends node and content to network",
// 			Action: func(c *cli.Context) {
// 				value := "node value" // TODO: retrieve actual value
// 				content := append(str, value)
// 				m := strings.Join(content, " ")
// 				fmt.Println(m)
// 			},
// 		},
// 	}
// }

// // StartCli TODO
// func StartCli() {
// 	info()
// 	commands()

// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
