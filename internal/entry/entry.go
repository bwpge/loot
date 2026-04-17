package entry

import (
	"regexp"
	"strings"
)

type Entry struct {
	Value   string   `json:"value"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
	Hosts   []string `json:"hosts"`
}

type Filter struct {
	ID    []string
	Tags  []string
	Hosts []string
}

type Flag struct {
	Type  string `json:"type"`
	Owner string `json:"owner"`
	Host  string `json:"host"`
}

func DetectValues(e Entry) ([]Entry, string) {
	var result []Entry

	// user@domain format
	user, domain, found := strings.Cut(e.Value, "@")
	if found {
		// this is by no means a good domain regex, but is good enough for guessing
		if match, _ := regexp.MatchString(`^[a-zA-Z0-9._-]{2,70}$`, domain); match {
			result = append(result,
				Entry{
					Value:   user,
					Comment: e.Comment,
					Tags:    append(e.Tags, "username"),
					Hosts:   e.Hosts,
				},
				Entry{
					Value:   domain,
					Comment: e.Comment,
					Tags:    append(e.Tags, "domain"),
					Hosts:   e.Hosts,
				},
			)
			return result, "user@domain"
		}
	}

	// user:pass format
	user, pass, found := strings.Cut(e.Value, ":")
	// avoid detecting NTLM hashes with length check
	if found && !strings.HasPrefix(pass, "//") && len(pass) <= 30 {
		result = append(result,
			Entry{
				Value:   user,
				Comment: e.Comment,
				Tags:    append(e.Tags, "username"),
				Hosts:   e.Hosts,
			},
			Entry{
				Value:   pass,
				Comment: e.Comment,
				Tags:    append(e.Tags, "password"),
				Hosts:   e.Hosts,
			},
		)
		return result, "user:pass"
	}

	return result, ""
}
