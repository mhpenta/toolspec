package toolspec

type MediaType string

const (
	MediaTypeJPEG MediaType = "image/jpeg"
	MediaTypePNG  MediaType = "image/png"
	MediaTypeGIF  MediaType = "image/gif"
	MediaTypeWebP MediaType = "image/webp"
	Unknown       MediaType = "unknown"
)
