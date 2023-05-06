package httpext

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	ioext "github.com/go-playground/pkg/v5/io"
)

// QueryParamsOption represents the options for including query parameters during Decode helper functions
type QueryParamsOption uint8

// QueryParamsOption's
const (
	QueryParams QueryParamsOption = iota
	NoQueryParams
)

var (
	xmlHeaderBytes = []byte(xml.Header)
)

func detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if t := mime.TypeByExtension(ext); t != "" {
		return t
	}
	switch ext {
	case ".md":
		return TextMarkdown
	default:
		return ApplicationOctetStream
	}
}

// AcceptedLanguages returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func AcceptedLanguages(r *http.Request) (languages []string) {
	accepted := r.Header.Get(AcceptedLanguage)
	if accepted == "" {
		return
	}
	options := strings.Split(accepted, ",")
	l := len(options)
	languages = make([]string, l)

	for i := 0; i < l; i++ {
		locale := strings.SplitN(options[i], ";", 2)
		languages[i] = strings.Trim(locale[0], " ")
	}
	return
}

// Attachment is a helper method for returning an attachement file
// to be downloaded, if you with to open inline see function Inline
func Attachment(w http.ResponseWriter, r io.Reader, filename string) (err error) {
	w.Header().Set(ContentDisposition, "attachment;filename="+filename)
	w.Header().Set(ContentType, detectContentType(filename))
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, r)
	return
}

// Inline is a helper method for returning a file inline to
// be rendered/opened by the browser
func Inline(w http.ResponseWriter, r io.Reader, filename string) (err error) {
	w.Header().Set(ContentDisposition, "inline;filename="+filename)
	w.Header().Set(ContentType, detectContentType(filename))
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, r)
	return
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func ClientIP(r *http.Request) (clientIP string) {
	values := r.Header[XRealIP]
	if len(values) > 0 {
		clientIP = strings.TrimSpace(values[0])
		if clientIP != "" {
			return
		}
	}
	if values = r.Header[XForwardedFor]; len(values) > 0 {
		clientIP = values[0]
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if clientIP != "" {
			return
		}
	}
	clientIP, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	return
}

// JSONStream uses json.Encoder to stream the JSON reponse body.
//
// This differs from the JSON helper which unmarshalls into memory first allowing the capture of JSON encoding errors.
func JSONStream(w http.ResponseWriter, status int, i interface{}) error {
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(i)
}

// JSON marshals provided interface + returns JSON + status code
func JSON(w http.ResponseWriter, status int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

// JSONBytes returns provided JSON response with status code
func JSONBytes(w http.ResponseWriter, status int, b []byte) (err error) {
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

// JSONP sends a JSONP response with status code and uses `callback` to construct
// the JSONP payload.
func JSONP(w http.ResponseWriter, status int, i interface{}, callback string) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(status)
	if _, err = w.Write([]byte(callback + "(")); err == nil {
		if _, err = w.Write(b); err == nil {
			_, err = w.Write([]byte(");"))
		}
	}
	return err
}

// XML marshals provided interface + returns XML + status code
func XML(w http.ResponseWriter, status int, i interface{}) error {
	b, err := xml.Marshal(i)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ApplicationXML)
	w.WriteHeader(status)
	if _, err = w.Write(xmlHeaderBytes); err == nil {
		_, err = w.Write(b)
	}
	return err
}

// XMLBytes returns provided XML response with status code
func XMLBytes(w http.ResponseWriter, status int, b []byte) (err error) {
	w.Header().Set(ContentType, ApplicationXML)
	w.WriteHeader(status)
	if _, err = w.Write(xmlHeaderBytes); err == nil {
		_, err = w.Write(b)
	}
	return
}

// DecodeForm parses the requests form data into the provided struct.
//
// The Content-Type and http method are not checked.
//
// NOTE: when QueryParamsOption=QueryParams the query params will be parsed and included eg. route /user?test=true 'test'
// is added to parsed Form.
func DecodeForm(r *http.Request, qp QueryParamsOption, v interface{}) (err error) {
	if err = r.ParseForm(); err == nil {
		switch qp {
		case QueryParams:
			err = DefaultFormDecoder.Decode(v, r.Form)
		case NoQueryParams:
			err = DefaultFormDecoder.Decode(v, r.PostForm)
		}
	}
	return
}

// DecodeMultipartForm parses the requests form data into the provided struct.
//
// The Content-Type and http method are not checked.
//
// NOTE: when includeQueryParams=true query params will be parsed and included eg. route /user?test=true 'test'
// is added to parsed MultipartForm.
func DecodeMultipartForm(r *http.Request, qp QueryParamsOption, maxMemory int64, v interface{}) (err error) {
	if err = r.ParseMultipartForm(maxMemory); err == nil {
		switch qp {
		case QueryParams:
			err = DefaultFormDecoder.Decode(v, r.Form)
		case NoQueryParams:
			err = DefaultFormDecoder.Decode(v, r.MultipartForm.Value)
		}
	}
	return
}

