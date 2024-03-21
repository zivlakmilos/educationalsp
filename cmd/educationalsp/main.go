package main

import (
	"bufio"
	"log"
	"os"

	"github.com/zivlakmilos/educationalsp/pkg/rpc"
)

func main() {
	logger := getLogger("/tmp/educationalsp.txt")
	logger.Println("LSP Started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handlMessage(logger, msg)
	}
}

func handlMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad logger file")
	}

	return log.New(logfile, "[educationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
