package apidoc

type ApiDocSpec struct {
	Title       string
	Version     string
	Description string
	Scheme      string
	BaseURL     string
	Groups      []*ApiGroupSpec
	Apis        []*ApiSpec
}

type ApiGroupSpec struct {
	Group       string
	Title       string
	Description string
	Apis        []*ApiSpec
}

type ApiSpec struct {
	Title       string
	HTTPMethod  string
	Api         string
	FullURL     string
	Version     string
	Description string
	Author      string
	Deprecated  bool
	Tags        []string
	Group       string
	// Parameters
}
