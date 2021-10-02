package cbot

import (
	"database/sql"
	"html/template"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/golangcollege/sessions"
	"github.com/pscompsci/cbot/internal/repository/pg"
)

type contextKey string

const (
	secret                    = "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"
	contextKeyIsAuthenticated = contextKey("isAuthenticated")
)

type cbot struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	users         *pg.UserModel
	templateCache map[string]*template.Template
}

func New(db *sql.DB) *cbot {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	templateCache, err := newTemplateCache("./web/templates/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &cbot{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		users:         &pg.UserModel{DB: db},
		templateCache: templateCache,
	}
}

func (b *cbot) Run() error {
	return b.serve()
}
