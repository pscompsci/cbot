package cbot

import "net/http"

func (b *cbot) admin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.render(w, r, "admin.page.tmpl", nil)
	}
}
