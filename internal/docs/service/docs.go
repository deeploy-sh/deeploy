package service

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Doc represents a documentation page
type Doc struct {
	Title       string
	Description string
	Slug        string        // URL path segment (e.g., "guides/installation")
	Path        string        // Full URL path (e.g., "/docs/guides/installation")
	Order       int           // From frontmatter
	Content     template.HTML // Rendered HTML
	Children    []*Doc        // Child pages (for categories)
	Parent      *Doc          // Parent reference
}

// DocsService manages documentation content
type DocsService struct {
	tree    *Doc             // Hierarchical tree
	bySlug  map[string]*Doc  // Lookup by slug
	flat    []*Doc           // Flat list for navigation
	baseURL string
	md      goldmark.Markdown
}

// NewDocsService creates a new documentation service
func NewDocsService(content embed.FS, baseURL string) (*DocsService, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
		),
	)

	s := &DocsService{
		baseURL: baseURL,
		md:      md,
		bySlug:  make(map[string]*Doc),
	}

	if err := s.buildTree(content); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *DocsService) buildTree(content embed.FS) error {
	s.tree = &Doc{
		Title:    "Documentation",
		Slug:     "",
		Path:     "/docs",
		Children: []*Doc{},
	}

	// Track directory metadata from _index.md files
	dirMetadata := make(map[string]*Doc)

	err := fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}

		// Skip embed.go
		if strings.HasSuffix(path, ".go") {
			return nil
		}

		data, err := content.ReadFile(path)
		if err != nil {
			return err
		}

		// Normalize path
		relPath := filepath.ToSlash(path)

		page, err := s.loadDocPage(data, relPath)
		if err != nil {
			return err
		}

		// Handle _index.md specially
		if strings.HasSuffix(relPath, "_index.md") {
			dir := strings.TrimSuffix(relPath, "/_index.md")
			if dir == "_index.md" {
				dir = ""
			}
			dirMetadata[dir] = page
			return nil
		}

		s.insertPage(page, relPath, dirMetadata)
		return nil
	})

	if err != nil {
		return err
	}

	s.sortTree(s.tree)
	s.buildFlatList()
	s.buildSlugIndex()

	return nil
}

func (s *DocsService) loadDocPage(data []byte, relPath string) (*Doc, error) {
	context := parser.NewContext()
	var buf bytes.Buffer

	if err := s.md.Convert(data, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	// Extract frontmatter
	meta := make(map[string]any)
	fm := frontmatter.Get(context)
	if fm != nil {
		_ = fm.Decode(&meta)
	}

	slug := strings.TrimSuffix(relPath, ".md")
	page := &Doc{
		Slug:     slug,
		Path:     "/docs/" + slug,
		Content:  template.HTML(buf.String()),
		Children: []*Doc{},
	}

	if title, ok := meta["title"].(string); ok {
		page.Title = title
	} else {
		page.Title = s.titleFromSlug(slug)
	}

	if desc, ok := meta["description"].(string); ok {
		page.Description = desc
	}

	if order, ok := meta["order"].(int); ok {
		page.Order = order
	} else if orderFloat, ok := meta["order"].(float64); ok {
		page.Order = int(orderFloat)
	}

	return page, nil
}

func (s *DocsService) insertPage(page *Doc, relPath string, dirMetadata map[string]*Doc) {
	parts := strings.Split(relPath, "/")
	current := s.tree

	// Build/traverse directory structure
	for i := 0; i < len(parts)-1; i++ {
		dirSlug := strings.Join(parts[:i+1], "/")

		// Check if directory exists
		var found *Doc
		for _, child := range current.Children {
			if child.Slug == dirSlug {
				found = child
				break
			}
		}

		if found == nil {
			// Create new directory node
			dirPage := &Doc{
				Slug:     dirSlug,
				Path:     "/docs/" + dirSlug,
				Children: []*Doc{},
				Parent:   current,
			}

			// Apply metadata from _index.md
			if meta, ok := dirMetadata[dirSlug]; ok {
				dirPage.Title = meta.Title
				dirPage.Description = meta.Description
				dirPage.Order = meta.Order
				dirPage.Content = meta.Content
			} else {
				dirPage.Title = s.titleFromSlug(parts[i])
			}

			current.Children = append(current.Children, dirPage)
			current = dirPage
		} else {
			current = found
		}
	}

	page.Parent = current
	current.Children = append(current.Children, page)
}

func (s *DocsService) sortTree(node *Doc) {
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].Order != node.Children[j].Order {
			return node.Children[i].Order < node.Children[j].Order
		}
		return node.Children[i].Title < node.Children[j].Title
	})

	for _, child := range node.Children {
		s.sortTree(child)
	}
}

