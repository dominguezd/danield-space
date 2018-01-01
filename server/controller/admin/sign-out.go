package admin

import (
	"context"
	"html/template"
	"net/http"

	"github.com/danield21/danield-space/server/controller"
	"github.com/danield21/danield-space/server/store"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

type SignOutController struct {
	Renderer            controller.Renderer
	SiteInfo            store.SiteInfoRepository
	SignOut             Processor
	InternalServerError controller.Controller
}

func (ctr SignOutController) Serve(ctx context.Context, pg *controller.Page, rqs *http.Request) controller.Controller {
	info := ctr.SiteInfo.Get(ctx)

	if rqs.Method == http.MethodPost {
		ctr.SignOut.Process(ctx, rqs, pg.Session)
	}

	cnt, err := ctr.Renderer.String("page/admin/sign-out", nil)

	if err != nil {
		log.Errorf(ctx, "%v", errors.Wrap(err, "unable to render content"))
		return ctr.InternalServerError
	}

	pg.Title = info.Title
	pg.Status = http.StatusSeeOther
	pg.Header["Location"] = "/"
	pg.Meta["description"] = info.ShortDescription()
	pg.Meta["author"] = info.Owner
	pg.Content = template.HTML(cnt)

	return nil
}