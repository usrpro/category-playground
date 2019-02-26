package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/valyala/fasthttp"

	log "gopkg.in/inconshreveable/log15.v2"
)

type cmdLine struct {
	test                       bool
	pgHost, pgUser, pgDatabase string
}

var flags cmdLine

func parseFlags() {
	flag.BoolVar(&flags.test, "test", false, "Load the test dataset")
	flag.StringVar(&flags.pgHost, "pg-host", "/run/postgresql", "Postgresql hostname")
	flag.StringVar(&flags.pgUser, "pg-user", "postgres", "Postgresql username")
	flag.StringVar(&flags.pgDatabase, "pg-db", "playground", "Postgresql database")
	flag.Parse()
	log.Debug("Starting with arguments:", "flags", flags)
}

func initLog() {
	rand.Seed(time.Now().UnixNano())
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	rl := log.New("ConnID", ctx.ConnID(), "ReqID", ctx.ID())
	rl.Debug("Request start", "uri", string(ctx.RequestURI()))
	h := "none"
	switch string(ctx.Path()) {
	case "/bc":
		h = "bcHandler"
		bcHandler(ctx, rl)
	case "/cat":
		h = "catHandler"
		catHandler(ctx, rl)
	default:
		ctx.NotFound()
	}
	dt := float64((time.Now().UnixNano() - ctx.Time().UnixNano())) / 1000000 // ms
	rl.Info("Request handled", "handler", h, "duration", dt)
}

func main() {
	parseFlags()
	initDb()
	if err := fasthttp.ListenAndServe("0.0.0.0:8012", requestHandler); err != nil {
		log.Warn("ListenAndServe terminated", "error", err)
	}
}
