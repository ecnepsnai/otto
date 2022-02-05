package server

import (
	"os"
	"path"

	"github.com/ecnepsnai/web"
)

func (v *view) Login(request web.Request, writer web.Writer) (response web.Response) {
	// Redirect users to index if they're already logged in
	session := sessionForHTTPRequest(request.HTTP, false)
	if session != nil {
		response.Headers = map[string]string{
			"Location": "/hosts",
		}
		response.Status = 307
		return
	}

	response.ContentType = "text/html; charset=utf-8"
	f, err := os.Open(path.Join(Directories.Static, "build", "login.html"))
	if err != nil {
		log.Error("Error reading static file: %s", err.Error())
		return web.Response{
			Status: 500,
		}
	}
	response.Reader = f
	return
}

func (v *view) JavaScript(request web.Request, writer web.Writer) web.Response {
	file, err := os.OpenFile(path.Join(Directories.Build, "index.html"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("Error serving javascript: %s", err.Error())
		return web.Response{
			Status: 500,
		}
	}
	return web.Response{
		Reader: file,
	}
}

func (v *view) Favicon(request web.Request, writer web.Writer) web.Response {
	file, err := os.OpenFile(path.Join(Directories.Build, "assets", "img", "favicon.ico"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("Error serving favicon: %s", err.Error())
		return web.Response{
			Status: 500,
		}
	}
	return web.Response{
		ContentType: "image/x-icon",
		Reader:      file,
	}
}

func (v *view) Redirect(request web.Request, writer web.Writer) (response web.Response) {
	redirectLocation := ""

	session := sessionForHTTPRequest(request.HTTP, false)
	if session != nil && !session.Partial {
		redirectLocation = "/hosts"
	} else {
		redirectLocation = "/login"
	}

	response.Headers = map[string]string{
		"Location": redirectLocation,
	}
	response.Status = 307
	return
}
