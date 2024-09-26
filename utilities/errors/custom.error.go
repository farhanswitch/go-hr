package errorutility

type HashidError struct{}

func (hashErr *HashidError) Error() string {
	return "[ERR109] Invalid value."
}
