package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type RenderRequest struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

type RenderResponse struct {
	HTML  string `json:"html,omitempty"`
	Error string `json:"error,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/render/html", withCORS(handleRenderHTML))

	addr := ":8080"
	log.Printf("render-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleRenderHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RenderRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "invalid json: " + err.Error()})
		return
	}

	if req.Template == "" {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "template is required"})
		return
	}
	if req.Data == nil {
		// allow empty object, but not null (keeps behavior predictable)
		req.Data = map[string]interface{}{}
	}

	tmpl, err := template.New("preview").Option("missingkey=default").Parse(req.Template)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "template parse error: " + err.Error()})
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req.Data); err != nil {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "template execute error: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, RenderResponse{HTML: buf.String()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Minimal dev CORS. Tighten if you ever expose beyond localhost.
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Allow local dev origins. Add more ports if needed.
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		next(w, r)
	}
}
