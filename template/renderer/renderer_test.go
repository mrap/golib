package renderer

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var pwd, _ = os.Getwd()

func getBodyString(u string) string {
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func TestRenderNestedTemplate(t *testing.T) {
	// Setup test templates
	templatesPath := filepath.Join(pwd, "test_templates")
	SetupTemplates(templatesPath)

	// Start test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/nested":
			RenderTemplate(w, "nested/test", nil)
		default:
			RenderTemplate(w, "test", nil)
		}
	}))
	defer ts.Close()

	// Test responses to make sure we get the correct template
	body := getBodyString(ts.URL)
	if strings.TrimSpace(body) != "test.html" {
		t.Error("%s should equal %s", body, "test.html")
	}

	body = getBodyString(ts.URL + "/nested")
	if strings.TrimSpace(body) != "nested/test.html" {
		t.Error("%s should equal %s", body, "nested/test.html")
	}

}
