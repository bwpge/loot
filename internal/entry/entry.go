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
	Owned   bool     `json:"owned,omitempty"`
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
	reLMNTLM    = regexp.MustCompile(`^([a-f0-9]{32}):([a-f0-9]{32}):*$`)
	reNetNTLMv2 = regexp.MustCompile(
		`^([^:]+)::([^:]+):[a-fA-F0-9]{16}:[a-fA-F0-9]{32}:[a-fA-F0-9]+$`,
	)
	reNTDSDump = regexp.MustCompile(
		`^([^:]+):\d+:[a-fA-F0-9]{32}:([a-fA-F0-9]{32}):.*$`,
	)
	userPat        = `([a-zA-Z][a-zA-Z0-9._-]{1,29})`
	domainPat      = `([a-zA-Z0-9][a-zA-Z0-9._-]{1,70})`
	reDomain       = regexp.MustCompile("^" + domainPat + "$")
	reUser         = regexp.MustCompile("^" + userPat + "$")
	reIdentityPass = regexp.MustCompile("^" + userPat + "@" + domainPat + ":(.{2,30})$")
)

type detector func(Entry) ([]Entry, string)

var detectors = []detector{
	detectLMNTLM,
	detectNetNTLMv2,
	detectNTDSDump,
	detectIdentityPassword,
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

	return []Entry{e}, ""
}

func derive(parent Entry, value, tag string) Entry {
	return Entry{
		Value:   value,
		Comment: parent.Comment,
		Tags:    append(parent.Tags, tag),
		Hosts:   parent.Hosts,
		Owned:   parent.Owned,
	}
}

func detectLMNTLM(e Entry) ([]Entry, string) {
	m := reLMNTLM.FindStringSubmatch(e.Value)
	if len(m) != 3 {
		return nil, ""
	}

	// usually format will be lm:ntlm, but some odd cases can have nt hash first,
	// sanity check known placeholder value to swap for correct hash
	left := m[1]
	right := m[2]
	nthash := right
	if strings.ToLower(right) == "aad3b435b51404eeaad3b435b51404ee" {
		nthash = left
	}

	return []Entry{
		derive(e, e.Value, "lmntlm"),
		derive(e, nthash, "nt-hash"),
	}, "lm:ntlm hash"
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

func detectNTDSDump(e Entry) ([]Entry, string) {
	m := reNTDSDump.FindStringSubmatch(e.Value)
	if len(m) != 3 {
		return nil, ""
	}

	var nbdomain, user string
	left, right, found := strings.Cut(m[1], `\`)
	if found {
		nbdomain = left
		user = right
	} else {
		user = left
	}

	result := []Entry{
		derive(e, e.Value, "ntds-dump"),
		derive(e, user, "username"),
		derive(e, m[2], "nt-hash"),
	}
	if nbdomain != "" {
		result = append(result,
			derive(e, m[1], "nbuser"),
			derive(e, nbdomain, "nbdomain"),
		)
	}

	return result, "ntds hash"
}

func detectIdentityPassword(e Entry) ([]Entry, string) {
	m := reIdentityPass.FindStringSubmatch(e.Value)
	if len(m) != 4 {
		return nil, ""
	}

	return []Entry{
		derive(e, e.Value, "credential"),
		derive(e, m[1]+"@"+m[2], "identity"),
		derive(e, m[1], "username"),
		derive(e, m[2], "domain"),
		derive(e, m[3], "password"),
	}, "user@domain:pass"
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
