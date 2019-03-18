package zengarden

type context map[string]interface{}

func str(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}

	return ""
}
