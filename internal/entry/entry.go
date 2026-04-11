package entry

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
	Owner string `json:"owner"`
	Host  string `json:"host"`
}
