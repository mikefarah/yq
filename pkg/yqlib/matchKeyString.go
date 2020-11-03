package yqlib

func Match(name string, pattern string) (matched bool) {
	if pattern == "" {
		return name == pattern
	}
	log.Debug("pattern: %v", pattern)
	if pattern == "*" {
		log.Debug("wild!")
		return true
	}
	return deepMatch([]rune(name), []rune(pattern))
}

func deepMatch(str, pattern []rune) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return deepMatch(str, pattern[1:]) ||
				(len(str) > 0 && deepMatch(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}
