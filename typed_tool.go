package toolspec

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhpenta/jobj/funcschema"
	"github.com/mhpenta/jobj/safeunmarshal"
)

type TypedTool[In, Out any] struct {
	spec    *ToolSpec
	handler func(context.Context, In) (Out, error)
}

func (t *TypedTool[In, Out]) Spec() *ToolSpec {
	return t.spec
}

func (t *TypedTool[In, Out]) Execute(ctx context.Context, params json.RawMessage) (*ToolResult, error) {
	var input In
	if len(params) > 0 {
		parsedInput, err := safeunmarshal.To[In](params)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		input = parsedInput
	}
	result, err := t.handler(ctx, input)
	if err != nil {
		return nil, err
	}
	return &ToolResult{
		Output: result,
		Error:  nil,
	}, nil
}

// ToolOption for functional configuration
type ToolOption func(*ToolSpec)

func WithType(toolType string) ToolOption {
	return func(spec *ToolSpec) {
		spec.Type = toolType
	}
}

func WithVerb(verb string) ToolOption {
	return func(spec *ToolSpec) {
		spec.UI.Verb = verb
	}
}

func WithLongRunning(longRunning bool) ToolOption {
	return func(spec *ToolSpec) {
		spec.UI.LongRunning = longRunning
	}
}

func WithCustomSchema(schema map[string]interface{}) ToolOption {
	return func(spec *ToolSpec) {
		spec.Parameters = schema
	}
}

func NewTool[In, Out any](
	name,
	description string,
	handler func(context.Context, In) (Out, error),
	opts ...ToolOption,
) Tool {

	schema, err := funcschema.SafeSchemaFromFunc(handler)
	if err != nil {
		return nil
	}

	spec := &ToolSpec{
		Name:        name,
		Type:        fmt.Sprintf("%s_v1", name),
		Description: description,
		Parameters:  schema,
		Sequential:  false,
		UI:          UI{},
	}

	for _, opt := range opts {
		opt(spec)
	}

	return &TypedTool[In, Out]{
		spec:    spec,
		handler: handler,
	}
}
