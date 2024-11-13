package models

type TemplateData struct {
	CurrentYear int
	Snippet     *Snippet
	Snippets    []*Snippet
	Form        any
}
