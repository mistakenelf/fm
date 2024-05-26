package utilities

// Contains checks if a string or int is in a slice.
func Contains[T string | int](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
