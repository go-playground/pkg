package httpext

const (
	charsetUTF8 = "; charset=" + UTF8
)

// Mime Type values for the Content-Type HTTP header
const (
	ApplicationJSON       string = "application/json" + charsetUTF8
	ApplicationJavaScript string = "application/javascript"
	ApplicationXML        string = "application/xml" + charsetUTF8
	ApplicationForm       string = "application/x-www-form-urlencoded"
	ApplicationProtobuf   string = "application/protobuf"
	ApplicationMsgpack    string = "application/msgpack"
	ApplicationWasm       string = "application/wasm"
	ApplicationPDF        string = "application/pdf"
	TextHTML              string = "text/html" + charsetUTF8
	TextPlain             string = "text/plain" + charsetUTF8
	TextMarkdown          string = "text/markdown" + charsetUTF8
	TextCSS               string = "text/css" + charsetUTF8
	ImagePNG              string = "image/png"
	ImageGIF              string = "image/gif"
	ImageSVG              string = "image/svg+xml"
	ImageJPEG             string = "image/jpeg"
	MultipartForm         string = "multipart/form-data"
	OctetStream           string = "application/octet-stream"
)
