package hosts

import (
	"bufio"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	profileStart = regexp.MustCompile(`(?i)# profile\s+([a-z0-9-_]+)\s*`)
	profileEnd   = regexp.MustCompile(`(?i)# end\s*`)
	disableRe    = regexp.MustCompile(`^#\s*`)
	spaceRemover = regexp.MustCompile(`\s+`)
	tabReplacer  = regexp.MustCompile(`\t+`)
)

type ListOptions struct {
	Profile string
}

type HostEntry []string

type HostFile struct {
	Profiles map[string]HostEntry
}

// ReadFromFile
// If you pass strict=true it would remove all comments.
func ParseFile(filename string, strict bool) (*HostFile, error) {
	hostFile := &HostFile{
		Profiles: make(map[string]HostEntry),
	}
	fileHandler, err := os.Open(filename)
	if err != nil {
		return hostFile, err
	}
	line := 0
	scanLine := bufio.NewScanner(fileHandler)
	open := ""
	for scanLine.Scan() {
		line++
		b := scanLine.Bytes()

		switch {
		case profileStart.Match(b):
			open = strings.TrimSpace(strings.Split(string(b), "# profile")[1])
			hostFile.Profiles[open] = []string{}

		case profileEnd.Match(b):
			open = ""

		case open != "":
			if strict && !IsHostLine(string(b)) {
				// skip
			} else {
				hostFile.Profiles[open] = append(hostFile.Profiles[open], string(b))
			}

		default:
			hostFile.Profiles["default"] = append(hostFile.Profiles["default"], string(b))
		}

		if err := scanLine.Err(); err != nil {
			return nil, err
		}
	}
	return hostFile, nil
}

// ListProfiles shows a table with profile names status and routing information
func (h *HostFile) ListProfiles(opts *ListOptions) error {
	var profile = ""
	if opts.Profile != "" {
		profile = opts.Profile
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Profile", "Status", "IP", "Domain"})

	if profile == "default" || profile == "" {
		appendProfile("default", table, h.Profiles["default"])

		// if len(h.profiles) > 1 {
		// 	table.AddSeparator()
		// }
	}

	i := 0
	for p, data := range h.Profiles {
		i++
		if profile != "" && p != profile {
			continue
		}
		if p == "default" {
			continue
		}

		appendProfile(p, table, data)

		// if i < len(h.profiles) {
		// 	table.AddSeparator()
		// }
	}
	table.Render()
	return nil
}

func appendProfile(profile string, table *tablewriter.Table, data HostEntry) {
	for _, r := range data {
		if r == "" {
			continue
		}
		if !IsHostLine(r) {
			continue
		}
		rs := strings.Split(cleanLine(r), " ")

		status := "on"
		ip, domain := rs[0], rs[1]
		if IsDisabled(r) {
			// skip empty comments lines
			if rs[1] == "" {
				continue
			}
			status = "off"
			ip, domain = rs[1], rs[2]
		}
		table.Append([]string{
			profile,
			status,
			ip,
			domain,
		})
	}
}

// IsDisabled check if a line starts with a # comment marker.
func IsDisabled(line string) bool {
	return disableRe.MatchString(line)
}

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

func cleanLine(line string) string {
	return tabReplacer.ReplaceAllString(spaceRemover.ReplaceAllString(line, " "), " ")
}
