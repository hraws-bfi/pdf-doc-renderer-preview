import { EditorView, basicSetup } from 'codemirror';
import { html } from '@codemirror/lang-html';
import { json } from '@codemirror/lang-json';
import { oneDark } from '@codemirror/theme-one-dark';

// =============================================================================
// DOM Elements
// =============================================================================

const elements = {
  // Controls
  renderBtn: document.querySelector("#render"),
  autoCheckbox: document.querySelector("#auto"),
  status: document.querySelector("#status"),
  saveStatus: document.querySelector("#saveStatus"),
  
  // Editors
  tplEditorContainer: document.querySelector("#tplEditor"),
  varsEditorContainer: document.querySelector("#varsEditor"),
  formatTplBtn: document.querySelector("#formatTpl"),
  formatVarsBtn: document.querySelector("#formatVars"),
  
  // Save
  templateNameInput: document.querySelector("#templateName"),
  saveBtn: document.querySelector("#save"),
  
  // Preview
  preview: document.querySelector("#preview"),
  downloadPdfBtn: document.querySelector("#downloadPdf"),
  
  // API Banner
  apiBanner: document.querySelector("#apiBanner"),
  wrapEl: document.querySelector(".wrap"),
  
  // DMS Upload Modal
  uploadModal: document.querySelector("#uploadModal"),
  openUploadModalBtn: document.querySelector("#openUploadModal"),
  closeUploadModalBtn: document.querySelector("#closeUploadModal"),
  uploadToDMSBtn: document.querySelector("#uploadToDMS"),
  dmsFilenameSelect: document.querySelector("#dmsFilename"),
  dmsRefIdInput: document.querySelector("#dmsRefId"),
  dmsIdTypeInput: document.querySelector("#dmsIdType"),
  dmsDocTypeInput: document.querySelector("#dmsDocType"),
  dmsSourceSystemInput: document.querySelector("#dmsSourceSystem"),
  dmsDocSequenceInput: document.querySelector("#dmsDocSequence"),
};

// =============================================================================
// Configuration
// =============================================================================

const API_BASE_URL = "http://localhost:8080";
const API_HEALTH_CHECK_INTERVAL = 5000;
const DEBOUNCE_DELAY = 300;
const STATUS_MESSAGE_DURATION = 5000;

// =============================================================================
// Default Content
// =============================================================================

const defaultTemplate = `<html>
<body style="font-family: Arial">
  <h1>Contract {{.contract_number}}</h1>
  <p>Name: {{.customer.name}}</p>
</body>
</html>`;

const defaultVariables = JSON.stringify({
  contract_number: "FIN-2024-001",
  customer: { name: "John Doe" }
}, null, 2);

// =============================================================================
// State
// =============================================================================

let apiOnline = true;
let debounceTimer = null;

// =============================================================================
// Editors
// =============================================================================

const tplEditor = new EditorView({
  doc: defaultTemplate,
  extensions: [
    basicSetup,
    html(),
    oneDark,
    EditorView.lineWrapping,
    EditorView.updateListener.of((update) => {
      if (update.docChanged && elements.autoCheckbox.checked) {
        debounceRender();
      }
    })
  ],
  parent: elements.tplEditorContainer
});

const varsEditor = new EditorView({
  doc: defaultVariables,
  extensions: [
    basicSetup,
    json(),
    oneDark,
    EditorView.lineWrapping,
    EditorView.updateListener.of((update) => {
      if (update.docChanged && elements.autoCheckbox.checked) {
        debounceRender();
      }
    })
  ],
  parent: elements.varsEditorContainer
});

// =============================================================================
// Utility Functions
// =============================================================================

