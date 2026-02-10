package main

import (
	"context"
	"encoding/json"
	"fmt"
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

// generatePDF converts HTML content to PDF using headless Chrome
func generatePDF(htmlContent string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for PDF generation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte

	// Navigate to a data URL with the HTML content
	// We need to use a data URL since we're rendering dynamic content
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
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
