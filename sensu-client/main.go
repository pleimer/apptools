package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	interval := flag.Int("interval", 15, "interval at which to execute script in seconds")
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

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

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

	ticker := time.NewTicker(time.Second * time.Duration(*interval))

	for {
		select {
		case <-c:
			closeContext, closeCancel := context.WithTimeout(ctx, time.Second*5)
			defer closeCancel()
			sender.close(closeContext)
			ticker.Stop()
			cancel()
			fmt.Println("keyboard interrupt: canceled")
		case <-ctx.Done():
			fmt.Println("sensu-client exited gracefully")
			return
		case <-ticker.C:
			scriptContext, scriptCancel := context.WithTimeout(ctx, time.Second*30)

			res, err := runner.Run(scriptContext)
			scriptCancel()
			if err != nil {
				log.Fatal(err)
			}

			sendContext, sendCancel := context.WithTimeout(ctx, time.Second*5)
			err = sender.send(sendContext, res)
			sendCancel()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
