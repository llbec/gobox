package electric

type Organization struct {
	name      string
	link      string // electric systerm link
	githubURL string // organization github
	subOrgs   map[string]*Organization
}

type ElecInfo struct {
	ArchiveMap map[string]string
	Orgs       map[string]*Organization
	linkMap    map[string]string
}

func NewElecInfo() *ElecInfo {
	archive, err := getContent()
	if err != nil {
		panic(err)
	}
	return &ElecInfo{
		archive,
		make(map[string]*Organization),
		make(map[string]string),
	}
}
