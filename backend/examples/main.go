package main

import (
	"log"
	"os"

	"ragify-backend/examples/chunking"
	"ragify-backend/examples/vectorstore"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run examples/main.go [chunking|vectorstore]")
	}

	example := os.Args[1]

	switch example {
	case "chunking":
		log.Println("Running chunking example...")
		chunking.RunChunkingExample()
	case "vectorstore":
		log.Println("Running vector store example...")
		vectorstore.RunVectorStoreExample()
	default:
		log.Fatalf("Unknown example: %s. Use 'chunking' or 'vectorstore'", example)
	}
}
