package cbot

import (
	"errors"
	"net/http"

	"github.com/pscompsci/cbot/internal/repository"
	"github.com/pscompsci/cbot/pkg/forms"
)

func (b *cbot) signupUserForm(w http.ResponseWriter, r *http.Request) {
	b.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (b *cbot) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		b.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 8)

	if !form.Valid() {
		b.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = b.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in user")
			b.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			b.serverError(w, err)
		}
		return
	}

	b.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (b *cbot) loginUserForm(w http.ResponseWriter, r *http.Request) {
	b.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (b *cbot) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		b.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := b.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			b.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else if errors.Is(err, repository.ErrUserNotActivated) {
			form.Errors.Add("generic", "User account not activated")
			b.render(w, r, "signin.page.tmpl", &templateData{Form: form})
		} else {
			b.serverError(w, err)
		}
		return
	}

	b.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (b *cbot) logoutUser(w http.ResponseWriter, r *http.Request) {
	b.session.Remove(r, "authenticatedUserID")
	b.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