// DecodeJSON decodes the request body into the provided struct and limits the request size via
// an ioext.LimitReader using the maxMemory param.
//
// The Content-Type e.g. "application/json" and http method are not checked.
//
// NOTE: when includeQueryParams=true query params will be parsed and included eg. route /user?test=true 'test'
// is added to parsed JSON and replaces any values that may have been present
func DecodeJSON(r *http.Request, qp QueryParamsOption, maxMemory int64, v interface{}) (err error) {
	return decodeJSON(r.Header, r.Body, qp, r.URL.Query(), maxMemory, v)
}

func decodeJSON(headers http.Header, body io.Reader, qp QueryParamsOption, values url.Values, maxMemory int64, v interface{}) (err error) {
	if encoding := headers.Get(ContentEncoding); encoding == Gzip {
		var gzr *gzip.Reader
		gzr, err = gzip.NewReader(body)
		if err != nil {
			return
		}
		defer func() {
			_ = gzr.Close()
		}()
		body = gzr
	}
	err = json.NewDecoder(ioext.LimitReader(body, maxMemory)).Decode(v)
	if qp == QueryParams && err == nil {
		err = decodeQueryParams(values, v)
	}
	return
}

// DecodeXML decodes the request body into the provided struct and limits the request size via
// an ioext.LimitReader using the maxMemory param.
//
// The Content-Type e.g. "application/xml" and http method are not checked.
//
// NOTE: when includeQueryParams=true query params will be parsed and included eg. route /user?test=true 'test'
// is added to parsed XML and replaces any values that may have been present
func DecodeXML(r *http.Request, qp QueryParamsOption, maxMemory int64, v interface{}) (err error) {
	return decodeXML(r.Header, r.Body, qp, r.URL.Query(), maxMemory, v)
}

func decodeXML(headers http.Header, body io.Reader, qp QueryParamsOption, values url.Values, maxMemory int64, v interface{}) (err error) {
	if encoding := headers.Get(ContentEncoding); encoding == Gzip {
		var gzr *gzip.Reader
		gzr, err = gzip.NewReader(body)
		if err != nil {
			return
		}
		defer func() {
			_ = gzr.Close()
		}()
		body = gzr
	}
	err = xml.NewDecoder(ioext.LimitReader(body, maxMemory)).Decode(v)
	if qp == QueryParams && err == nil {
		err = decodeQueryParams(values, v)
	}
	return
}

// DecodeQueryParams takes the URL Query params flag.
func DecodeQueryParams(r *http.Request, v interface{}) (err error) {
	return decodeQueryParams(r.URL.Query(), v)
}

func decodeQueryParams(values url.Values, v interface{}) (err error) {
	err = DefaultFormDecoder.Decode(v, values)
	return
}

const (
	nakedApplicationJSON string = "application/json"
	nakedApplicationXML  string = "application/xml"
)

// Decode takes the request and attempts to discover its content type via
// the http headers and then decode the request body into the provided struct.
// Example if header was "application/json" would decode using
// json.NewDecoder(ioext.LimitReader(r.Body, maxMemory)).Decode(v).
//
// This default to parsing query params if includeQueryParams=true and no other content type matches.
//
// NOTE: when includeQueryParams=true query params will be parsed and included eg. route /user?test=true 'test'
// is added to parsed XML and replaces any values that may have been present
func Decode(r *http.Request, qp QueryParamsOption, maxMemory int64, v interface{}) (err error) {
	typ := r.Header.Get(ContentType)
	if idx := strings.Index(typ, ";"); idx != -1 {
		typ = typ[:idx]
	}
	switch typ {
	case nakedApplicationJSON:
		err = DecodeJSON(r, qp, maxMemory, v)
	case nakedApplicationXML:
		err = DecodeXML(r, qp, maxMemory, v)
	case ApplicationForm:
		err = DecodeForm(r, qp, v)
	case MultipartForm:
		err = DecodeMultipartForm(r, qp, maxMemory, v)
	default:
		if qp == QueryParams {
			err = DecodeQueryParams(r, v)
		}
	}
	return
}

// DecodeResponse takes the response and attempts to discover its content type via
// the http headers and then decode the request body into the provided type.
//
// Example if header was "application/json" would decode using
// json.NewDecoder(ioext.LimitReader(r.Body, maxMemory)).Decode(v).
func DecodeResponse[T any](r *http.Response, maxMemory int64) (result T, err error) {
	typ := r.Header.Get(ContentType)
	if idx := strings.Index(typ, ";"); idx != -1 {
		typ = typ[:idx]
	}
	switch typ {
	case nakedApplicationJSON:
		err = decodeJSON(r.Header, r.Body, NoQueryParams, nil, maxMemory, &result)
	case nakedApplicationXML:
		err = decodeXML(r.Header, r.Body, NoQueryParams, nil, maxMemory, &result)
	default:
		err = errors.New("unsupported content type")
	}
	return
}
