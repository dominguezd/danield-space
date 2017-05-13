package server

import (
	"html/template"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

//ProductionEnvironment contains all the data required to run the server
type ProductionEnvironment struct {
	Templates         *template.Template
	GenerateTemplates <-chan *template.Template
	WaitForView       sync.Mutex
}

//View generates a view based on the templates stored
func (p *ProductionEnvironment) View(w io.Writer, theme, view string, data interface{}) error {
	p.WaitForView.Lock()
	if p.Templates == nil {
		p.Templates = <-p.GenerateTemplates
	}
	p.WaitForView.Unlock()

	return RenderTemplateWithTheme(p.Templates, w, view, data)
}

//Session gets the session using a secure key
func (p *ProductionEnvironment) Session(r *http.Request) (session *sessions.Session) {
	session = GetSession(r)
	return
}

//Context gets a new app engine context
func (p *ProductionEnvironment) Context(r *http.Request) (ctx context.Context) {
	ctx = appengine.NewContext(r)
	return
}
