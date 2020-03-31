package server

import (
	"os"
	"path"

	"github.com/ecnepsnai/web"
)

func (v *view) Login(request web.Request, writer web.Writer) (response web.Response) {
	// Redirect users to index if they're already logged in
	session := IsAuthenticated(request.HTTP)
	if session != nil {
		response.Headers = map[string]string{
			"Location": "/hosts/",
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

func (v *view) AngularJS(request web.Request, writer web.Writer) (response web.Response) {
	response.ContentType = "text/html; charset=utf-8"
	f, err := os.Open(path.Join(Directories.Static, "build", "ng.html"))
	if err != nil {
		log.Error("Error reading static file: %s", err.Error())
		return web.Response{
			Status: 500,
		}
	}
	response.Reader = f
	return
}

func (v *view) Redirect(request web.Request, writer web.Writer) (response web.Response) {
	redirectLocation := ""

	session := IsAuthenticated(request.HTTP)
	if session != nil {
		redirectLocation = "/hosts/"
	} else {
		redirectLocation = "/login"
	}

	response.Headers = map[string]string{
		"Location": redirectLocation,
	}
	response.Status = 307
	return
}
