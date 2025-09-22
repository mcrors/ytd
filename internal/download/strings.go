package download

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

func normalizeTwoLevel(rel string) (string, error) {
	if rel == "" {
		return "", errors.New("path required: use genre/channel")
	}

	rel = filepath.ToSlash(strings.TrimSpace(rel))

	if strings.HasPrefix(rel, "/") || strings.Contains(rel, `..`) {
		return "", errors.New("invalid path")
	}

	parts := strings.Split(rel, "/")
	if len(parts) != 2 {
		return "", errors.New("path must be exactly two levels: genre/channel")
	}

	g := slugify(parts[0])
	c := slugify(parts[1])
	if g == "" || c == "" {
		return "", errors.New("invalid path after slugify")
	}

	clean := g + "/" + c

	if !twoSegRe.MatchString(clean) {
		return "", errors.New("invalid characters in path")
	}
	return clean, nil
}

var twoSegRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*/[a-z0-9]+(?:-[a-z0-9]+)*$`)

// slug: lowercase, spaces/_ -> "-", keep a-z0-9-, collapse dashes
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")

	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case r == '-':
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		default:
			// drop everything else
		}
	}
	out := strings.Trim(b.String(), "-")
	return out
}
