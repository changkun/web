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
