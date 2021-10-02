package cbot

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (b *cbot) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	b.errorLog.Output(2, trace)
}

func (b *cbot) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (b *cbot) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = b.session.PopString(r, "flash")
	td.IsAuthenticated = b.IsAuthenticated(r)
	return td
}

func (b *cbot) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := b.templateCache[name]
	if !ok {
		b.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, b.addDefaultData(td, r))
	if err != nil {
		b.serverError(w, err)
	}

	buf.WriteTo(w)
}

func (b *cbot) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
