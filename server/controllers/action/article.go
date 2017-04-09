package action

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/danield21/danield-space/server/form"
	"github.com/danield21/danield-space/server/handler"
	"github.com/danield21/danield-space/server/repository"
	"github.com/danield21/danield-space/server/repository/articles"
	"github.com/danield21/danield-space/server/repository/categories"
	"golang.org/x/net/context"
)

const titleKey = "title"
const authorKey = "author"
const urlKey = "url"
const publishKey = "publish"
const abstractKey = "abstract"
const contentKey = "content"
const catKey = "category"

func UnpackArticle(ctx context.Context, values url.Values) (*articles.Article, form.Form) {
	var (
		err         error
		category    *categories.Category
		publishDate time.Time
		content     template.HTML
	)

	frm := form.MakeForm()

	titleFld := frm.AddFieldFromValue(titleKey, values)
	form.NotEmpty(titleFld, "title is required")

	authorFld := frm.AddFieldFromValue(authorKey, values)
	form.NotEmpty(authorFld, "author is required")

	urlFld := frm.AddFieldFromValue(urlKey, values)
	if form.NotEmpty(urlFld, "url is required") && !repository.ValidURLPart(urlFld.Get()) {
		form.Fail(urlFld, "url is not in a proper format")
	}

	catFld := frm.AddFieldFromValue(catKey, values)
	if form.NotEmpty(catFld, "category is required") {
		if !repository.ValidURLPart(catFld.Get()) {
			form.Fail(catFld, "category is not in a proper format")
		} else if category, err = categories.Get(ctx, catFld.Get()); err != nil {
			form.Fail(catFld, "unable to find specified category")
		}
	}

	publishFld := frm.AddFieldFromValue(publishKey, values)
	if form.NotEmpty(publishFld, "publish is required") {
		if publishDate, err = time.Parse("2006-01-02T15:04", publishFld.Get()); err != nil {
			form.Fail(publishFld, "unable to parse time")
		}
	}

	abstractFld := frm.AddFieldFromValue(abstractKey, values)
	form.NotEmpty(abstractFld, "abstract is required")

	contentFld := frm.AddFieldFromValue(contentKey, values)
	if form.NotEmpty(contentFld, "publish is required") {
		if content, err = repository.CleanHTML([]byte(contentFld.Get())); err != nil {
			form.Fail(contentFld, "unable to parse content")
		}
	}

	frm.Submitted = true

	if frm.HasErrors() {
		return nil, frm
	}

	a := new(articles.Article)
	*a = articles.Article{
		Title:       titleFld.Get(),
		Author:      authorFld.Get(),
		Category:    category,
		URL:         urlFld.Get(),
		PublishDate: publishDate,
		Abstract:    abstractFld.Get(),
		HTMLContent: []byte(content),
	}

	return a, frm
}

func PutArticleLink(h handler.Handler) handler.Handler {
	return func(ctx context.Context, e handler.Environment, w http.ResponseWriter) (context.Context, error) {
		r := handler.Request(ctx)
		err := r.ParseForm()
		if err != nil {
			return h(WithForm(ctx, form.Form{Error: errors.New("Unable to parse form")}), e, w)
		}

		art, frm := UnpackArticle(ctx, r.Form)
		if art == nil {
			return h(WithForm(ctx, frm), e, w)
		}

		err = articles.Set(ctx, art)
		if err != nil {
			frm.Error = errors.New("Unable to put into database")
			return h(WithForm(ctx, frm), e, w)
		}

		return h(WithForm(ctx, frm), e, w)
	}
}
