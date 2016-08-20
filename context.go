package apimock

import (
	"net/http"
)

type Context interface {
	GetHandlerFunc() http.HandlerFunc
	CheckActions()
}

type contextImpl struct {
	actionManager ActionManager
}

func NewContext(root string) (Context, error) {
	ctx := &contextImpl{}
	actionManager, err := NewActionManager(root)
	if err != nil {
		return nil, err
	}
	ctx.actionManager = actionManager
	return ctx, nil
}

func (ctx *contextImpl) GetHandlerFunc() http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx.handle(w, r)
	}
	return f
}

func (ctx *contextImpl) handle(w http.ResponseWriter, r *http.Request) {
	response, _ := ctx.actionManager.DoAction(r)
	w.WriteHeader(response.Status)
	w.Write(response.Body)
}

func (ctx *contextImpl) CheckActions() {
	ctx.actionManager.CheckActions()
}
