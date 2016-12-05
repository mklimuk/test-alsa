package annon

import (
	"bytes"
	"fmt"
	html "html/template"
	"log"
	txt "text/template"
)

/*
Templates Wrapper over templates text and instances
*/
type Templates struct {
	Human string `gorethink:"human" json:"human" yaml:"human"`
	Tts   string `gorethink:"tts" json:"tts" yaml:"tts"`
	HTML  string `gorethink:"html" json:"html" yaml:"html"`
	human *txt.Template
	tts   *txt.Template
	html  *html.Template
}

/*
Parse parses arguments present in templates text fields and creates
their templating system equivalents
*/
func (t *Templates) Parse(name string) {
	t.human = txt.Must(txt.New(fmt.Sprintf("Human text template for %v", name)).Parse(t.Human))
	t.tts = txt.Must(txt.New(fmt.Sprintf("Tts text template for %v", name)).Parse(t.Tts))
	t.html = html.Must(html.New(fmt.Sprintf("HTML text template for %v", name)).Parse(t.HTML))
}

//TemplateParams is a struct that holds template variables
type TemplateParams map[string]interface{}

//ExecuteTemplates executes given templates with an appropriate set of parameters
func ExecuteTemplates(tpl *Templates, params TemplateParams, ttsParams TemplateParams) (human string, tts string, html string, err error) {
	human = ExecuteTextTemplate(tpl.human, &params)
	html = ExecuteHTMLTemplate(tpl.html, &params)
	tts = ExecuteTextTemplate(tpl.tts, &ttsParams)
	return human, tts, html, err
}

//ExecuteTextTemplate executes text based template
func ExecuteTextTemplate(tpl *txt.Template, params *TemplateParams) string {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, params)
	if err != nil {
		log.Printf("Could not execute text template: %v", err)
		return ""
	}
	return buf.String()
}

//ExecuteHTMLTemplate executes HTML based template
func ExecuteHTMLTemplate(tpl *html.Template, params *TemplateParams) string {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, params)
	if err != nil {
		log.Printf("Could not execute html template: %v", err)
		return ""
	}
	return buf.String()
}
