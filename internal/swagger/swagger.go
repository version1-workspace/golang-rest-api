package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

type config struct {
	SchemaURL string
	DomID     string
}

//go:embed swagger.yml
var yml embed.FS

func SwaggerYAML() []byte {
	data, err := yml.ReadFile("swagger.yml")
	if err != nil {
		panic(err)
	}

	return data
}

//go:embed template/*.tmpl
var embededFS embed.FS

func IndexHTML(host string) []byte {
	t, err := template.ParseFS(embededFS, "template/*.tmpl")
	if err != nil {
		panic(err)
	}
	c := config{
		SchemaURL: fmt.Sprintf("%s/api/v1/spec/swagger.yml", host),
		DomID:     "#root",
	}
	var buffer bytes.Buffer
	if err := t.ExecuteTemplate(&buffer, "index.html.tmpl", c); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
