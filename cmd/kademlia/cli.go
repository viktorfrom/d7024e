import "github.com/urfave/cli" 

func cli() {
app := &cli.App{
			Name: "greet",
			Usage: "say a greeting",
			Action: func(c *cli.Context) error {
				fmt.Println("Greetings")
				return nil
			},
		}

app.Run(os.Args)
}