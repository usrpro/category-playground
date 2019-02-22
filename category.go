package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

func logHTTPError(w http.ResponseWriter, r *http.Request, err error, code int, info ...string) {
	var f func(string, ...interface{})
	switch {
	case code < 300:
		f = log.Debug
	case code < 400:
		f = log.Info
	case code < 500:
		f = log.Warn
	default:
		f = log.Error
	}
	guru := strconv.Itoa(rand.Int())
	info = append(info, "Guru meditation:", guru)
	msg := strings.Join(info, " ")
	f(err.Error(), "code", code, "uri", r.RequestURI, "guru", guru, "msg", msg)
	http.Error(w, msg, code)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	log.Warn("Bad request URI", "uri", r.RequestURI)
	http.Error(w, "Invalid URI format", http.StatusBadRequest)
}

type categoryMap map[int]*category

func (cm categoryMap) sort() (index []int) {
	for k := range cm {
		index = append(index, k)
	}
	sort.Ints(index)
	return
}

type category struct {
	ID       int
	Data     string
	Parent   int
	Children []*category
}

func catHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Request", "uri", r.RequestURI)
	a := time.Now().UnixNano()
	// Still hard-coded, needs to be taken from request params.
	offset, depth := 0, 6
	rows, err := db.Query("cat-tree-v2", offset, depth)
	if err != nil {
		logHTTPError(w, r, err, http.StatusInternalServerError)
		return
	}
	cm := make(categoryMap)
	for rows.Next() {
		var i, p int
		var d string
		if err := rows.Scan(&i, &d, &p); err != nil {
			logHTTPError(w, r, err, http.StatusInternalServerError)
			return
		}
		c := &category{
			ID:     i,
			Data:   d,
			Parent: p,
		}
		cm[c.ID] = c
	}
	root := []*category{}
	for _, i := range cm.sort() {
		c := cm[i]
		// Are we at the root of the tree?
		if c.Parent == offset {
			root = append(root, c)
			continue
		}
		p := cm[c.Parent]
		// Append category to its parent's children
		p.Children = append(p.Children, c)
	}
	js, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		logHTTPError(w, r, err, http.StatusInternalServerError)
		return
	}
	dt := float64((time.Now().UnixNano() - a)) / 1000000 // ms
	log.Info("Request completed", "uri", r.RequestURI, "ms", dt)
	w.Write(js)
}

func bcHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Request", "uri", r.RequestURI)
	if r.RequestURI == "/bc/" || r.RequestURI == "/bc/0" {
		logHTTPError(w, r, errors.New("No content"), http.StatusNoContent, "No content")
		return
	}
	p := strings.Split(r.RequestURI, "/")
	log.Debug("Parsed", "path", p, "level", len(p))
	cid, err := strconv.Atoi(p[len(p)-1])
	if err != nil {
		logHTTPError(w, r, err, http.StatusBadRequest, "Bad request")
		return
	}
	row, err := db.QueryRow("bc-json", cid)
	if err != nil {
		logHTTPError(w, r, err, http.StatusInternalServerError)
		return
	}
	var j string
	if err := row.Scan(&j); err != nil {
		logHTTPError(w, r, err, http.StatusBadRequest, "Category id doesn't exist")
		return
	}
	w.Write([]byte(j))
}
