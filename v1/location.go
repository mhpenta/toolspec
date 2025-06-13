package v1

import "errors"

type CitationLocationType string

const (
	CitationLocationTypeChar  CitationLocationType = "char"
	CitationLocationTypePage  CitationLocationType = "page"
	CitationLocationTypeBlock CitationLocationType = "block"
)

type Location struct {
	Type  string `json:"type"` // "char", "page", "block"
	Start int    `json:"start"`
	End   int    `json:"end"`
}

func (l *Location) Validate() error {
	if l.Type != string(CitationLocationTypeChar) && l.Type != string(CitationLocationTypePage) && l.Type != string(CitationLocationTypeBlock) {
		return errors.New("invalid location type")
	}
	return nil
}
