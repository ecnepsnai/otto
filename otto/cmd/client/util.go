package main

func addressFromSocketString(s string) string {
	// Remove the port first
	portIdx := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			portIdx = i
			break
		}
	}

	s = s[:portIdx]

	if s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}

	return s
}
