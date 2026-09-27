package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVanityHandlerGoGet(t *testing.T) {
	handler := vanityHandler(
		"poly.red git https://github.com/changkun/polyred",
		"poly.red https://github.com/changkun/polyred https://github.com/changkun/polyred/tree/main{/dir} https://github.com/changkun/polyred/blob/main{/dir}/{file}#L{line}",
		"https://github.com/changkun/polyred",
	)

	req := httptest.NewRequest(http.MethodGet, "https://poly.red/render?go-get=1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	content := string(body)
	for _, want := range []string{
		`<meta name="go-import" content="poly.red git https://github.com/changkun/polyred">`,
		`<meta name="go-source" content="poly.red https://github.com/changkun/polyred https://github.com/changkun/polyred/tree/main{/dir} https://github.com/changkun/polyred/blob/main{/dir}/{file}#L{line}">`,
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("response body missing %q:\n%s", want, content)
		}
	}
}

func TestVanityHandlerRedirect(t *testing.T) {
	handler := vanityHandler(
		"poly.red git https://github.com/changkun/polyred",
		"",
		"https://github.com/changkun/polyred",
	)

	req := httptest.NewRequest(http.MethodGet, "https://poly.red/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusMovedPermanently)
	}
	if got := res.Header.Get("Location"); got != "https://github.com/changkun/polyred" {
		t.Fatalf("Location = %q, want %q", got, "https://github.com/changkun/polyred")
	}
}

func TestCountedHandler(t *testing.T) {
	var reported []string
	handler := countedHandler("/research/", "https://changkun.de", func(page, ua string) {
		reported = append(reported, page+" "+ua)
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	for _, path := range []string{
		"/research/",
		"/research/index.html",
		"/research/talks/phd.pdf",
		"/research/papers/",
		"/research",
		"/blog/",
	} {
		req := httptest.NewRequest(http.MethodGet, "https://www.changkun.de"+path, nil)
		req.Header.Set("User-Agent", "reader")
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}

	want := []string{
		"https://changkun.de/research/talks/phd.pdf reader",
		"https://changkun.de/research/papers/ reader",
	}
	if strings.Join(reported, "\n") != strings.Join(want, "\n") {
		t.Fatalf("reported:\n%s\nwant:\n%s", strings.Join(reported, "\n"), strings.Join(want, "\n"))
	}
}

func TestReportUrlstat(t *testing.T) {
	got := make(chan *http.Request, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r
	}))
	defer srv.Close()
	old := urlstatEndpoint
	urlstatEndpoint = srv.URL + "/urlstat?report=page+site"
	defer func() { urlstatEndpoint = old }()

	reportUrlstat("https://changkun.de/research/talks/phd.pdf", "reader")

	r := <-got
	if r.URL.Query().Get("report") != "page site" {
		t.Errorf("report = %q, want %q", r.URL.Query().Get("report"), "page site")
	}
	if h := r.Header.Get("urlstat-url"); h != "https://changkun.de/research/talks/phd.pdf" {
		t.Errorf("urlstat-url = %q", h)
	}
	if h := r.Header.Get("urlstat-ua"); h != "reader" {
		t.Errorf("urlstat-ua = %q", h)
	}
}
