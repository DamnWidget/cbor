package test

const (
	Succeed string = "\x1b[32m\u2713\x1b[0m"
	Failed  string = "\x1b[31m\u2717\x1b[0m"
)

// BytesEqual is used to compare equality between two bytes sequences
func BytesEqual(a, b []byte) bool {

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
