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

var (
	reNTLMLM    = regexp.MustCompile(`^([a-f0-9]{32}):[a-f0-9]{32}$`)
	reNetNTLMv2 = regexp.MustCompile(`^([^:]+)::([^:]+):[a-fA-F0-9]{16}:[a-fA-F0-9]{32}:[a-fA-F0-9]+$`)
	reDomain    = regexp.MustCompile(`^[a-zA-Z0-9._-]{2,70}$`)
)

func DetectValues(e Entry) ([]Entry, string) {
	var result []Entry

	matches := reNTLMLM.FindStringSubmatch(e.Value)
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

	matches = reNetNTLMv2.FindStringSubmatch(e.Value)
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

	user, domain, found := strings.Cut(e.Value, "@")
	if found {
		// not a good domain regex, but good enough for guessing
		if reDomain.MatchString(domain) {
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
