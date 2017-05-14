package articles

import (
	"errors"
	"time"

	"github.com/danield21/danield-space/server/repository"
	"github.com/danield21/danield-space/server/repository/categories"
	"golang.org/x/net/context"
)

//FormArticle contains information about articles written on this website.
type FormArticle struct {
	Title       string `schema:"title"`
	Author      string `schema:"author"`
	URL         string `schema:"url"`
	PublishDate string `schema:"publish"`
	Abstract    string `schema:"abstract"`
	Content     string `schema:"content"`
	Category    string `schema:"category"`
}

var ErrNoTitle = errors.New("No title")
var ErrCategoryBadFormat = errors.New("Bad category format")
var ErrURLBadFormat = errors.New("Bad url format")
var ErrNoPublishDate = errors.New("No publish date")
var ErrNoAbstract = errors.New("No abstract")
var ErrNoContent = errors.New("No content")

func (f FormArticle) Unpack(ctx context.Context) (*Article, error) {

	if !repository.ValidURLPart(f.URL) {
		return nil, ErrURLBadFormat
	}

	category, err := parseCategory(ctx, f.Category)
	if err != nil {
		return nil, err
	}

	publish, err := parsePublish(f.PublishDate)
	if err != nil {
		return nil, err
	}

	content, err := parseContent(f.Content)
	if err != nil {
		return nil, err
	}

	a := new(Article)
	a.Title = f.Title
	a.Author = f.Author
	a.URL = f.URL
	a.PublishDate = publish
	a.Abstract = f.Abstract
	a.HTMLContent = content
	a.Category = category

	return a, nil
}

func parseURL(url string) (string, error) {
	if !repository.ValidURLPart(url) {
		return "", ErrURLBadFormat
	}
	return url, nil
}

func parseCategory(ctx context.Context, catURL string) (*categories.Category, error) {
	if !repository.ValidURLPart(catURL) {
		return nil, ErrCategoryBadFormat
	}
	return categories.Get(ctx, catURL)
}

func parsePublish(publish string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", publish)
}

func parseContent(content string) ([]byte, error) {
	html, err := repository.CleanHTML([]byte(content))
	return []byte(html), err
}
