package gomfk

import (
	"html/template"
)

func Templates() *template.Template {
	return template.Must(template.ParseFiles(
		"mfk/templ/make.html",
		"mfk/templ/header.html",
		"mfk/templ/triple.html",
	))
}
