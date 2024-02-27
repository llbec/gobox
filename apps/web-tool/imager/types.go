package imager

type WebSite interface {
	Name() string
	SetUrl(string)
	Images() []string
}
