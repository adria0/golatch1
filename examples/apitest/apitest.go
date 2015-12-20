package main

import "fmt"
import "github.com/amassanet/golatch1"
import "os"
import "log"

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Needs parameters <id> <secret> <token>")
	}
	id, secret, token := os.Args[1], os.Args[2], os.Args[3]
	la := golatch1.NewLatchApp(id, secret)

	var accountId string
	var err error

	if accountId, err = la.Pair(token); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Paired ok. Account id is %v\n", accountId)

	var statusIsOn bool

	if statusIsOn, err = la.StatusIsOn(accountId); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status is on %v\n", statusIsOn)

	if err := la.Unpair(accountId); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Unpaired ok.\n")
}
