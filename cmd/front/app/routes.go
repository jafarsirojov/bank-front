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
	Profile  = "/profile"
	Transfer = "/transfer"
	Payment  = "/payment"
	Register = "/register"
	AddCard  = "/add/card"
	ErrorPage = "/page/error/client"
	Block = "/card/block"
	UnBlock = "/card/unblock"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwtmux.SourceCookie, reflect.TypeOf((*Payload)(nil)).Elem(), s.secret)
	authMW := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)
	s.router.GET(Root, s.handleFrontPage(), logger.Logger("HTTP"))
	// GET -> html

	s.router.GET(ErrorPage, s.handlePageErrorClient(), logger.Logger("HTTP"))
	s.router.POST(ErrorPage, s.handlePageErrorClient(), logger.Logger("HTTP"))

	s.router.GET(Login, s.handleLoginPage(), logger.Logger("HTTP"))
	s.router.GET(Logout, s.handleLogout(), logger.Logger("HTTP"))
	// POST -> form handling + return HTML
	s.router.POST(Login, s.handleLogin(), logger.Logger("HTTP"))
	s.router.GET(Profile, s.handleProfile(), jwtMW, logger.Logger("HTTP"))  //authMW deleted
	s.router.POST(Profile, s.handleProfile(), jwtMW, logger.Logger("HTTP")) //authMW deleted

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

	s.router.GET(Payment, s.handlePayment(), jwtMW, logger.Logger("HTTP"))
	s.router.POST(Payment, s.handlePayment(), jwtMW, logger.Logger("HTTP"))

	s.router.GET("/cards", s.handleCardsPage(), jwtMW, logger.Logger("HTTP"))
	//s.router.GET("/cards", s.handleCards(), jwtMW, logger.Logger("HTTP"))
	s.router.POST("/cards", s.handleCards(), jwtMW, logger.Logger("HTTP"))

	// chat service
	s.router.GET("/message/all", s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/message", s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST("/message", s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
}
