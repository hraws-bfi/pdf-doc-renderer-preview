package main

// RenderRequest represents a template rendering request
type RenderRequest struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

// RenderResponse represents a template rendering response
type RenderResponse struct {
	HTML  string `json:"html,omitempty"`
	Error string `json:"error,omitempty"`
}

// SaveRequest represents a template save request
type SaveRequest struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Data     string `json:"data,omitempty"` // optional JSON data to save alongside
}

// SaveResponse represents a template save response
type SaveResponse struct {
	Filename string `json:"filename,omitempty"`
	Version  int    `json:"version,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ListResponse represents a template list response
type ListResponse struct {
	Templates []TemplateInfo `json:"templates,omitempty"`
	Error     string         `json:"error,omitempty"`
}

// TemplateInfo represents metadata about a template file
type TemplateInfo struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  int    `json:"version"`
}

// UploadDMSRequest represents a DMS upload request
type UploadDMSRequest struct {
	Filename         string `json:"filename"`          // Template filename to upload
	RefID            string `json:"ref_id"`            // UUID reference
	IDType           string `json:"id_type"`           // e.g., "test"
	DocumentType     string `json:"document_type"`     // e.g., "html"
	SourceSystem     string `json:"source_system"`     // e.g., "LORA"
	DocumentSequence string `json:"document_sequence"` // e.g., "1"
}

// PDFRequest represents a PDF generation request
type PDFRequest struct {
	Template      string                 `json:"template"`
	Data          map[string]interface{} `json:"data"`
	Filename      string                 `json:"filename,omitempty"`        // optional filename for download
	WaitAfterLoad int                    `json:"wait_after_load,omitempty"` // milliseconds to wait for images to load (default: 500)
}

// UploadDMSResponse represents a DMS upload response
type UploadDMSResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
	Response string `json:"response,omitempty"` // Raw DMS response
}