function escapeHtml(str) {
  return String(str)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function setEditorContent(editor, content) {
  editor.dispatch({
    changes: { from: 0, to: editor.state.doc.length, insert: content }
  });
}

function showStatus(message, isError = true) {
  elements.status.textContent = isError ? message : "";
}

function showSaveStatus(message, duration = STATUS_MESSAGE_DURATION) {
  elements.saveStatus.textContent = message;
  if (duration > 0) {
    setTimeout(() => { elements.saveStatus.textContent = ""; }, duration);
  }
}

function parseJsonData(jsonStr) {
  try {
    return { data: JSON.parse(jsonStr || "{}"), error: null };
  } catch (e) {
    return { data: null, error: e.message };
  }
}

function formatHtml(htmlStr) {
  let formatted = '';
  let indent = 0;
  const tab = '  ';
  
  const tokens = htmlStr.replace(/>\s*</g, '>\n<').split('\n');
  
  tokens.forEach(token => {
    token = token.trim();
    if (!token) return;
    
    if (token.match(/^<\/\w/)) {
      indent = Math.max(0, indent - 1);
    }
    
    formatted += tab.repeat(indent) + token + '\n';
    
    if (token.match(/^<\w[^>]*[^\/]>$/) && !token.match(/^<(br|hr|img|input|meta|link|area|base|col|embed|param|source|track|wbr)/i)) {
      indent++;
    }
  });
  
  return formatted.trim();
}

// =============================================================================
// API Functions
// =============================================================================

async function apiRequest(endpoint, options = {}) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: { "Content-Type": "application/json" },
    ...options
  });
  return response;
}

async function checkApiHealth() {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 3000);
    
    await fetch(`${API_BASE_URL}/render/html`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ template: "", data: {} }),
      signal: controller.signal
    });
    clearTimeout(timeoutId);
    
    setApiStatus(true);
  } catch (e) {
    setApiStatus(false);
  }
}

function setApiStatus(online) {
  if (apiOnline === online) return;
  apiOnline = online;
  
  if (online) {
    elements.apiBanner.classList.remove("show");
    elements.wrapEl.classList.remove("api-offline");
    doRender();
  } else {
    elements.apiBanner.classList.add("show");
    elements.wrapEl.classList.add("api-offline");
    elements.preview.srcdoc = `<div style="padding:20px;color:#666;font-family:system-ui">API server is offline</div>`;
  }
}

// =============================================================================
// Render Functions
// =============================================================================

function debounceRender() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(doRender, DEBOUNCE_DELAY);
}

async function doRender() {
  showStatus("");
  
  const tplValue = tplEditor.state.doc.toString();
  const varsValue = varsEditor.state.doc.toString();
  
  const { data, error } = parseJsonData(varsValue);
  if (error) {
    showStatus("JSON error: " + error);
    return;
  }

  try {
    const res = await apiRequest("/render/html", {
      method: "POST",
      body: JSON.stringify({ template: tplValue, data })
    });

    const payload = await res.json().catch(() => ({}));
    if (!res.ok) {
      const errorMsg = payload.error || `Render failed (${res.status})`;
      showStatus(errorMsg);
      elements.preview.srcdoc = `<pre style="color:#b00020;white-space:pre-wrap">${escapeHtml(errorMsg)}</pre>`;
      return;
    }

    elements.preview.srcdoc = payload.html ?? "";
  } catch (e) {
    showStatus("Render error: " + e.message);
  }
}

// =============================================================================
// PDF Download
// =============================================================================

async function downloadPdf() {
  const tplValue = tplEditor.state.doc.toString();
  const varsValue = varsEditor.state.doc.toString();

  if (!tplValue) {
    showStatus("Template is empty");
    return;
  }

  const { data, error } = parseJsonData(varsValue);
  if (error) {
    showStatus("JSON error: " + error);
    return;
  }

  elements.downloadPdfBtn.disabled = true;
  elements.downloadPdfBtn.textContent = "Generating...";
  showStatus("");

  try {
    const filename = elements.templateNameInput.value.trim() || "document";
    const res = await apiRequest("/render/pdf", {
      method: "POST",
      body: JSON.stringify({ template: tplValue, data, filename })
    });

    if (!res.ok) {
      const payload = await res.json().catch(() => ({}));
      showStatus(payload.error || "PDF generation failed");
      return;
    }

    // Download the PDF
    const blob = await res.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename + ".pdf";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);

    showSaveStatus("✓ PDF downloaded", 3000);
  } catch (e) {
    showStatus("PDF error: " + e.message);
  } finally {
    elements.downloadPdfBtn.disabled = false;
    elements.downloadPdfBtn.textContent = "Download PDF";
  }
}

// =============================================================================
// Save Template
// =============================================================================

