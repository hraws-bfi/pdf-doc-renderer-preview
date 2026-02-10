package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// handleRenderHTML handles POST /render/html - renders a Go template with data
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

// handleSaveTemplate handles POST /templates/save - saves a template to disk
func handleSaveTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, SaveResponse{Error: "invalid json: " + err.Error()})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, SaveResponse{Error: "name is required"})
		return
	}
	if req.Template == "" {
		writeJSON(w, http.StatusBadRequest, SaveResponse{Error: "template is required"})
		return
	}

	safeName := sanitizeName(req.Name)
	if safeName == "" {
		writeJSON(w, http.StatusBadRequest, SaveResponse{Error: "invalid name"})
		return
	}

	version := getNextVersion(safeName)
	filename := fmt.Sprintf("%s-v%d.html", safeName, version)

	if err := saveFile(filename, []byte(req.Template)); err != nil {
		writeJSON(w, http.StatusInternalServerError, SaveResponse{Error: err.Error()})
		return
	}

	// Optionally save the JSON data alongside
	if req.Data != "" {
		dataFilename := fmt.Sprintf("%s-v%d.json", safeName, version)
		if err := saveFile(dataFilename, []byte(req.Data)); err != nil {
			log.Printf("warning: failed to save data file: %v", err)
		}
	}

	log.Printf("saved template: %s (version %d)", filename, version)
	writeJSON(w, http.StatusOK, SaveResponse{Filename: filename, Version: version})
}

// handleListTemplates handles GET /templates/list - lists all templates
func handleListTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries, err := os.ReadDir(config.TemplatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			writeJSON(w, http.StatusOK, ListResponse{Templates: []TemplateInfo{}})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ListResponse{Error: "failed to read templates directory"})
		return
	}

	var templates []TemplateInfo
	versionPattern := regexp.MustCompile(`^(.+)-v(\d+)\.html$`)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		matches := versionPattern.FindStringSubmatch(entry.Name())
		if matches != nil {
			version, _ := strconv.Atoi(matches[2])
			templates = append(templates, TemplateInfo{
				Name:     matches[1],
				Filename: entry.Name(),
				Version:  version,
			})
		} else {
			name := strings.TrimSuffix(entry.Name(), ".html")
			templates = append(templates, TemplateInfo{
				Name:     name,
				Filename: entry.Name(),
				Version:  0,
			})
		}
	}

	// Sort by name, then by version descending
	sort.Slice(templates, func(i, j int) bool {
		if templates[i].Name != templates[j].Name {
			return templates[i].Name < templates[j].Name
		}
		return templates[i].Version > templates[j].Version
	})

	writeJSON(w, http.StatusOK, ListResponse{Templates: templates})
}

// handleUploadDMS handles POST /templates/upload-dms - uploads a template to DMS
func handleUploadDMS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check DMS configuration
	if config.DMS.APIURL == "" || config.DMS.APISecret == "" {
		writeJSON(w, http.StatusInternalServerError, UploadDMSResponse{
			Error: "DMS not configured. Set DMS_API_URL and DMS_API_SECRET in .env",
		})
		return
	}

	var req UploadDMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, UploadDMSResponse{Error: "invalid json: " + err.Error()})
		return
	}

	// Validate required fields
	if req.Filename == "" {
		writeJSON(w, http.StatusBadRequest, UploadDMSResponse{Error: "filename is required"})
		return
	}
	if req.RefID == "" {
		writeJSON(w, http.StatusBadRequest, UploadDMSResponse{Error: "ref_id is required"})
		return
	}

	// Set defaults
	if req.IDType == "" {
		req.IDType = "template"
	}
	if req.DocumentType == "" {
		req.DocumentType = "html"
	}
	if req.SourceSystem == "" {
		req.SourceSystem = "LORA"
	}
	if req.DocumentSequence == "" {
		req.DocumentSequence = "1"
	}

	// Read the template file
	fileContent, err := readTemplateFile(req.Filename)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, UploadDMSResponse{
			Error: fmt.Sprintf("failed to read template file: %v", err),
		})
		return
	}

	// Create multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the file
	part, err := writer.CreateFormFile("file", req.Filename)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, UploadDMSResponse{Error: "failed to create form file"})
		return
	}
	if _, err := part.Write(fileContent); err != nil {
		writeJSON(w, http.StatusInternalServerError, UploadDMSResponse{Error: "failed to write file content"})
		return
	}

	// Add form fields
	writer.WriteField("ref_id", req.RefID)
	writer.WriteField("id_type", req.IDType)
	writer.WriteField("document_type", req.DocumentType)
	writer.WriteField("source_system", req.SourceSystem)
	writer.WriteField("document_sequence", req.DocumentSequence)

	if err := writer.Close(); err != nil {
		writeJSON(w, http.StatusInternalServerError, UploadDMSResponse{Error: "failed to close multipart writer"})
		return
	}

	// Create the request to DMS
	dmsReq, err := http.NewRequest(http.MethodPost, config.DMS.APIURL, &body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, UploadDMSResponse{Error: "failed to create DMS request"})
		return
	}

	dmsReq.Header.Set("Content-Type", writer.FormDataContentType())
	dmsReq.Header.Set("api-secret", config.DMS.APISecret)

	// Send the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(dmsReq)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, UploadDMSResponse{
			Error: fmt.Sprintf("DMS request failed: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, _ := io.ReadAll(resp.Body)
	respString := string(respBody)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("uploaded template to DMS: %s (ref_id: %s)", req.Filename, req.RefID)
		writeJSON(w, http.StatusOK, UploadDMSResponse{
			Success:  true,
			Message:  fmt.Sprintf("Template %s uploaded successfully", req.Filename),
			Response: respString,
		})
	} else {
		writeJSON(w, http.StatusBadGateway, UploadDMSResponse{
			Error:    fmt.Sprintf("DMS returned status %d", resp.StatusCode),
			Response: respString,
		})
	}
}

// handleRenderPDF handles POST /render/pdf - renders a Go template to PDF
func handleRenderPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PDFRequest
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
		req.Data = map[string]interface{}{}
	}

	// Parse and execute Go template
	tmpl, err := template.New("pdf").Option("missingkey=default").Parse(req.Template)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "template parse error: " + err.Error()})
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req.Data); err != nil {
		writeJSON(w, http.StatusBadRequest, RenderResponse{Error: "template execute error: " + err.Error()})
		return
	}

	htmlContent := buf.String()

	// Generate PDF using chromedp
	pdfBytes, err := generatePDF(htmlContent)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, RenderResponse{Error: "pdf generation error: " + err.Error()})
		return
	}

	// Set filename for download
	filename := "document.pdf"
	if req.Filename != "" {
		filename = sanitizeName(req.Filename) + ".pdf"
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.Write(pdfBytes)
}
