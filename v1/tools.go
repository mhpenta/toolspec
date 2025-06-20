package v1

import (
	"context"
	"encoding/json"
)

type ToolImage struct {
	Base64Image string `json:"base64_image"`
	ContentType string `json:"content_type"`
}

// ToolResult represents the outcome of a tool execution, providing structured output
// for different types of tool responses. It separates concerns between normal output,
// errors, image data, and system-level information to facilitate proper handling.
type ToolResult struct {
	// Name contains the name of the tool
	Name string `json:"name,omitempty"`

	// Output contains the standard/success result of a tool execution.
	// This field should be used for the primary, expected output of the tool.
	Output any `json:"output,omitempty"`

	// Error contains any error messages when tool execution fails.
	// This is kept separate from Output to clearly distinguish between
	// success and failure cases.
	Error *string `json:"error,omitempty"`

	// System contains metadata or system-level messages about the tool's execution.
	// This field is used for information about the execution environment,
	// tool initialization messages, or other system-level status updates
	// that are separate from the tool's primary output or errors.
	System *string `json:"system,omitempty"`

	// Image contains any image data generated by the tool execution. E.g., for tools
	// that perform screen captures or generate visual output. The ToolImage type encapsulates the
	// image data and related metadata.
	Image *ToolImage `json:"image,omitempty"`

	// CitableDocuments contains a list of citable documents related to the tool's output. This is critical
	CitableDocuments []*CitableDocument `json:"citable_document,omitempty"`
}

// Tool defines the interface that all tools must implement
type Tool interface {
	// Name returns the tool's identifier
	Name() string
	// Type returns the tool's API type
	Type() string
	// Execute runs the tool with given parameters
	Execute(ctx context.Context, params json.RawMessage) (*ToolResult, error)
	// Parameters returns the tool's parameter schema
	Parameters() map[string]interface{}
	// Description returns the tool's description
	Description() string
}
