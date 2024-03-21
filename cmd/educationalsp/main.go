package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/zivlakmilos/educationalsp/pkg/analysis"
	"github.com/zivlakmilos/educationalsp/pkg/lsp"
	"github.com/zivlakmilos/educationalsp/pkg/rpc"
)

func main() {
	logger := getLogger("/tmp/educationalsp.txt")
	logger.Println("LSP Started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("error: %s\n", err)
			continue
		}

		handlMessage(logger, state, method, contents)
	}
}

func handlMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Printf("received msg with method: %s\n", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error initialize: %s\n", err)
			return
		}

		logger.Printf("connected to: %s %s\n",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Printf("sent reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error textDocument/didOpen: %s\n", err)
			return
		}

		logger.Printf("opened: %s\n", request.Params.TextDocument.URI)

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error textDocument/didChange: %s\n", err)
			return
		}

		logger.Printf("changed: %s\n", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			state.OpenDocument(request.Params.TextDocument.URI, change.Text)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad logger file")
	}

	return log.New(logfile, "[educationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
