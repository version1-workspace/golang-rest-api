package app

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Route struct {
	Method  string
	Matcher string
	Path    string

	params *Params
}

func (r Route) Params() (Params, error) {
	if r.params != nil {
		return *r.params, nil
	}

	u, err := url.Parse(r.Path)
	if err != nil {
		return nil, err
	}

	p := parseQueryValues(r.Matcher, u)
	r.params = &p

	return p, nil
}

type Params map[string]string

func parseQueryValues(matcher string, u *url.URL) Params {
	params := make(Params)
	key := ""
	value := ""
	mSegments := strings.Split(matcher, "/")
	rSegments := strings.Split(u.Path, "/")
	for i, segment := range mSegments {
		if len(segment) > 0 && segment[0] == '{' && segment[len(segment)-1] == '}' {
			key = segment[1 : len(segment)-1]
			value = rSegments[i]
			params[key] = value
		}
	}

	return params
}

func (p Params) Has(key string) (string, error) {
	v, ok := p[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}
	return v, nil
}

func (p Params) String(key string) (string, error) {
	v, err := p.Has(key)
	if err != nil {
		return "", err
	}
	return v, nil
}

func (p Params) Int(key string) (int, error) {
	v, err := p.Has(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(v)
}

func (p Params) Bool(key string) (bool, error) {
	v, err := p.Has(key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(v)
}
