package zengarden

// Context is a simple dictionary for use in text/template files.
type Context map[string]interface{}

func str(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}

	return ""
}
