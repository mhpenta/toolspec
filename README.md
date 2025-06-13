# ToolSpec

A standardized Go interface for building interchangeable tools.

## The Tool Interface

Every tool implements this simple interface:

```go
type Tool interface {
    // Spec returns the tool's specification including metadata and parameters
    Spec() ToolSpec
    
    // Execute runs the tool with the given parameters
    Execute(ctx context.Context, params map[string]interface{}) ToolResult
}
```

## Tool Specification

Tool metadata is encapsulated in the `ToolSpec` struct:

```go
type ToolSpec struct {
    Name        string                 `json:"name,omitempty"`
    Type        string                 `json:"type,omitempty"`
    Description string                 `json:"description,omitempty"`
    Parameters  map[string]interface{} `json:"parameters,omitempty"`
    UI          UI                     `json:"ui,omitempty"`
}
```

### UI Hints

The `UI` field provides hints for tool presentation:

```go
type UI struct {
    // Verb is a present progressive verb phrase describing what the tool is doing
    // Example: "Searching for companies", "Analyzing data"
    Verb string `json:"verb,omitempty"`
    
    // LongRunning indicates if this tool typically takes a long time to execute
    LongRunning bool `json:"longRunning,omitempty"`
}
```

## Standardized Results

All tools return a `ToolResult` with consistent fields:

```go
type ToolResult struct {
    Name   string      `json:"name"`
    Output interface{} `json:"output,omitempty"`
    Error  string      `json:"error,omitempty"`
}
```

## Why Use This Pattern?

Interface-based tools provide several advantages over function registration patterns:

1. **Composition** - Objects implementing the Tool interface can contain state and configuration
2. **Testability** - Each tool can be isolated and unit tested effectively
3. **Polymorphism** - Tools can be used interchangeably regardless of implementation
4. **Discoverability** - Tools expose their parameters and description as part of the interface
5. **Result consistency** - The standard ToolResult structure ensures uniform error handling and flexible output 
6. **Extensibility** - Tools can be extended via composition or embedding
7. **Better encapsulation** - All metadata consolidated in a single `ToolSpec` struct
8. **UI/UX support** - Built-in hints for better tool presentation in user interfaces