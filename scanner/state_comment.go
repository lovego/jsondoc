package scanner

// stateInCommentWhenBeginString is the state after reading `#`, when in begin string state.
func stateInCommentWhenBeginString(s *scanner, c byte) int {
	if c == '\n' {
		s.step = stateBeginString
		return scanEndComment
	}
	return scanContinue
}

// stateInCommentWhenBeginValue is the state after reading `#`, when in begin value state.
func stateInCommentWhenBeginValue(s *scanner, c byte) int {
	if c == '\n' {
		s.step = stateBeginValue
		return scanEndComment
	}
	return scanContinue
}

// stateInCommentWhenEndValue is the state after reading `#`, when in end value state.
func stateInCommentWhenEndValue(s *scanner, c byte) int {
	if c == '\n' {
		s.step = stateEndValue
		return scanSkipSpace
	}
	return scanContinue
}
