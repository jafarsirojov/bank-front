package app

import (
	"context"
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
	// TODO: change to {id} and add own formatting functions
	Post      = "/posts/%s"
	PostEdit  = "/posts/%s/edit"
	Profile   = "/profile"
	Transfer  = "/transfer"
	Payment   = "/payment"
	History   = "/history"
	Register  = "/register"
	AddCard   = "/add/card"
	ErrorPage = "/page/error/client"
	Block     = "/card/block"
	UnBlock   = "/card/unblock"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwtmux.SourceCookie, reflect.TypeOf((*Payload)(nil)).Elem(), s.secret)
	authMWNotOK := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)
	authMWOK := authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, true, Profile)
	s.router.GET(Root, s.handleFrontPage(), authMWOK, jwtMW, logger.Logger("HTTP"))
	// GET -> html

	s.router.GET(ErrorPage, s.handlePageErrorClient(), logger.Logger("HTTP"))
	s.router.POST(ErrorPage, s.handlePageErrorClient(), logger.Logger("HTTP"))

	s.router.GET(Login, s.handleLoginPage(), logger.Logger("HTTP"))
	s.router.GET(Logout, s.handleLogout(), logger.Logger("HTTP"))
	// POST -> form handling + return HTML
	s.router.POST(Login, s.handleLogin(), logger.Logger("HTTP"))
	s.router.GET(Profile, s.handleProfile(), authMWNotOK, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Profile, s.handleProfile(), authMWNotOK, jwtMW, logger.Logger("HTTP"))

	s.router.GET(Transfer, s.handleTransferPage(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(Transfer, s.handleTransfer(), jwtMW, logger.Logger("HTTP"))

	s.router.GET(Block, s.handleBlockPage(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(Block, s.handleBlock(), jwtMW, logger.Logger("HTTP"))

	s.router.GET(UnBlock, s.handleUnBlockPage(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(UnBlock, s.handleUnBlock(), jwtMW, logger.Logger("HTTP"))

	s.router.GET(Register, s.handleRegisterPage(), logger.Logger("HTTP"))
	s.router.POST(Register, s.handleRegister(), logger.Logger("HTTP"))

	s.router.GET(AddCard, s.handleAddCardPage(), logger.Logger("HTTP"))
	s.router.POST(AddCard, s.handleAddCard(), logger.Logger("HTTP"))

	//s.router.GET(Transfer, s.handleHistoryPage(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(History, s.handleHistory(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(History, s.handleHistory(), jwtMW, logger.Logger("HTTP"))

	s.router.GET(Payment, s.handlePayment(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(Payment, s.handlePayment(), jwtMW, logger.Logger("HTTP"))

	//s.router.GET("/cards", s.handleCards(), jwtMW, logger.Logger("HTTP"))
	//s.router.POST("/cards", s.handleCards(), jwtMW, logger.Logger("HTTP"))

	// список постов
	s.router.GET(Posts, s.handlePostsPage(), authMWNotOK, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Posts, s.handlePostsPage(), authMWNotOK, jwtMW, logger.Logger("HTTP"))
	// форма создания/редактирования
	s.router.GET(PostEdit, s.handlePostEditPage(), authMWNotOK, jwtMW, logger.Logger("HTTP"))
	// сохранение
	s.router.POST(PostEdit, s.handlePostEdit(), authMWNotOK, jwtMW, logger.Logger("HTTP"))
}
