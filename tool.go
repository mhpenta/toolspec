package toolspec

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tool defines the interface that all tools must implement
type Tool interface {
	// Spec returns the tool's specification, including name, description, parameters, and UI hints.
	Spec() *ToolSpec
	// Execute runs the tool with given parameters
	Execute(ctx context.Context, params json.RawMessage) (*ToolResult, error)
}

type ToolSpec struct {
	// Name returns the tool's identifier
	Name string `json:"name,omitempty"`
	// Type returns the tool's type, which is used for categorization
	Type string `json:"type,omitempty"`
	// Description returns the tool's description
	Description string `json:"description,omitempty"`
	// Parameters returns the tool's parameter schema
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	// UI provides additional UI hints for the tool
	UI UI `json:"ui,omitempty"`
}

type UI struct {
	// Verb is a present progressive verb phrase for UI display (e.g., "Searching for companies")
	Verb string `json:"verb,omitempty"`
	// LongRunning indicates if the tool is expected to run for a long time, resulting in different handling in the UI
	LongRunning bool `json:"long_running,omitempty"`
}

func Validate(t Tool) error {
	return validateSpec(t.Spec())
}

const (
	maxToolNameLength = 64
)

func validateSpec(m *ToolSpec) error {
	if m.Name == "" {
		return fmt.Errorf("tool spec must include a non-empty name")
	}

	if len(m.Name) > maxToolNameLength {
		return fmt.Errorf("tool name must not exceed 64 characters")
	}

	for _, char := range m.Name {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-' {
			continue
		}
		return fmt.Errorf("tool name must contain only alphanumeric characters, underscores, or hyphens")
	}

	if m.Description == "" {
		return fmt.Errorf("tool spec description cannot be empty")
	}

	if m.Parameters == nil {
		return fmt.Errorf("tool spec parameters cannot be nil")
	}

	return nil
}
