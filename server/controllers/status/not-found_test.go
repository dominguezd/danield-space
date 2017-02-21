package status_test

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine/aetest"

	"github.com/danield21/danield-space/server/handler"
	"github.com/danield21/danield-space/server/handler"
	"github.com/danield21/danield-space/server/handler/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotFound(t *testing.T) {
	client := &http.Client{}

	view := template.New("page/status/not-found")
	view.Parse(", ")
	head := template.New("")
	head.Parse("Hello")
	view.AddParseTree("theme/balloon/head", head.Tree)
	foot := template.New("")
	foot.Parse("World")
	view.AddParseTree("theme/balloon/footer", foot.Tree)

	ctx, done, err := aetest.NewContext()
	require.NoError(t, err, "TestIndexHead.TestIndex - Error in creating context")
	e := handler.TestingEnvironment{Templates: view, Ctx: ctx}
	defer done()

	server := httptest.NewServer(handler.Prepare(e, status.NotFoundHandler))
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL, bytes.NewBuffer(nil))
	require.NoError(t, err, "Error in creating GET request for Index: %v", err)
	request.Header.Add("Content-Type", "text/html")

	response, err := client.Do(request)
	require.NoError(t, err, "Error in creating GET request for Index: %v", err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusNotFound, response.StatusCode, "Expected response status 404, received %s", response.Status)
	assert.Equal(t, "text/html; charset=utf-8", response.Header.Get("Content-Type"))

	assert.NotEmpty(t, response.ContentLength)
	assert.Equal(t, len("Hello, World"), int(response.ContentLength))
}
