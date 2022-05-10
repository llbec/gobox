package electric

import "fmt"

func (o *Organization) SetName(name string) {
	o.name = name
}

func (o *Organization) SetLink(url string) {
	o.link = url
}

func (o *Organization) SetGithub(url string) {
	o.githubURL = url
}

func (o *Organization) AddSub(name string) {
	if o.subOrgs == nil {
		o.subOrgs = make(map[string]*Organization)
	}
	org, err := GetOrg(name)
	if err != nil {
		fmt.Println("AddSub:", err)
	}
	o.subOrgs[name] = org
}
