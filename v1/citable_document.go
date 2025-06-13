package v1

import (
	"fmt"
)

type DocumentSourceType string
type CitationType string

const (
	CitationTypeCharLocation CitationType = "char_location"
	CitationTypePageNumber   CitationType = "page_number"
	CitationTypeBlockIndex   CitationType = "block_index"
)

const (
	SourceTypeText   DocumentSourceType = "text"
	SourceTypeBase64 DocumentSourceType = "base64"
	SourceTypeCustom DocumentSourceType = "content"
)

// ChunkedDocument is an interface for documents that can be chunked for custom document citation
type ChunkedDocument interface {
	GetChunks() []string
	GetFormattedCitationAtLocation(location Location) string
	GetFormattedCitationFromText(text string) string
}

type CitableDocument struct {
	// UniqueTitle title, used to identify the document among documents
	UniqueTitle string `json:"unique_title,omitempty"`

	// Title, official title
	Title string `json:"title,omitempty"`

	// Document description or context
	Description string `json:"context,omitempty"`

	// Whether citations are enabled for document
	CitationsEnabled bool `json:"citations,omitempty"`

	// Source of the document
	Source DocumentContent `json:"source"`
}

func (c *CitableDocument) GetFormattedCitationAtLocation(location Location) string {
	if c.Source.Content != nil {
		return c.Source.Content.GetFormattedCitationAtLocation(location)
	}
	return ""
}

type DocumentContent struct {
	// Type of the document, either text (for plaintext), base64 (for pdf) or content, for custom chunking
	Type DocumentSourceType `json:"type"`

	// Media type of the document, one of text/plain, application/pdf (for custom chunking)
	// See anthropic api docs https://docs.anthropic.com/en/docs/build-with-claude/citations#example-plain-text-citation
	MediaType MediaType `json:"media_type,omitempty"`

	// Base64 encoded data for pdfs or plain text for plain text documents, empty if using custom chunking
	Data string `json:"data,omitempty"`

	// Content of the document, if using custom chunking - nil if using data or pdfs
	Content ChunkedDocument `json:"content"`
}

type ContentBlockDocumentChunk struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func createContentBlockDocumentChunkSlice(chunks []string) []ContentBlockDocumentChunk {
	var contentBlockChunks []ContentBlockDocumentChunk
	for _, chunk := range chunks {
		contentBlockChunks = append(contentBlockChunks, ContentBlockDocumentChunk{
			Type: "text",
			Text: chunk,
		})
	}
	return contentBlockChunks
}

func (d *DocumentContent) Validate() error {
	switch d.Type {
	case SourceTypeText, SourceTypeBase64, SourceTypeCustom:
		// valid
	default:
		return fmt.Errorf("invalid source type: %s", d.Type)
	}

	switch d.Type {
	case SourceTypeCustom:
		if d.Content == nil || len(d.Content.GetChunks()) == 0 {
			return fmt.Errorf("content required for content type")
		}
		if d.Data != "" {
			return fmt.Errorf("data not allowed with content type")
		}
	default:
		if d.Content != nil {
			return fmt.Errorf("content not allowed with %s type", d.Type)
		}
		if d.Data == "" {
			return fmt.Errorf("data required for %s type", d.Type)
		}
	}

	return nil
}

func (c *CitableDocument) Validate() error {
	return c.Source.Validate()
}
