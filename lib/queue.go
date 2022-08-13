package lib

type Queue []string

// IsEmpty: check if Queue is empty
func (s *Queue) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the Queue
func (s *Queue) Push(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the Queue
}

// Return the top element of the Queue, but don't modify the Queue
func (s *Queue) Peek() (string, bool) {
	if s.IsEmpty() {
		return "", true
	} else {
		element := (*s)[0]
		return element, false
	}
}

// Remove and return top element of Queue. Return true if Queue is empty.
func (s *Queue) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", true
	} else {
		element := (*s)[0] // Index into the slice and obtain the element.
		*s = (*s)[1:]      // Remove it from the Queue by slicing it off.
		return element, false
	}
}
