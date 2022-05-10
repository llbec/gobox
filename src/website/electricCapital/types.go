package electric

type Organization struct {
	name      string
	link      string // electric systerm link
	githubURL string // organization github
	subOrgs   map[string]*Organization
}