func (s *DocsService) buildFlatList() {
	s.flat = []*Doc{}
	s.collectPagesInOrder(s.tree, &s.flat)
}

func (s *DocsService) collectPagesInOrder(node *Doc, pages *[]*Doc) {
	// Only add content pages (no children = not a category)
	if node.Slug != "" && len(node.Children) == 0 {
		*pages = append(*pages, node)
	}

	for _, child := range node.Children {
		s.collectPagesInOrder(child, pages)
	}
}

func (s *DocsService) buildSlugIndex() {
	s.indexPage(s.tree)
}

func (s *DocsService) indexPage(node *Doc) {
	if node.Slug != "" {
		s.bySlug[node.Slug] = node
	}
	for _, child := range node.Children {
		s.indexPage(child)
	}
}

func (s *DocsService) titleFromSlug(slug string) string {
	parts := strings.Split(slug, "/")
	lastPart := parts[len(parts)-1]

	lastPart = strings.ReplaceAll(lastPart, "-", " ")
	lastPart = strings.ReplaceAll(lastPart, "_", " ")

	words := strings.Fields(lastPart)
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}

	return strings.Join(words, " ")
}

// GetDoc returns a documentation page by slug
func (s *DocsService) GetDoc(slug string) (*Doc, error) {
	if slug == "" {
		// Return first content page
		if len(s.flat) > 0 {
			return s.flat[0], nil
		}
		return nil, fmt.Errorf("no documentation pages found")
	}

	page, ok := s.bySlug[slug]
	if !ok {
		return nil, fmt.Errorf("documentation page not found: %s", slug)
	}

	return page, nil
}

// GetTree returns the documentation tree
func (s *DocsService) GetTree() *Doc {
	return s.tree
}

// GetFlatList returns all documentation pages in navigation order
func (s *DocsService) GetFlatList() []*Doc {
	return s.flat
}

// GetPrevNext returns the previous and next pages
func (s *DocsService) GetPrevNext(slug string) (prev, next *Doc) {
	for i, page := range s.flat {
		if page.Slug == slug {
			if i > 0 {
				prev = s.flat[i-1]
			}
			if i < len(s.flat)-1 {
				next = s.flat[i+1]
			}
			break
		}
	}
	return
}

// GenerateSitemap returns the sitemap XML
func (s *DocsService) GenerateSitemap() []byte {
	var buf bytes.Buffer
	today := time.Now().Format("2006-01-02")

	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	// Homepage
	buf.WriteString(fmt.Sprintf(`<url><loc>%s/</loc><lastmod>%s</lastmod><priority>1.0</priority></url>`, s.baseURL, today))

	// Docs index
	buf.WriteString(fmt.Sprintf(`<url><loc>%s/docs</loc><lastmod>%s</lastmod><priority>0.9</priority></url>`, s.baseURL, today))

	// All doc pages
	for _, doc := range s.flat {
		priority := "0.8"
		if strings.Contains(doc.Slug, "/") {
			priority = "0.7"
		}
		buf.WriteString(fmt.Sprintf(`<url><loc>%s%s</loc><lastmod>%s</lastmod><priority>%s</priority></url>`,
			s.baseURL, doc.Path, today, priority))
	}

	buf.WriteString(`</urlset>`)
	return buf.Bytes()
}

// GenerateRobots returns the robots.txt content
func (s *DocsService) GenerateRobots() []byte {
	return []byte(fmt.Sprintf(`User-agent: *
Allow: /

Sitemap: %s/sitemap.xml
`, s.baseURL))
}
