package server

import (
	"os"
	"path"

	"github.com/ecnepsnai/web"
)

func (v *view) Login(request web.Request) (response web.HTTPResponse) {
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
	f, err := os.Open(path.Join(Directories.Static, "login.html"))
	if err != nil {
		log.Error("Error reading static file: %s", err.Error())
		return web.HTTPResponse{
			Status: 500,
		}
	}
	response.Reader = f
	return
}

func (v *view) JavaScript(request web.Request) web.HTTPResponse {
	file, err := os.OpenFile(path.Join(Directories.Static, "index.html"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("Error serving javascript: %s", err.Error())
		return web.HTTPResponse{
			Status: 500,
		}
	}
	return web.HTTPResponse{
		Reader: file,
	}
}

func (v *view) Favicon(request web.Request) web.HTTPResponse {
	file, err := os.OpenFile(path.Join(Directories.Static, "assets", "img", "favicon.ico"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("Error serving favicon: %s", err.Error())
		return web.HTTPResponse{
			Status: 500,
		}
	}
	return web.HTTPResponse{
		ContentType: "image/x-icon",
		Reader:      file,
	}
}

func (v *view) Redirect(request web.Request) (response web.HTTPResponse) {
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