async function saveTemplate() {
  const name = elements.templateNameInput.value.trim();
  if (!name) {
    showStatus("Please enter a template name");
    return;
  }

  const tplValue = tplEditor.state.doc.toString();
  const varsValue = varsEditor.state.doc.toString();

  if (!tplValue) {
    showStatus("Template is empty");
    return;
  }

  elements.saveBtn.disabled = true;
  showSaveStatus("Saving...", 0);
  showStatus("");

  try {
    const res = await apiRequest("/templates/save", {
      method: "POST",
      body: JSON.stringify({
        name: name,
        template: tplValue,
        data: varsValue
      })
    });

    const payload = await res.json();
    if (!res.ok) {
      showStatus(payload.error || "Save failed");
      showSaveStatus("");
    } else {
      showSaveStatus(`✓ Saved as ${payload.filename}`);
    }
  } catch (e) {
    showStatus("Save error: " + e.message);
    showSaveStatus("");
  } finally {
    elements.saveBtn.disabled = false;
  }
}

// =============================================================================
// DMS Upload
// =============================================================================

async function loadTemplateList() {
  try {
    const res = await apiRequest("/templates/list");
    const payload = await res.json();
    if (res.ok && payload.templates) {
      elements.dmsFilenameSelect.innerHTML = payload.templates
        .map(t => `<option value="${t.filename}">${t.filename}</option>`)
        .join('');
    }
  } catch (e) {
    console.error("Failed to load templates:", e);
  }
}

async function openUploadModal() {
  await loadTemplateList();
  elements.dmsRefIdInput.value = crypto.randomUUID();
  elements.uploadModal.classList.add("open");
}

function closeUploadModal() {
  elements.uploadModal.classList.remove("open");
}

async function uploadToDMS() {
  const filename = elements.dmsFilenameSelect.value;
  const refId = elements.dmsRefIdInput.value.trim();

  if (!filename) {
    showStatus("Please select a template file");
    return;
  }
  if (!refId) {
    showStatus("Reference ID is required");
    return;
  }

  elements.uploadToDMSBtn.disabled = true;
  elements.uploadToDMSBtn.textContent = "Uploading...";
  showStatus("");

  try {
    const res = await apiRequest("/templates/upload-dms", {
      method: "POST",
      body: JSON.stringify({
        filename: filename,
        ref_id: refId,
        id_type: elements.dmsIdTypeInput.value.trim(),
        document_type: elements.dmsDocTypeInput.value.trim(),
        source_system: elements.dmsSourceSystemInput.value.trim(),
        document_sequence: elements.dmsDocSequenceInput.value.trim()
      })
    });

    const payload = await res.json();
    if (!res.ok || !payload.success) {
      showStatus(payload.error || "Upload failed");
    } else {
      showSaveStatus(`✓ Uploaded ${filename} to DMS`);
      closeUploadModal();
    }
  } catch (e) {
    showStatus("Upload error: " + e.message);
  } finally {
    elements.uploadToDMSBtn.disabled = false;
    elements.uploadToDMSBtn.textContent = "Upload";
  }
}

// =============================================================================
// Event Listeners
// =============================================================================

// Render
elements.renderBtn.addEventListener("click", doRender);

// PDF Download
elements.downloadPdfBtn.addEventListener("click", downloadPdf);

// Format buttons
elements.formatVarsBtn.addEventListener("click", () => {
  const { data, error } = parseJsonData(varsEditor.state.doc.toString());
  if (error) {
    showStatus("JSON format error: " + error);
  } else {
    setEditorContent(varsEditor, JSON.stringify(data, null, 2));
    showStatus("");
  }
});

elements.formatTplBtn.addEventListener("click", () => {
  try {
    const formatted = formatHtml(tplEditor.state.doc.toString());
    setEditorContent(tplEditor, formatted);
    showStatus("");
  } catch (e) {
    showStatus("HTML format error: " + e.message);
  }
});

// Save
elements.saveBtn.addEventListener("click", saveTemplate);

// DMS Upload Modal
elements.openUploadModalBtn.addEventListener("click", openUploadModal);
elements.closeUploadModalBtn.addEventListener("click", closeUploadModal);
elements.uploadToDMSBtn.addEventListener("click", uploadToDMS);

elements.uploadModal.addEventListener("click", (e) => {
  if (e.target === elements.uploadModal) {
    closeUploadModal();
  }
});

// =============================================================================
// Initialize
// =============================================================================

checkApiHealth();
setInterval(checkApiHealth, API_HEALTH_CHECK_INTERVAL);
doRender();
