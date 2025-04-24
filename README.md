# ToolSpec

A standardized Go interface for building interchangeable tools.

## The Tool Interface

```go
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
```

## Standardized Results

Tools return a consistent `ToolResult` structure that separates concerns:

```go
type ToolResult struct {
    Name             string             `json:"name,omitempty"`
    Output           *string            `json:"output,omitempty"`
    Error            *string            `json:"error,omitempty"`
    System           *string            `json:"system,omitempty"`
    Image            *ToolImage         `json:"image,omitempty"`
    CitableDocuments []*CitableDocument `json:"citable_document,omitempty"`
}
```

## Why Use This Pattern?

Interface-based tools provide several advantages over function registration patterns:

1. **Composition** - Objects implementing the Tool interface can contain state and configuration
2. **Testability** - Each tool can be isolated and unit tested effectively
3. **Polymorphism** - Tools can be used interchangeably regardless of implementation
4. **Discoverability** - Tools expose their parameters and description as part of the interface
5. **Result consistency** - The standard ToolResult structure ensures uniform error handling and output formats
6. **Type safety** - Interface constraints enforce contracts at compile time
7. **Extensibility** - Tools can be extended via composition or embedding

This approach aligns with Go's interface philosophy of "design by contract" while enabling a flexible plugin system without sacrificing type safety.