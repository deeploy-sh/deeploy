package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/deeploy-sh/deeploy/internal/docs/ctxkeys"
	"github.com/deeploy-sh/deeploy/internal/docs/service"
	"github.com/deeploy-sh/deeploy/internal/docs/ui/pages"
)

// DocsHandler handles documentation requests
type DocsHandler struct {
	docs *service.DocsService
}

// NewDocsHandler creates a new docs handler
func NewDocsHandler(docs *service.DocsService) *DocsHandler {
	return &DocsHandler{docs: docs}
}

// DocPage renders a documentation page
func (h *DocsHandler) DocPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	slug = strings.TrimPrefix(slug, "/")

	doc, err := h.docs.GetDoc(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tree := h.docs.GetTree()
	prev, next := h.docs.GetPrevNext(doc.Slug)

	// Add URL path to context
	ctx := context.WithValue(r.Context(), ctxkeys.URLPathKey, doc.Path)

	// Check if this is an HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		// Only render the content fragment
		pages.DocContent(doc, prev, next).Render(ctx, w)
		return
	}

	// Full page render
	pages.Doc(doc, tree, prev, next).Render(ctx, w)
}

// Sitemap returns the sitemap.xml
func (h *DocsHandler) Sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write(h.docs.GenerateSitemap())
}

// Robots returns the robots.txt
func (h *DocsHandler) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write(h.docs.GenerateRobots())
}
