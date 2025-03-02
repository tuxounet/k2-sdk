package files

import (
	"strings"
)

func GetExentionFromFilename(filename string) string {
	extension := ""

	extension = strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	return extension

}

func GetContentTypeFromExtension(ext string) string {

	switch ext {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "svg":
		return "image/svg+xml"
	case "png":
		return "image/png"
	case "jpg":
		return "image/jpeg"
	case "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "ico":
		return "image/x-icon"
	case "webp":
		return "image/webp"
	case "woff":
		return "font/woff"
	case "woff2":
		return "font/woff2"
	case "ttf":
		return "font/ttf"
	case "otf":
		return "font/otf"
	case "eot":
		return "application/vnd.ms-fontobject"
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	case "ogg":
		return "video/ogg"
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "weba":
		return "audio/webm"
	case "flac":
		return "audio/flac"
	case "txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}

}
