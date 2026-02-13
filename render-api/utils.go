package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// sanitizeName converts a name to a safe filename
func sanitizeName(name string) string {
	// Convert to lowercase and replace spaces with dashes
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any characters that aren't alphanumeric, dash, or underscore
	reg := regexp.MustCompile(`[^a-z0-9\-_]`)
	name = reg.ReplaceAllString(name, "")

	// Remove consecutive dashes
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	return strings.Trim(name, "-")
}

// getNextVersion finds the next version number for a template
func getNextVersion(baseName string) int {
	entries, err := os.ReadDir(config.TemplatesDir)
	if err != nil {
		return 1
	}

	maxVersion := 0
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s-v(\d+)\.html$`, regexp.QuoteMeta(baseName)))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := pattern.FindStringSubmatch(entry.Name())
		if matches != nil {
			v, _ := strconv.Atoi(matches[1])
			if v > maxVersion {
				maxVersion = v
			}
		}
	}

	return maxVersion + 1
}

// saveFile saves content to a file in the templates directory
func saveFile(filename string, content []byte) error {
	if err := os.MkdirAll(config.TemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	filePath := filepath.Join(config.TemplatesDir, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// readTemplateFile reads a template file from the templates directory
func readTemplateFile(filename string) ([]byte, error) {
	filePath := filepath.Join(config.TemplatesDir, filename)
	return os.ReadFile(filePath)
}

// isURLString checks if a string looks like a URL that should be marked as safe
func isURLString(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "data:") ||
		strings.HasPrefix(s, "file://")
}

// sanitizeDataForTemplate recursively converts URL-like strings to template.URL
// so they render correctly without needing explicit safeURL in templates
func sanitizeDataForTemplate(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = sanitizeValue(v)
	}
	return result
}

func sanitizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		if isURLString(val) {
			return template.URL(val)
		}
		return val
	case map[string]interface{}:
		return sanitizeDataForTemplate(val)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = sanitizeValue(item)
		}
		return result
	default:
		return val
	}
}

// generatePDF converts HTML content to PDF using headless Chrome
// waitMs specifies how long to wait (in milliseconds) for images to load after page ready
func generatePDF(htmlContent string, waitMs int) ([]byte, error) {
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-gpu", true),
	)
	if config.ChromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(config.ChromePath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for PDF generation
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Default wait time if not specified
	if waitMs <= 0 {
		waitMs = 500
	}

	// Write HTML to a temporary file so Chrome can load it with external resources
	tmpFile, err := os.CreateTemp("", "pdf-render-*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	var pdfBuf []byte

	// Navigate to the temp file URL so Chrome can fetch external resources (images)
	fileURL := "file://" + tmpPath

	if err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(time.Duration(waitMs)*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp error: %w", err)
	}

	return pdfBuf, nil
}
