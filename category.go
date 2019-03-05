package main

import (
	"github.com/usrpro/cats"
	"github.com/valyala/fasthttp"

	log "gopkg.in/inconshreveable/log15.v2"
)

func internalError(ctx *fasthttp.RequestCtx, rl log.Logger, err error) {

}

func catQuery(offset, depth int) (cm *cats.CategoryMap, err error) {
	rows, err := db.Query("cat-tree-v2", offset, depth)
	if err != nil {
		return
	}
	cm = cats.NewCm()
	for rows.Next() {
		var i, p int
		var n string
		if err = rows.Scan(&i, &n, &p); err != nil {
			return
		}
		c := &cats.Category{
			ID:     i,
			Name:   n,
			Parent: p,
		}
		cm.Set(c)
	}
	return
}

func catHandler(ctx *fasthttp.RequestCtx, rl log.Logger) {
	a := ctx.QueryArgs()
	offset, depth := a.GetUintOrZero("offset"), a.GetUintOrZero("depth")
	if depth == 0 {
		rl.Warn("Requested depth of 0 or emtpy in catHandler")
		ctx.Error("depth cannot be 0", fasthttp.StatusBadRequest)
		return
	}
	cm, err := catQuery(offset, depth)
	if err != nil {
		rl.Error(err.Error())
		ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
		return
	}
	js, err := cm.JSONTree(offset)
	if err != nil {
		rl.Error(err.Error())
		ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Write(js)
}

func bcHandler(ctx *fasthttp.RequestCtx, rl log.Logger) {
	a := ctx.QueryArgs()
	cid := a.GetUintOrZero("cat-id")
	if cid == 0 {
		rl.Info("No content in bcHandler")
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
	row, err := db.QueryRow("bc-json", cid)
	if err != nil {
		internalError(ctx, rl, err)
		return
	}
	var js []byte
	if err := row.Scan(&js); err != nil {
		rl.Warn("Category id does not exist in bcHandler")
		ctx.Error("Category id does not exist", fasthttp.StatusBadRequest)
		return
	}
	ctx.Write(js)
}
