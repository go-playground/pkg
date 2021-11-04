package httpext

const (
	charsetUTF8 = "; charset=" + UTF8
)

// Mime Type values for the Content-Type HTTP header
const (
	ApplicationJSONNoCharset string = "application/json"
	ApplicationJSON          string = ApplicationJSONNoCharset + charsetUTF8
	ApplicationJavaScript    string = "application/javascript"
	ApplicationXMLNoCharset  string = "application/xml"
	ApplicationXML           string = ApplicationXMLNoCharset + charsetUTF8
	ApplicationForm          string = "application/x-www-form-urlencoded"
	ApplicationProtobuf      string = "application/protobuf"
	ApplicationMsgpack       string = "application/msgpack"
	ApplicationWasm          string = "application/wasm"
	ApplicationPDF           string = "application/pdf"
	ApplicationOctetStream   string = "application/octet-stream"
	TextHTMLNoCharset               = "text/html"
	TextHTML                 string = TextHTMLNoCharset + charsetUTF8
	TextPlainNoCharset              = "text/plain"
	TextPlain                string = TextPlainNoCharset + charsetUTF8
	TextMarkdownNoCharset    string = "text/markdown"
	TextMarkdown             string = TextMarkdownNoCharset + charsetUTF8
	TextCSSNoCharset         string = "text/css"
	TextCSS                  string = TextCSSNoCharset + charsetUTF8
	TextCSV                  string = "text/csv"
	ImagePNG                 string = "image/png"
	ImageGIF                 string = "image/gif"
	ImageSVG                 string = "image/svg+xml"
	ImageJPEG                string = "image/jpeg"
	MultipartForm            string = "multipart/form-data"
)
