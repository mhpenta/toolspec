package httprequest

import (
	"bytes"
	"context"
	"fmt"
	tools "github.com/mhpenta/toolspec"
	"golang.org/x/net/html"
	"net/http"

	"io"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// NewHTTPTool creates a new HTTP tool with a timeout using the NewTool function
func NewHTTPTool(timeoutSeconds int) tools.Tool {
	httpTool := &HTTPTool{
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}

	return tools.NewTool(
		"http_request",
		"Make HTTP requests to external services and APIs, with optional plain text conversion",
		httpTool.executeRequest,
		tools.WithType("function"),
		tools.WithVerb("Making HTTP request"),
	)
}

// HTTPParams defines the parameters for the HTTP tool
type HTTPParams struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	RequestBody string            `json:"request_body,omitempty"`
	PlainText   bool              `json:"plain_text,omitempty"`
}

// HTTPTool is a tool for making HTTP requests
type HTTPTool struct {
	client *http.Client
}

// executeRequest is the handler function that processes HTTP requests
func (t *HTTPTool) executeRequest(ctx context.Context, params HTTPParams) (string, error) {
	resp, err := t.ExecuteRequest(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if params.PlainText {
		plaintext := htmlToPlaintext(body)
		return string(plaintext), nil
	}

	return string(body), nil
}

// ExecuteRequest performs the actual HTTP request (kept separate for potential reuse)
func (t *HTTPTool) ExecuteRequest(ctx context.Context, params HTTPParams) (*http.Response, error) {
	_, err := url.Parse(params.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, params.Method, params.URL, strings.NewReader(params.RequestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

const (
	whiteSpaceRegex  = `(\n\s*){3,}`
	nonBreakingSpace = `\u00a0`
	zeroWidthSpace   = `\u200b`
)

// htmlToPlaintext removes all HTML tags from an array of bytes, adding a \n when a <div> tag is encountered
// with a maximum of three consecutive newlines
func htmlToPlaintext(htmlBytes []byte) []byte {
	plainText := cleanHTMLWhitespace(stripHtmlAsBytesDivsToWhiteSpace(htmlBytes, '\n'))
	re := regexp.MustCompile(whiteSpaceRegex)
	return re.ReplaceAll(plainText, []byte(`\n`))
}

// stripHtmlAsBytesDivsToWhiteSpace processes HTML bytes, replacing <div> and <tr> tags with a specified whitespace character
func stripHtmlAsBytesDivsToWhiteSpace(htmlBytes []byte, whitespaceCharForDivBreaks rune) []byte {
	var b bytes.Buffer
	r := bytes.NewReader(htmlBytes)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return b.Bytes()
		case html.TextToken:
			b.Write(z.Text())
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if token.Data == "div" {
				b.WriteRune(whitespaceCharForDivBreaks)
			} else if token.Data == "tr" {
				b.WriteRune(whitespaceCharForDivBreaks)
			}
		default:
			// Ignore other token types
		}
	}
}

// cleanHTMLWhitespace removes all HTML whitespace characters from the input string.
func cleanHTMLWhitespace(input []byte) []byte {
	input = bytes.ReplaceAll(input, []byte("“"), []byte("\""))
	input = bytes.ReplaceAll(input, []byte("”"), []byte("\""))
	return bytes.ReplaceAll(bytes.ReplaceAll(input, []byte(zeroWidthSpace), []byte("")), []byte(nonBreakingSpace), []byte(" "))
}
