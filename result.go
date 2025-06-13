package toolspec

// ToolResult represents the outcome of a tool execution, providing structured output
// for different types of tool responses.
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
}
