package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

func redirectErr(w http.ResponseWriter, r *http.Request, slug string, err error) {
	http.Redirect(w, r, fmt.Sprintf("/%s?error=%s", slug, url.QueryEscape(err.Error())), http.StatusSeeOther)
}

func redirectThreadByUUID(w http.ResponseWriter, r *http.Request, slug string, uuid string) {
	http.Redirect(w, r, fmt.Sprintf("/%s/%s", slug, uuid), http.StatusSeeOther)
}

func redirectThreadByID(w http.ResponseWriter, r *http.Request, slug string, id int) {
	http.Redirect(w, r, fmt.Sprintf("/%s/%d", slug, id), http.StatusSeeOther)
}
