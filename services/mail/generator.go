package mail

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type templateGenerator interface {
	generateEmailContent(templateSource string, model interface{}) (string, error)
}
type fileTemplateGenerator struct{}

func (fileTemplateGenerator) generateEmailContent(templateSource string, model interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateSource)
	if err != nil {
		log.Errorf("Error loading template %v", err)
		return "", err
	}
	tmpl.Option("missingkey=error")
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, model)
	if err != nil {
		log.Errorf("Error generating content %v", err)
		return "", err
	}
	return tpl.String(), nil
}
