package app

type Route struct {
	Method  string
	Path    string
	Matcher string

	params *Params
}

type Params map[string]string
