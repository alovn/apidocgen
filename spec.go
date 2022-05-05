package apidoc

type ApiDocSpec struct {
	Title       string
	Version     string
	Description string
	Host        string
	BasePath    string
	Apis        []ApiSpec
}

type ApiSpec struct {
	Title       string
	HTTPMethod  string
	Api         string
	Version     string
	Description string
	Author      string
	Deprecated  bool
	Tags        []string
	// Parameters
}
