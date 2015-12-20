package main

import "fmt"
import "github.com/amassanet/golatch1"
import "github.com/codegangsta/cli"
import "os"
import "log"

func getIdAndSecret(c *cli.Context) (id string, secret string) {
	if !c.GlobalIsSet("id") || !c.GlobalIsSet("secret") {
		log.Fatalf("--id and --secret requiered")
	}
	id = c.GlobalString("id")
	secret = c.GlobalString("secret")
	return id, secret
}

func main() {
	app := cli.NewApp()
	app.Name = "latchcli"
	app.Usage = "a latch client in go"
	app.Commands = []cli.Command{
		{
			Name:    "pair",
			Aliases: []string{"p"},
			Usage:   "pair <token> : pairs with supplied token. Returns the accountId",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					log.Fatalf("Missing token argument")
				}
				id, secret := getIdAndSecret(c)
				token := c.Args().First()
				la := golatch1.NewLatchApp(id, secret)
				if accountId, err := la.Pair(token); err == nil {
					fmt.Printf("%v", accountId)
				} else {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "unpair",
			Aliases: []string{"u"},
			Usage:   "unpair <accountId> : unpairs the account.",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					log.Fatalf("Missing accountId argument")
				}
				id, secret := getIdAndSecret(c)
				accountId := c.Args().First()
				la := golatch1.NewLatchApp(id, secret)
				if err := la.Unpair(accountId); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "status <accountId> : gets status of account",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					log.Fatalf("Missing accountId argument")
				}
				id, secret := getIdAndSecret(c)
				accountId := c.Args().First()
				la := golatch1.NewLatchApp(id, secret)
				if status, err := la.StatusIsOn(accountId); err == nil {
					if status {
						fmt.Printf("on")
					} else {
						fmt.Printf("off")
					}
				} else {
					log.Fatal(err)
				}
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "Application id",
		},
		cli.StringFlag{
			Name:  "secret",
			Usage: "Application secret",
		},
	}
	app.Run(os.Args)
}
