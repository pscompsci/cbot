package cbot

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (b *cbot) routes() http.Handler {
	standardMiddleware := alice.New(b.recoverPanic, b.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(b.session.Enable, b.authenticate)

	router := pat.New()

	router.Get("/", dynamicMiddleware.ThenFunc(b.loginUserForm))

	router.Post("/user/login", dynamicMiddleware.ThenFunc(b.loginUser))
	router.Get("/user/signup", dynamicMiddleware.ThenFunc(b.signupUserForm))
	router.Post("/user/signup", dynamicMiddleware.ThenFunc(b.signupUser))
	router.Post("/user/logout", dynamicMiddleware.Append(b.requireAuthentication).ThenFunc(b.logoutUser))

	router.Get("/admin", dynamicMiddleware.Append(b.requireAuthentication).Then(b.admin()))

	fileserver := http.FileServer(http.Dir("./web/static/"))
	router.Get("/static/", http.StripPrefix("/static", fileserver))

	return standardMiddleware.Then(router)
}
