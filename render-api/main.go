package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Template rendering
	mux.HandleFunc("/render/html", withCORS(handleRenderHTML))
	mux.HandleFunc("/render/pdf", withCORS(handleRenderPDF))

	// Template management
	mux.HandleFunc("/templates/save", withCORS(handleSaveTemplate))
	mux.HandleFunc("/templates/list", withCORS(handleListTemplates))
	mux.HandleFunc("/templates/upload-dms", withCORS(handleUploadDMS))

	log.Printf("render-api listening on %s", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, mux); err != nil {
		log.Fatal(err)
	}
}
