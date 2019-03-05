package main

import (
	"encoding/json"
	"encoding/xml"
	"sort"

	"github.com/valyala/fasthttp"

	log "gopkg.in/inconshreveable/log15.v2"
)

func internalError(ctx *fasthttp.RequestCtx, rl log.Logger, err error) {

}

// Category model.
type Category struct {
	ID       int
	Name     string
	Parent   int
	Children []*Category
}

// CategoryMap a map of Category pointers and associates an index with it, for ordered output.
// The order of the index is based on the order items where added.
type CategoryMap struct {
	cats  map[int]*Category
	index []int
}

// NewCm initializes and returs a pointer to a new CategoryMap.
func NewCm() *CategoryMap {
	cats := make(map[int]*Category)
	return &CategoryMap{
		cats: cats,
	}
}

// Sort re-indexes the CategoryMap, increasing order on category id.
// It is advised to add the categories in a sorted way instead of using this method.
func (cm *CategoryMap) Sort() {
	sort.Ints(cm.index)
}

// Set a Category point to the map, only if it was not already added.
// It will also be appeneded to the index, keeping the order of calling this function.
func (cm *CategoryMap) Set(c *Category) {
	if cm.cats[c.ID] != nil {
		return
	}
	cm.cats[c.ID] = c
	cm.index = append(cm.index, c.ID)
}

// Get a Category pointer by its id.
func (cm *CategoryMap) Get(id int) (c *Category) {
	return cm.cats[id]
}

// Index returns a struct of category id. It allows for a sorted range loop.
func (cm *CategoryMap) Index() []int {
	return cm.index
}

// Tree creates decendant tree by population the Category's children.
// It returns a slice of Category pointers, representing the root of the tree.
func (cm *CategoryMap) Tree(offset int) (root []*Category) {
	for _, i := range cm.Index() {
		c := cm.Get(i)
		// Are we at the root of the tree?
		if c.Parent == offset {
			root = append(root, c)
			continue
		}
		p := cm.Get(c.Parent)
		// Append Category to its parent's children
		p.Children = append(p.Children, c)
	}
	return
}

// JSONTree returns a JSON document, representing the parent /
// child relationship of the categories in a tree.
// It returns an error if json.Marshall does.
func (cm *CategoryMap) JSONTree(offset int) ([]byte, error) {
	return json.MarshalIndent(cm.Tree(offset), "", "  ")
}

func (cm *CategoryMap) XMLTree(offset int) ([]byte, error) {
	return xml.MarshalIndent(cm.Tree(offset), "", "  ")
}

func catQuery(offset, depth int) (cm *CategoryMap, err error) {
	rows, err := db.Query("cat-tree-v2", offset, depth)
	if err != nil {
		return
	}
	cm = NewCm()
	for rows.Next() {
		var i, p int
		var n string
		if err = rows.Scan(&i, &n, &p); err != nil {
			return
		}
		c := &Category{
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
