package httpext

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestAcceptedLanguages(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set(AcceptedLanguage, "da, en-GB;q=0.8, en;q=0.7")

	languages := AcceptedLanguages(req)

	Equal(t, languages[0], "da")
	Equal(t, languages[1], "en-GB")
	Equal(t, languages[2], "en")

	req.Header.Del(AcceptedLanguage)

	languages = AcceptedLanguages(req)
	Equal(t, len(languages), 0)

	req.Header.Set(AcceptedLanguage, "")
	languages = AcceptedLanguages(req)
	Equal(t, len(languages), 0)
}

func TestAttachment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("../../README.md")
		if err := Attachment(w, f, "README.md"); err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/dl-unknown-type", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("../../README.md")
		if err := Attachment(w, f, "readme"); err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/dl-fake-png", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("../../README.md")
		if err := Attachment(w, f, "logo.png"); err != nil {
			panic(err)
		}
	})

	tests := []struct {
		name        string
		code        int
		disposition string
		typ         string
		url         string
	}{
		{
			code:        http.StatusOK,
			disposition: "attachment;filename=README.md",
			typ:         TextMarkdown,
			url:         "/dl",
		},
		{
			code:        http.StatusOK,
			disposition: "attachment;filename=readme",
			typ:         ApplicationOctetStream,
			url:         "/dl-unknown-type",
		},
		{
			code:        http.StatusOK,
			disposition: "attachment;filename=logo.png",
			typ:         ImagePNG,
			url:         "/dl-fake-png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			Equal(t, err, nil)

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if tt.code != w.Code {
				t.Errorf("Status Code = %d, want %d", w.Code, tt.code)
			}
			if tt.disposition != w.Header().Get(ContentDisposition) {
				t.Errorf("Content Disaposition = %s, want %s", w.Header().Get(ContentDisposition), tt.disposition)
			}
			if tt.typ != w.Header().Get(ContentType) {
				t.Errorf("Content Type = %s, want %s", w.Header().Get(ContentType), tt.typ)
			}
		})
	}
}

func TestInline(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/dl-inline", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("../../README.md")
		if err := Inline(w, f, "README.md"); err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/dl-unknown-type-inline", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("../../README.md")
		if err := Inline(w, f, "readme"); err != nil {
			panic(err)
		}
	})

	tests := []struct {
		name        string
		code        int
		disposition string
		typ         string
		url         string
	}{
		{
			code:        http.StatusOK,
			disposition: "inline;filename=README.md",
			typ:         TextMarkdown,
			url:         "/dl-inline",
		},
		{
			code:        http.StatusOK,
			disposition: "inline;filename=readme",
			typ:         ApplicationOctetStream,
			url:         "/dl-unknown-type-inline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			Equal(t, err, nil)

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if tt.code != w.Code {
				t.Errorf("Status Code = %d, want %d", w.Code, tt.code)
			}
			if tt.disposition != w.Header().Get(ContentDisposition) {
				t.Errorf("Content Disaposition = %s, want %s", w.Header().Get(ContentDisposition), tt.disposition)
			}
			if tt.typ != w.Header().Get(ContentType) {
				t.Errorf("Content Type = %s, want %s", w.Header().Get(ContentType), tt.typ)
			}
		})
	}
}

func TestClientIP(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("X-Real-IP", " 10.10.10.10  ")
	req.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	req.RemoteAddr = "  40.40.40.40:42123 "

	Equal(t, ClientIP(req), "10.10.10.10")

	req.Header.Del("X-Real-IP")
	Equal(t, ClientIP(req), "20.20.20.20")

	req.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	Equal(t, ClientIP(req), "30.30.30.30")

	req.Header.Del("X-Forwarded-For")
	Equal(t, ClientIP(req), "40.40.40.40")
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	type test struct {
		Field string `json:"field"`
	}
	tst := test{Field: "myfield"}
	b, err := json.Marshal(tst)
	Equal(t, err, nil)

	err = JSON(w, http.StatusOK, tst)
	Equal(t, err, nil)
	Equal(t, w.Header().Get(ContentType), ApplicationJSON)
	Equal(t, w.Body.Bytes(), b)

	err = JSON(w, http.StatusOK, func() {})
	NotEqual(t, err, nil)
}

func TestJSONBytes(t *testing.T) {
	w := httptest.NewRecorder()
	type test struct {
		Field string `json:"field"`
	}
	tst := test{Field: "myfield"}
	b, err := json.Marshal(tst)
	Equal(t, err, nil)

	err = JSONBytes(w, http.StatusOK, b)
	Equal(t, err, nil)
	Equal(t, w.Header().Get(ContentType), ApplicationJSON)
	Equal(t, w.Body.Bytes(), b)
}

func TestJSONP(t *testing.T) {
	callbackFunc := "CallbackFunc"
	w := httptest.NewRecorder()
	type test struct {
		Field string `json:"field"`
	}
	tst := test{Field: "myfield"}
	err := JSONP(w, http.StatusOK, tst, callbackFunc)
	Equal(t, err, nil)
	Equal(t, w.Header().Get(ContentType), ApplicationJSON)

	err = JSON(w, http.StatusOK, func() {})
	NotEqual(t, err, nil)
}

func TestXML(t *testing.T) {
	w := httptest.NewRecorder()
	type zombie struct {
		ID   int    `json:"id"   xml:"id"`
		Name string `json:"name" xml:"name"`
	}
	tst := zombie{1, "Patient Zero"}
	xmlData := `<zombie><id>1</id><name>Patient Zero</name></zombie>`
	err := XML(w, http.StatusOK, tst)
	Equal(t, err, nil)
	Equal(t, w.Header().Get(ContentType), ApplicationXML)
	Equal(t, w.Body.Bytes(), []byte(xml.Header+xmlData))

	err = JSON(w, http.StatusOK, func() {})
	NotEqual(t, err, nil)
}

