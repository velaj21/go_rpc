package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type HashMap struct {
	Key, Value string
	Procedure  int
}

type OutPut struct {
	Value string
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	service := os.Args[1]
	client, err := rpc.Dial("tcp", service)

	if err != nil {
		log.Fatal("dialing:", err)
	}

	outPut := OutPut{""}

	for {
		// Lexojme komandat dhe ne baze te procedures kryhen veprimet perkatese
		fmt.Println("Vendosni komandat: ")
		scanner.Scan()
		line := scanner.Text()
		cmds := strings.Split(line, " ")

		if cmds[0] == "exit" {
			fmt.Println("Exiting program...")
			return
		} else if cmds[0] == "put" {
			hashMap := HashMap{Key: cmds[1], Value: cmds[2], Procedure: 1}
			err = client.Call("Orm.Transaction", hashMap, &outPut)
			if err != nil {
				log.Fatal("ORM error:", err)
			}
		} else if cmds[0] == "get" {
			hashMap := HashMap{Key: cmds[1], Procedure: 2}
			err = client.Call("Orm.Transaction", hashMap, &outPut)
			if err != nil {
				log.Println("ORM error:", err)
			}
			fmt.Println(outPut)
		} else {
			fmt.Println("Komanda nuk njihet ose ka gabime ne sintakse!")
		}
	}

}
