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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		login := request.PostFormValue("login")
		if login == "" {
			// TODO: show error page
			log.Print("login can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		password := request.PostFormValue("password")
		if password == "" {
			// TODO: show error page
			log.Print("password can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
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

type Auth struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Expired int64  `json:"exp"`
}

func (s *Server) handleProfile() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "profile.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print("start handle profile  2")
		ctx, _ := context.WithTimeout(context.Background(), 210*time.Second)
		token, err := request.Cookie("token")
		///*authentication*/_, ok := jwt2.FromContext(request.Context()).(*Auth)
		//if !ok {
		//	log.Print("can't authentication is not ok")
		//	http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
		//	return
		//}
		//authentication.Id==0
		allCards, err := s.cardsSvc.AllCards(ctx, token.Value)

		log.Print("start handle profile  3")

		if err != nil {
			log.Printf("error------------------------------------------ : %s", err)
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
		log.Print("start handle profile  2")
		err = tpl.Execute(writer, allCards)
		log.Print("start handle profile  2")
		if err != nil {
			log.Printf("can't execute: %d", err)
			http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
			return
		}
		//http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleTransfer() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "transfer.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		numberCard := request.PostFormValue("numberCard")
		if numberCard == "" {
			// TODO: show error page
			log.Print("numberCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		idCard := request.PostFormValue("idCard")
		if idCard == "" {
			// TODO: show error page
			log.Print("idCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		count := request.PostFormValue("count")
		if count == "" {
			// TODO: show error page
			log.Print("count can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		token, err := request.Cookie("token")
		if err != nil {
			log.Print("can't token in cookie")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		err = s.cardsSvc.Transfer(request.Context(), numberCard, idCard,count, token.Value)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handlePayment() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		_   *template.Template //tpl *template.Template
		err error
	)
	_, err = template.ParseFiles(filepath.Join("web/templates", "payment.gohtml")) //tpl
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		//	log.Print("start handle profile  2")
		//	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		//
		//	//allCards, err := s.cardsSvc.AllCards(ctx)
		//	log.Print("start handle profile  3")
		//
		//	if err != nil {
		//		log.Printf("error------------------------------------------ : %s", err)
		//		switch {
		//		case errors.Is(err, context.DeadlineExceeded):
		//			log.Print("auth service didn't response in given time")
		//			log.Print("another err")
		//			http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
		//		case errors.Is(err, context.Canceled):
		//			log.Print("auth service didn't response in given time")
		//			log.Print("another err")
		//			http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
		//		case errors.Is(err, auth.ErrResponse):
		//			var typedErr *auth.ErrorResponse
		//			ok := errors.As(err, &typedErr)
		//			if ok {
		//				tplData := struct {
		//					Err string
		//				}{
		//					Err: "",
		//				}
		//				if utils.StringInSlice("err.password_mismatch", typedErr.Errors) {
		//					tplData.Err = "err.password_mismatch"
		//				}
		//
		//				err := tpl.Execute(writer, tplData)
		//				if err != nil {
		//					log.Print(err)
		//				}
		//			}
		//		}
		//		return
		//	}
		//	log.Print("start handle profile  2")
		//	err = tpl.Execute(writer, allCards)
		//	log.Print("start handle profile  2")
		//	if err != nil {
		//		log.Printf("can't execute: %d", err)
		//		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
		//		return
		//	}
		//	//http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
	}
}

//func (s *Server) handleCards() http.HandlerFunc {
//	log.Print("start handle profile")
//	var (
//		tpl *template.Template
//		err error
//	)
//	tpl, err = template.ParseFiles(filepath.Join("web/templates", "profile.gohtml"))
//	if err != nil {
//		log.Printf("-----------------------------------%s", err)
//		panic(err)
//	}
//	return func(writer http.ResponseWriter, request *http.Request) {
//		log.Print("start handle profile  2")
//		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
//		asd := ""
//		allCards, err := s.cardsSvc.AllCards(ctx, asd)
//		log.Print("start handle profile  3")
//
//		if err != nil {
//			log.Printf("error------------------------------------------ : %s", err)
//			switch {
//			case errors.Is(err, context.DeadlineExceeded):
//				log.Print("auth service didn't response in given time")
//				log.Print("another err")
//				http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
//			case errors.Is(err, context.Canceled):
//				log.Print("auth service didn't response in given time")
//				log.Print("another err")
//				http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
//			case errors.Is(err, auth.ErrResponse):
//				var typedErr *auth.ErrorResponse
//				ok := errors.As(err, &typedErr)
//				if ok {
//					tplData := struct {
//						Err string
//					}{
//						Err: "",
//					}
//					if utils.StringInSlice("err.password_mismatch", typedErr.Errors) {
//						tplData.Err = "err.password_mismatch"
//					}
//
//					err := tpl.Execute(writer, tplData)
//					if err != nil {
//						log.Print(err)
//					}
//				}
//			}
//			return
//		}
//		log.Print("start handle profile  2")
//		err = tpl.Execute(writer, allCards)
//		log.Print("start handle profile  2")
//		if err != nil {
//			log.Printf("can't execute: %d", err)
//			http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
//			return
//		}
//		//http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
//	}
//}

func (s *Server) handleTransferPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "transfer.gohtml"))
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

func (s *Server) handleHistory() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "history.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print("start handle profile  2")
		ctx, _ := context.WithTimeout(context.Background(), 210*time.Second)
		token, err := request.Cookie("token")
		if err != nil {
			log.Printf("can't token is nil: %d",err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		///*authentication*/_, ok := jwt2.FromContext(request.Context()).(*Auth)
		//if !ok {
		//	log.Print("can't authentication is not ok")
		//	http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
		//	return
		//}
		//authentication.Id==0
		AllHistory, err := s.historySvc.AllHistory(ctx, token.Value)

		log.Print("start handle profile  3")

		if err != nil {
			log.Printf("error------------------------------------------ : %s", err)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		log.Print("start handle profile  2")

		tplData2 := struct {
			Data []history.ModelOperationsLog
		}{
			Data: AllHistory,
		}
		err = tpl.Execute(writer, tplData2)
		log.Print("start handle profile  2")
		if err != nil {
			log.Printf("can't execute: %d", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleRegisterPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "register.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "register.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		name := request.PostFormValue("name")
		if name == "" {
			// TODO: show error page
			log.Print("numberCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		login := request.PostFormValue("login")
		if login == "" {
			// TODO: show error page
			log.Print("idCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		password := request.PostFormValue("password")
		if password == "" {
			// TODO: show error page
			log.Print("count can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		phone := request.PostFormValue("phone")
		if phone == "" {
			// TODO: show error page
			log.Print("count can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		err = s.authSvc.Register(request.Context(), name, login,password, phone)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Login, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleAddCardPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "addcard.gohtml"))
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

func (s *Server) handleAddCard() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "addcard.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		name := request.PostFormValue("name")
		if name == "" {
			// TODO: show error page
			log.Print("name can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		balance := request.PostFormValue("balance")
		if balance == "" {
			// TODO: show error page
			log.Print("balance can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		ownerid := request.PostFormValue("ownerid")
		if ownerid == "" {
			// TODO: show error page
			log.Print("ownerid can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		token, err := request.Cookie("token")

		err = s.cardsSvc.AddCard(request.Context(), name, balance, ownerid, token.Value)
		if err != nil {
			log.Printf("can't add card %s",err)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handlePageErrorClient() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "errorclient.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
			return
		}
		//http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleBlockPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "block.gohtml"))
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

func (s *Server) handleBlock() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "block.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		numberCard := request.PostFormValue("numberCard")
		if numberCard == "" {
			// TODO: show error page
			log.Print("numberCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		idCard := request.PostFormValue("idCard")
		if idCard == "" {
			// TODO: show error page
			log.Print("idCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		count := request.PostFormValue("count")
		if count == "" {
			// TODO: show error page
			log.Print("count can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		token, err := request.Cookie("token")
		if err != nil {
			log.Print("can't token in cookie")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		err = s.cardsSvc.Transfer(request.Context(), numberCard, idCard,count, token.Value)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleUnBlockPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "unblock.gohtml"))
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

func (s *Server) handleUnBlock() http.HandlerFunc {
	log.Print("start handle profile")
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "unblock.gohtml"))
	if err != nil {
		log.Printf("-----------------------------------%s", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// TODO: show error page
			log.Printf("error while parse login form: %v", err)
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		numberCard := request.PostFormValue("numberCard")
		if numberCard == "" {
			// TODO: show error page
			log.Print("numberCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		idCard := request.PostFormValue("idCard")
		if idCard == "" {
			// TODO: show error page
			log.Print("idCard can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		count := request.PostFormValue("count")
		if count == "" {
			// TODO: show error page
			log.Print("count can't be empty")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		token, err := request.Cookie("token")
		if err != nil {
			log.Print("can't token in cookie")
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}

		err = s.cardsSvc.Transfer(request.Context(), numberCard, idCard,count, token.Value)
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
			http.Redirect(writer, request, ErrorPage, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(writer, request, Profile, http.StatusTemporaryRedirect)
	}
}