func TestXMLBytes(t *testing.T) {
	xmlData := `<zombie><id>1</id><name>Patient Zero</name></zombie>`
	w := httptest.NewRecorder()
	err := XMLBytes(w, http.StatusOK, []byte(xmlData))
	Equal(t, err, nil)
	Equal(t, w.Header().Get(ContentType), ApplicationXML)
	Equal(t, w.Body.Bytes(), []byte(xml.Header+xmlData))
}

func TestDecode(t *testing.T) {
	type TestStruct struct {
		ID              int `form:"id"`
		Posted          string
		MultiPartPosted string
	}

	test := new(TestStruct)

	mux := http.NewServeMux()
	mux.HandleFunc("/decode-noquery", func(w http.ResponseWriter, r *http.Request) {
		err := Decode(r, NoQueryParams, 16<<10, test)
		Equal(t, err, nil)
	})
	mux.HandleFunc("/decode-query", func(w http.ResponseWriter, r *http.Request) {
		err := Decode(r, QueryParams, 16<<10, test)
		Equal(t, err, nil)
	})

	// test query params
	r, _ := http.NewRequest(http.MethodGet, "/decode-query?id=5", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 5)
	Equal(t, test.Posted, "")
	Equal(t, test.MultiPartPosted, "")

	// test Form decode
	form := url.Values{}
	form.Add("Posted", "values")

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-query?id=13", strings.NewReader(form.Encode()))
	r.Header.Set(ContentType, ApplicationForm)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 13)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "")

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-noquery?id=14", strings.NewReader(form.Encode()))
	r.Header.Set(ContentType, ApplicationForm)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 0)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "")

	// test MultipartForm
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("MultiPartPosted", "values")
	Equal(t, err, nil)

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	err = writer.Close()
	Equal(t, err, nil)

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-query?id=12", body)
	r.Header.Set(ContentType, writer.FormDataContentType())
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 12)
	Equal(t, test.Posted, "")
	Equal(t, test.MultiPartPosted, "values")

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	err = writer.WriteField("MultiPartPosted", "values")
	Equal(t, err, nil)

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	err = writer.Close()
	Equal(t, err, nil)

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-noquery?id=13", body)
	r.Header.Set(ContentType, writer.FormDataContentType())
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 0)
	Equal(t, test.Posted, "")
	Equal(t, test.MultiPartPosted, "values")

	// test JSON
	jsonBody := `{"ID":13,"Posted":"values","MultiPartPosted":"values"}`
	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-query?id=13", strings.NewReader(jsonBody))
	r.Header.Set(ContentType, ApplicationJSON)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 13)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")

	var buff bytes.Buffer
	gzw := gzip.NewWriter(&buff)
	defer func() {
		_ = gzw.Close()
	}()
	_, err = gzw.Write([]byte(jsonBody))
	Equal(t, err, nil)

	err = gzw.Close()
	Equal(t, err, nil)

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-query?id=14", &buff)
	r.Header.Set(ContentType, ApplicationJSON)
	r.Header.Set(ContentEncoding, Gzip)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 14)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-noquery?id=14", strings.NewReader(jsonBody))
	r.Header.Set(ContentType, ApplicationJSON)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 13)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")

	// test XML
	xmlBody := `<TestStruct><ID>13</ID><Posted>values</Posted><MultiPartPosted>values</MultiPartPosted></TestStruct>`
	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-noquery?id=14", strings.NewReader(xmlBody))
	r.Header.Set(ContentType, ApplicationXML)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 13)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-query?id=14", strings.NewReader(xmlBody))
	r.Header.Set(ContentType, ApplicationXML)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 14)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")

	buff.Reset()
	gzw = gzip.NewWriter(&buff)
	defer func() {
		_ = gzw.Close()
	}()
	_, err = gzw.Write([]byte(xmlBody))
	Equal(t, err, nil)

	err = gzw.Close()
	Equal(t, err, nil)

	test = new(TestStruct)
	r, _ = http.NewRequest(http.MethodPost, "/decode-noquery?id=14", &buff)
	r.Header.Set(ContentType, ApplicationXML)
	r.Header.Set(ContentEncoding, Gzip)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	Equal(t, w.Code, http.StatusOK)
	Equal(t, test.ID, 13)
	Equal(t, test.Posted, "values")
	Equal(t, test.MultiPartPosted, "values")
}

func TestDecodeResponse(t *testing.T) {

	type result struct {
		ID int `json:"id" xml:"id"`
	}

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		expected result
	}{
		{
			name: "Test JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				Equal(t, JSON(w, http.StatusOK, result{ID: 3}), nil)
			},
			expected: result{ID: 3},
		},
		{
			name: "Test XML",
			handler: func(w http.ResponseWriter, r *http.Request) {
				Equal(t, XML(w, http.StatusOK, result{ID: 5}), nil)
			},
			expected: result{ID: 5},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", tc.handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			Equal(t, err, nil)

			resp, err := http.DefaultClient.Do(req)
			Equal(t, err, nil)
			Equal(t, resp.StatusCode, http.StatusOK)

			res, err := DecodeResponse[result](resp, bytesext.MiB)
			Equal(t, err, nil)
			Equal(t, tc.expected.ID, res.ID)
		})
	}
}
