package host

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strings"
)

var (
	profileStart = regexp.MustCompile(`(?i)# profile\s+([a-z0-9-_]+)\s*`)
	profileEnd   = regexp.MustCompile(`(?i)# end\s*`)
	disableRe    = regexp.MustCompile(`^#\s*`)
	spaceRemover = regexp.MustCompile(`\s+`)
	tabReplacer  = regexp.MustCompile(`\t+`)
)

type hostFile struct {
	profiles profileMap
}

type profileMap map[string]hostLines

type hostLines []string

// IsHostLine checks if a line is a host line or a comment line.
func IsHostLine(line string) bool {
	p := strings.Split(cleanLine(line), " ")
	i := 0
	if p[0] == "#" && len(p) > 1 {
		i = 1
	}
	ip := net.ParseIP(p[i])

	return ip != nil
}

// Read returns hosts file content grouped by profiles.
// If you pass strict=true it would remove all comments.
func Read(r io.Reader, strict bool) (*hostFile, error) {
	h := &hostFile{
		profiles: profileMap{},
	}

	ln := 0
	s := bufio.NewScanner(r)
	open := ""
	for s.Scan() {
		ln++
		b := s.Bytes()

		switch {

		case profileStart.Match(b):
			open = strings.TrimSpace(strings.Split(string(b), "# profile")[1])
			h.profiles[open] = []string{}

		case profileEnd.Match(b):
			open = ""

		case open != "":
			if strict && !IsHostLine(string(b)) {
				// skip
			} else {
				h.profiles[open] = append(h.profiles[open], string(b))
			}

		default:
			h.profiles["default"] = append(h.profiles["default"], string(b))
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func cleanLine(line string) string {
	return tabReplacer.ReplaceAllString(spaceRemover.ReplaceAllString(line, " "), " ")
}

// IsDisabled check if a line starts with a # comment marker.
func IsDisabled(line string) bool {
	return disableRe.MatchString(line)
}
