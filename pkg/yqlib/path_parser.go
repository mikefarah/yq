package yqlib

type PathParser interface {
	ParsePath(path string) []string
}

type pathParser struct{}

func NewPathParser() PathParser {
	return &pathParser{}
}

func (p *pathParser) ParsePath(path string) []string {
	return p.parsePathAccum([]string{}, path)
}

func (p *pathParser) parsePathAccum(paths []string, remaining string) []string {
	head, tail := p.nextYamlPath(remaining)
	if tail == "" {
		return append(paths, head)
	}
	return p.parsePathAccum(append(paths, head), tail)
}

func (p *pathParser) nextYamlPath(path string) (pathElement string, remaining string) {
	switch path[0] {
	case '[':
		// e.g [0].blah.cat -> we need to return "0" and "blah.cat"
		return p.search(path[1:], []uint8{']'}, true)
	case '"':
		// e.g "a.b".blah.cat -> we need to return "a.b" and "blah.cat"
		return p.search(path[1:], []uint8{'"'}, true)
	default:
		// e.g "a.blah.cat" -> return "a" and "blah.cat"
		return p.search(path[0:], []uint8{'.', '['}, false)
	}
}

func (p *pathParser) search(path string, matchingChars []uint8, skipNext bool) (pathElement string, remaining string) {
	for i := 0; i < len(path); i++ {
		var char = path[i]
		if p.contains(matchingChars, char) {
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

func (p *pathParser) contains(matchingChars []uint8, candidate uint8) bool {
	for _, a := range matchingChars {
		if a == candidate {
			return true
		}
	}
	return false
}
