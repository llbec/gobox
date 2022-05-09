package govcncode

//AddressCode defines a interface to get the Administrative division Code
type AddressCode interface {
	GetProvinceList() []*Province
}

//Address defines a interface for all the address structs
type Address interface {
	GetName() string
	GetCode() int
}

//Console defines a inteface for a simple console shell
type Console interface {
	Run()
	Stutus() string
}

//struct

//County defines a struct for county
type County struct {
	Code   int
	Parent *City
	Name   string
}

//GetCode for join address interface
func (a *County) GetCode() int {
	return a.Code
}

//GetName for join address interface
func (a *County) GetName() string {
	return a.Name
}

//City defines a struct for City
type City struct {
	Code    int
	Countys []*County
	Parent  *Province
	Name    string
}

//GetCode for join address interface
func (a *City) GetCode() int {
	return a.Code
}

//GetName for join address interface
func (a *City) GetName() string {
	return a.Name
}

//GetCounty search county object by name
func (a *City) GetCounty(name string) *County {
	for _, v := range a.Countys {
		if v.Name == name {
			return v
		}
	}
	return nil
}

//Province defines a struct for Province
type Province struct {
	Code  int
	Citys []*City
	Name  string
}

//GetCode for join address interface
func (a *Province) GetCode() int {
	return a.Code
}

//GetName for join address interface
func (a *Province) GetName() string {
	return a.Name
}

//GetCity search city object by name
func (a *Province) GetCity(name string) *City {
	for _, v := range a.Citys {
		if v.Name == name {
			return v
		}
	}
	return nil
}
