package main

import (
	"bufio"
	"encoding/json"
	"io"
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
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("error: %s\n", err)
			continue
		}

		handlMessage(logger, writer, state, method, contents)
	}
}

func handlMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
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
		writeResponse(writer, msg)

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
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error textDocument/hover: %s\n", err)
			return
		}

		logger.Printf("hover: %s at pos row(%d), col(%d)\n", request.Params.TextDocument.URI, request.Params.Position.Line, request.Params.Position.Character)

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error textDocument/definition: %s\n", err)
			return
		}

		logger.Printf("definition: %s at pos row(%d), col(%d)\n", request.Params.TextDocument.URI, request.Params.Position.Line, request.Params.Position.Character)

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad logger file")
	}

	return log.New(logfile, "[educationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
