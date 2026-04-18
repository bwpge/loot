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
	reNetNTLMv2 = regexp.MustCompile(
		`^([^:]+)::([^:]+):[a-fA-F0-9]{16}:[a-fA-F0-9]{32}:[a-fA-F0-9]+$`,
	)
	reDomain = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{1,70}$`)
	reUser   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]{1,29}$`)
)

type detector func(Entry) ([]Entry, string)

var detectors = []detector{
	detectNTLMLM,
	detectNetNTLMv2,
	detectUserAtDomain,
	detectUserPass,
	detectDomainUser,
}

func DetectValues(e Entry) ([]Entry, string) {
	for _, d := range detectors {
		if result, label := d(e); result != nil {
			return result, label
		}
	}
	return nil, ""
}

func derive(parent Entry, value, tag string) Entry {
	return Entry{
		Value:   value,
		Comment: parent.Comment,
		Tags:    append(parent.Tags, tag),
		Hosts:   parent.Hosts,
	}
}

func detectNTLMLM(e Entry) ([]Entry, string) {
	m := reNTLMLM.FindStringSubmatch(e.Value)
	if len(m) != 2 {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "ntlmlm"),
		derive(e, m[1], "ntlm"),
	}, "ntlm:lm hash"
}

func detectNetNTLMv2(e Entry) ([]Entry, string) {
	m := reNetNTLMv2.FindStringSubmatch(e.Value)
	if len(m) != 3 {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "netntlmv2"),
		derive(e, m[1], "username"),
		derive(e, m[2], "nbdomain"),
	}, "net-ntlmv2"
}

func detectUserAtDomain(e Entry) ([]Entry, string) {
	user, domain, found := strings.Cut(e.Value, "@")
	if !found || !reUser.MatchString(user) || !reDomain.MatchString(domain) {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "identity"),
		derive(e, user, "username"),
		derive(e, domain, "domain"),
	}, "user@domain"
}

func detectUserPass(e Entry) ([]Entry, string) {
	user, pass, found := strings.Cut(e.Value, ":")
	if !found || !reUser.MatchString(user) || strings.HasPrefix(pass, "//") || len(pass) > 30 {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "credential"),
		derive(e, user, "username"),
		derive(e, pass, "password"),
	}, "user:pass"
}

func detectDomainUser(e Entry) ([]Entry, string) {
	domain, user, found := strings.Cut(e.Value, `\`)
	if !found || !reDomain.MatchString(domain) {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "nbuser"),
		derive(e, domain, "nbdomain"),
		derive(e, user, "username"),
	}, "domain\\user"
}
