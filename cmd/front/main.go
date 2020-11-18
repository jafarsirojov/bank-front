package main

import (
	"flag"
	"github.com/jafarsirojov/bank-front/cmd/front/app"
	"github.com/jafarsirojov/bank-front/pkg/core/auth"
	"github.com/jafarsirojov/bank-front/pkg/core/cards"
	"github.com/jafarsirojov/bank-front/pkg/core/history"
	"github.com/jafarsirojov/bank-front/pkg/jwt"
	"github.com/jafarsirojov/bank-front/pkg/mux"
	"net"
	"net/http"
)

var (
	host       = flag.String("host", "", "Server host")
	port       = flag.String("port", "", "Server port")
	authUrl    = flag.String("authUrl", "", "Auth Service URL")
	cardsUrl   = flag.String("cardsUrl", "", "Cards Service URL")
	historyUrl = flag.String("historyUrl", "", "Transfer Service URL")
)

	//-host 0.0.0.0 -port 9012 -authUrl "http://localhost:9011" -cardsUrl "http://localhost:9009" -historyUrl "http://localhost:9010"

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	secret := jwt.Secret("top secret")
	start(addr, secret, auth.Url(*authUrl), cards.Url(*cardsUrl), history.Url(*historyUrl))
}

func start(addr string, secret jwt.Secret, authURL auth.Url, cardsURL cards.Url, historyURL history.Url) {
	exactMux := mux.NewExactMux()
	authSvc := auth.NewClient(authURL)
	cardsSvc := cards.NewCard(cardsURL)
	historySvc := history.NewHistory(historyURL)
	server := app.NewServer(exactMux, secret, authSvc, cardsSvc, historySvc)
	server.Start()

	panic(http.ListenAndServe(addr, server))
}
