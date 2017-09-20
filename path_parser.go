package main

func parsePath(path string) []string {
	return parsePathAccum([]string{}, path)
}

func parsePathAccum(paths []string, remaining string) []string {
	head, tail := nextYamlPath(remaining)
	if tail == "" {
		return append(paths, head)
	}
	return parsePathAccum(append(paths, head), tail)
}

func nextYamlPath(path string) (pathElement string, remaining string) {
	switch path[0] {
	case '[':
		// e.g [0].blah.cat -> we need to return "0" and "blah.cat"
		return search(path[1:], []uint8{']'}, true)
	case '"':
		// e.g "a.b".blah.cat -> we need to return "a.b" and "blah.cat"
		return search(path[1:], []uint8{'"'}, true)
	default:
		// e.g "a.blah.cat" -> return "a" and "blah.cat"
		return search(path[0:], []uint8{'.', '['}, false)
	}
}

func search(path string, matchingChars []uint8, skipNext bool) (pathElement string, remaining string) {
	for i := 0; i < len(path); i++ {
		var char = path[i]
		if contains(matchingChars, char) {
			var remainingStart = i + 1
			if skipNext {
				remainingStart = remainingStart + 1
			} else if !skipNext && char != '.' {
				remainingStart = i
			}
			if remainingStart > len(path) {
				remainingStart = len(path)
			}
			return path[0:i], path[remainingStart:]
		}
	}
	return path, ""
}

func contains(matchingChars []uint8, candidate uint8) bool {
	for _, a := range matchingChars {
		if a == candidate {
			return true
		}
	}
	return false
}
