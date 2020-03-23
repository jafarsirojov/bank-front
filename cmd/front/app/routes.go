package app

import (
	"github.com/jafarsirojov/bank-front/pkg/mux/middleware/authenticated"
	"github.com/jafarsirojov/bank-front/pkg/mux/middleware/jwt"
	jwtmux "github.com/jafarsirojov/bank-front/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/bank-front/pkg/mux/middleware/logger"
	"reflect"
)

var (
	Root   = "/"
	Login  = "/login"
	Logout = "/logout"
	Posts  = "/posts"
	// TODO: change to {id} and add own formattin functions
	Post     = "/posts/%s"
	PostEdit = "/posts/%s/edit"
	Profile  = "/profile"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwtmux.SourceCookie, reflect.TypeOf((*Payload)(nil)).Elem(), s.secret)
	authMW := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)
	s.router.GET(Root, s.handleFrontPage(), jwtMW, logger.Logger("HTTP"))
	// GET -> html
	s.router.GET(Login, s.handleLoginPage(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(Logout, s.handleLogout(), jwtMW, logger.Logger("HTTP"))
	// POST -> form handling + return HTML
	s.router.POST(Login, s.handleLogin(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(Profile, s.handleProfile(), jwtMW, authMW, logger.Logger("HTTP"))

	// список постов
	s.router.GET(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	// форма создания/редактирования
	s.router.GET(PostEdit, s.handlePostEditPage(), authMW, jwtMW, logger.Logger("HTTP"))
	// сохранение
	s.router.POST(PostEdit, s.handlePostEdit(), authMW, jwtMW, logger.Logger("HTTP"))
}
