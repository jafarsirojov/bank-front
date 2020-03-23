package app

import (
	"context"
	"errors"
	"github.com/jafarsirojov/bank-front/pkg/core/auth"
	"github.com/jafarsirojov/bank-front/pkg/core/cards"
	"github.com/jafarsirojov/bank-front/pkg/core/history"
	"github.com/jafarsirojov/bank-front/pkg/core/utils"
	"github.com/jafarsirojov/bank-front/pkg/jwt"
	"github.com/jafarsirojov/bank-front/pkg/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Server struct {
	router     *mux.ExactMux
	secret     jwt.Secret
	authSvc    *auth.Client
	cardsSvc   *cards.Card
	historySvc *history.History
}

func NewServer(router *mux.ExactMux, secret jwt.Secret, authSvc *auth.Client, cardsSvc *cards.Card, historySvc *history.History) *Server {
	return &Server{router: router, secret: secret, authSvc: authSvc, cardsSvc: cardsSvc, historySvc: historySvc}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleFrontPage() http.HandlerFunc {
	// executes in one goroutine
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "index.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		// executes in many goroutines
		// TODO: fetch data from multiple upstream services
		err := tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cookie := &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleLoginPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "login.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "login.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			return
		}
		login := request.PostFormValue("login")
		if login == "" {
			// TODO: show error page
			log.Print("login can't be empty")
			return
		}
		password := request.PostFormValue("password")
		if password == "" {
			// TODO: show error page
			log.Print("password can't be empty")
			return
		}

		token, err := s.authSvc.Login(request.Context(), login, password)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				// TODO: show error page (for deadline)
				log.Print("auth service didn't response in given time")
				log.Print("another err") // parse it
			case errors.Is(err, context.Canceled):
				// TODO: show error page (for deadline)
				log.Print("auth service didn't response in given time")
				log.Print("another err") // parse it
			case errors.Is(err, auth.ErrResponse):
				var typedErr *auth.ErrorResponse
				ok := errors.As(err, &typedErr)
				if ok {
					tplData := struct {
						Err string
					}{
						Err: "",
					}
					// TODO: work with another
					if utils.StringInSlice("err.password_mismatch", typedErr.Errors) {
						tplData.Err = "err.password_mismatch"
					}

					err := tpl.Execute(writer, tplData)
					if err != nil {
						log.Print(err)
					}
				}
			}
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handlePostsPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handlePostEditPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handlePostEdit() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handleProfile() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "profile.gohtml"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

		allCards, err := s.cardsSvc.AllCards(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				log.Print("auth service didn't response in given time")
				log.Print("another err")
				http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
			case errors.Is(err, context.Canceled):
				log.Print("auth service didn't response in given time")
				log.Print("another err")
				http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
			case errors.Is(err, auth.ErrResponse):
				var typedErr *auth.ErrorResponse
				ok := errors.As(err, &typedErr)
				if ok {
					tplData := struct {
						Err string
					}{
						Err: "",
					}
					if utils.StringInSlice("err.password_mismatch", typedErr.Errors) {
						tplData.Err = "err.password_mismatch"
					}

					err := tpl.Execute(writer, tplData)
					if err != nil {
						log.Print(err)
					}
				}
			}
			return
		}
		err = tpl.Execute(writer, allCards)
		if err != nil {
			log.Printf("can't execute: %d", err)
			http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
	}
}
