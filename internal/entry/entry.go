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

	// ntlm:lm
	ntlmlmexp := regexp.MustCompile(`^([a-f0-9]{32}):[a-f0-9]{32}$`)
	matches := ntlmlmexp.FindStringSubmatch(e.Value)
	if len(matches) == 2 {
		result = append(result,
			Entry{
				Value:   matches[1],
				Comment: e.Comment,
				Tags:    append(e.Tags, "ntlm"),
				Hosts:   e.Hosts,
			},
		)
		return result, "ntlm:lm hash"
	}

	// net-ntlmv2
	ntlmv2exp := regexp.MustCompile(
		`^([^:]+)::([^:]+):[a-fA-F0-9]{16}:[a-fA-F0-9]{32}:[a-fA-F0-9]+$`,
	)
	matches = ntlmv2exp.FindStringSubmatch(e.Value)
	if len(matches) == 3 {
		result = append(result,
			Entry{
				Value:   matches[1],
				Comment: e.Comment,
				Tags:    append(e.Tags, "username"),
				Hosts:   e.Hosts,
			},
			Entry{
				Value:   matches[2],
				Comment: e.Comment,
				Tags:    append(e.Tags, "nbdomain"),
				Hosts:   e.Hosts,
			},
		)
		return result, "net-ntlmv2"
	}

	// user@domain
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

	// user:pass
	user, pass, found := strings.Cut(e.Value, ":")
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
