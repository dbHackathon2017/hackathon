package main

import (
	"fmt"
	"net/http"
)

var _ = fmt.Sprintf("")

func HandleIndexPage(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.Delims(templateDelims[0], templateDelims[1]).ExecuteTemplate(w, "index", nil)
	return nil
}

func HandleNotFoundError(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "notFoundError", nil)
	return nil
}
