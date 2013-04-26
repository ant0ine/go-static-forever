package forever

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func codeIs(t *testing.T, r *httptest.ResponseRecorder, expectedCode int) {
	if r.Code != expectedCode {
		t.Errorf("Code %d expected, got: %d", expectedCode, r.Code)
	}
}

func TestWrongVersion(t *testing.T) {
	// wring version string, expect a 404

	handler := NewStaticHandler(
		http.Dir("."),
		"1234567",
		nil,
		true,
	)

	urlObj, err := url.Parse("http://1.2.3.4/wrong_version/handler.go")
	if err != nil {
		t.Fatal(err)
	}
	r := http.Request{
		Method: "GET",
		URL:    urlObj,
	}
	r.Header = http.Header{}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, &r)

	codeIs(t, recorder, 404)
}

func TestDevelopement(t *testing.T) {
	// right version string, expect a 200

	handler := NewStaticHandler(
		http.Dir("."),
		"1234567",
		nil,
		true,
	)

	urlObj, err := url.Parse("http://1.2.3.4/1234567/handler.go")
	if err != nil {
		t.Fatal(err)
	}
	r := http.Request{
		Method: "GET",
		URL:    urlObj,
	}
	r.Header = http.Header{}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, &r)

	codeIs(t, recorder, 200)
}

func TestFirstRequest(t *testing.T) {
	// expect a 200 with all the custom headers set

	handler := NewStaticHandler(
		http.Dir("."),
		"1234567",
		nil,
		false,
	)

	handler = handler.(*staticHandler)

	urlObj, err := url.Parse("http://1.2.3.4/1234567/handler.go")
	if err != nil {
		t.Fatal(err)
	}
	r := http.Request{
		Method: "GET",
		URL:    urlObj,
	}
	r.Header = http.Header{}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, &r)

	codeIs(t, recorder, 200)

	expires := recorder.HeaderMap.Get("Expires")
	if expires == "" {
		t.Errorf("Expires expected, got: %s", expires)
	}

	control := recorder.HeaderMap.Get("Cache-Control")
	if control != "public; max-age=31536000; s-maxage=31536000" {
		t.Errorf("Cache-Control public; max-age=31536000; s-maxage=31536000 expected, got: %s", control)
	}

	last := recorder.HeaderMap.Get("Last-Modified")
	if last != "Thu, 01 Jan 1970 00:00:00 GMT" {
		t.Errorf("Last-Modified Thu, 01 Jan 1970 00:00:00 GMT expected, got: %s", last)
	}
}

func TestSecondRequest(t *testing.T) {
	// sent with If-Modified-Since, expect a 304

	handler := NewStaticHandler(
		http.Dir("."),
		"1234567",
		nil,
		false,
	)

	handler = handler.(*staticHandler)

	urlObj, err := url.Parse("http://1.2.3.4/1234567/handler.go")
	if err != nil {
		t.Fatal(err)
	}
	r := http.Request{
		Method: "GET",
		URL:    urlObj,
	}
	r.Header = http.Header{}
	r.Header.Set("If-Modified-Since", "Sat, 01 Apr 2113 04:15:01 GMT")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, &r)

	codeIs(t, recorder, 304)
}
