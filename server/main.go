package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
)

// Deklarojme nje map per ta perdorur si DB per te mbajtur te dhenat
var fakeDatabase map[string]string

// Deklarojme nje strukture per te dhenat e vendosura nga perdoruesi dhe proceduren nqs eshte get ose put
type HashMap struct {
	Key, Value string
	Procedure  int
}

// Deklarojme strukturen e marjes se te dhenave nqs komanda eshte get
type OutPut struct {
	Value string
}

// Deklarojme tipin FileStream per te lexuar dhe shkruar ne file
type FileStream int

func (f *FileStream) ReadFile() map[string]string {
	jsonFile, err := os.Open("db.txt")
	checkError(err)
	defer jsonFile.Close()

	content, err := ioutil.ReadAll(jsonFile)
	checkError(err)
	json.Unmarshal(content, &fakeDatabase)

	return fakeDatabase
}

func (f *FileStream) WriteFile() error {
	content, _ := json.Marshal(&fakeDatabase)
	err := ioutil.WriteFile("db.txt", content, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}

// Deklarojme tipin Orm per te menaxhuar "transaksionin"
type Orm int

func (o *Orm) Transaction(args *HashMap, out *OutPut) error {
	fileStream := new(FileStream)

	// Nqs procedura eshte 1 ath e shtojme vleren ne DB perndryshe e kthejme nga DB
	if args.Procedure == 1 {
		fakeDatabase[args.Key] = args.Value
		fileStream.WriteFile()
	} else if args.Procedure == 2 {
		// Kontrollojme nqs vlera ekziston apo jo ne DB
		value, exists := fakeDatabase[args.Key]
		if !exists {
			return fmt.Errorf("nuk ekziston rekord me kete key %v!", args.Key)
		}
		out.Value = value
	}

	// Afishojme "loget" e DB
	fmt.Println(args)

	return nil
}

func main() {
	orm := new(Orm)
	fileStream := new(FileStream)

	err := rpc.Register(orm)
	if err != nil {
		return
	}

	fakeDatabase = fileStream.ReadFile()

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

