package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("%s [OPTION] \n", os.Args[0])
		flag.PrintDefaults()
	}

	address := flag.String("address", "amqp://0.0.0.0:5672", "amqp1 url and address")
	target := flag.String("target", "", "target address e.g. /example-queue-name")
	fileScript := flag.String("file", "", "filepath to script")
	command := flag.String("command", "", "command to run")
	flag.Parse()

	if *fileScript == "" && *command == "" {
		fmt.Println("filepath or command must be specified")
		flag.Usage()
		return
	}

	if *target == "" {
		fmt.Println("target option must be specified")
		flag.Usage()
		return
	}

	runner := Runner{}
	runner.LoadScriptFromCommand([]byte(*command))
	if *fileScript != "" {
		runner.LoadScriptFromFile(*fileScript)
	}

	sender, err := newSender(*address, *target)
	if err != nil {
		log.Panic(err)
		return
	}

	ctx := context.Background()
	scriptContext, scriptCancel := context.WithTimeout(ctx, time.Second*30)

	res, err := runner.Run(scriptContext)
	if err != nil {
		log.Fatal(err)
		return
	}
	scriptCancel()

	sendContext, sendCancel := context.WithTimeout(ctx, time.Second*5)
	defer sendCancel()

	err = sender.send(sendContext, res)
	if err != nil {
		log.Fatal(err)
		return
	}

}
