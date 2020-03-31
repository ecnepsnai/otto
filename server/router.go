package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ecnepsnai/web"
)

var bindAddress = "localhost:8080"

// RouterSetup set up the HTTP router
func RouterSetup() {
	server := web.New(bindAddress)

	maxBodyLength := uint64(10240)

	authenticatedOptions := web.HandleOptions{
		AuthenticateMethod: func(request *http.Request) interface{} {
			return IsAuthenticated(request)
		},
		MaxBodyLength:      maxBodyLength,
		UnauthorizedMethod: unauthorizedHandle,
	}
	unauthenticatedOptions := web.HandleOptions{
		AuthenticateMethod: func(request *http.Request) interface{} {
			return 1
		},
		MaxBodyLength: maxBodyLength,
	}

	h := handle{}
	v := view{}

	staticPath := path.Join(Directories.StaticBuild, "assets")
	server.HTTP.Static("/static/*filepath", staticPath)
	server.HTTP.Static(fmt.Sprintf("/otto%s/static/*filepath", ServerVersion), staticPath)
	server.HTTP.Static("/clients/*filepath", Directories.Clients)

	// Authentication
	server.HTTP.GET("/login", v.Login, unauthenticatedOptions)
	server.API.POST("/api/login", h.Login, unauthenticatedOptions)
	server.API.POST("/api/logout", h.Logout, authenticatedOptions)

	// Hosts
	server.API.GET("/api/hosts", h.HostList, authenticatedOptions)
	server.API.PUT("/api/hosts/host", h.HostNew, authenticatedOptions)
	server.API.GET("/api/hosts/host/:id", h.HostGet, authenticatedOptions)
	server.API.GET("/api/hosts/host/:id/scripts", h.HostGetScripts, authenticatedOptions)
	server.API.POST("/api/hosts/host/:id", h.HostEdit, authenticatedOptions)
	server.API.DELETE("/api/hosts/host/:id", h.HostDelete, authenticatedOptions)

	// Register
	server.API.PUT("/api/register", h.Register, unauthenticatedOptions)

	// Groups
	server.API.GET("/api/groups", h.GroupList, authenticatedOptions)
	server.API.GET("/api/groups/membership", h.GroupGetMembership, authenticatedOptions)
	server.API.PUT("/api/groups/group", h.GroupNew, authenticatedOptions)
	server.API.GET("/api/groups/group/:id", h.GroupGet, authenticatedOptions)
	server.API.GET("/api/groups/group/:id/scripts", h.GroupGetScripts, authenticatedOptions)
	server.API.GET("/api/groups/group/:id/hosts", h.GroupGetHosts, authenticatedOptions)
	server.API.POST("/api/groups/group/:id", h.GroupEdit, authenticatedOptions)
	server.API.DELETE("/api/groups/group/:id", h.GroupDelete, authenticatedOptions)

	// Heartbeats
	server.API.GET("/api/heartbeat", h.HeartbeatLast, authenticatedOptions)

	// Scripts
	server.API.GET("/api/scripts", h.ScriptList, authenticatedOptions)
	server.API.PUT("/api/scripts/script", h.ScriptNew, authenticatedOptions)
	server.API.GET("/api/scripts/script/:id", h.ScriptGet, authenticatedOptions)
	server.API.GET("/api/scripts/script/:id/hosts", h.ScriptGetHosts, authenticatedOptions)
	server.API.GET("/api/scripts/script/:id/groups", h.ScriptGetGroups, authenticatedOptions)
	server.API.POST("/api/scripts/script/:id/groups", h.ScriptSetGroups, authenticatedOptions)
	server.API.POST("/api/scripts/script/:id", h.ScriptEdit, authenticatedOptions)
	server.API.DELETE("/api/scripts/script/:id", h.ScriptDelete, authenticatedOptions)

	// Request
	server.API.PUT("/api/request", h.RequestNew, authenticatedOptions)

	// State
	server.API.GET("/api/state", h.State, authenticatedOptions)

	// Users
	server.API.GET("/api/users", h.UserList, authenticatedOptions)
	server.API.PUT("/api/users/user", h.UserNew, authenticatedOptions)
	server.API.GET("/api/users/user/:username", h.UserGet, authenticatedOptions)
	server.API.POST("/api/users/user/:username", h.UserEdit, authenticatedOptions)
	server.API.DELETE("/api/users/user/:username", h.UserDelete, authenticatedOptions)

	// Options
	server.API.GET("/api/options", h.OptionsGet, authenticatedOptions)
	server.API.POST("/api/options", h.OptionsSet, authenticatedOptions)

	// Redirect
	server.HTTP.GET("/", v.Redirect, unauthenticatedOptions)

	server.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "application/json") {
			json.NewEncoder(w).Encode(web.CommonErrors.NotFound)
		} else {
			file, err := os.OpenFile(path.Join(Directories.Static, "build", "not_found.html"), os.O_RDONLY, os.ModePerm)
			defer file.Close()
			if err != nil {
				panic(err)
			}
			io.CopyBuffer(w, file, nil)
		}
	}

	ngRoutes := []string{
		"/hosts/",
		"/hosts/host/",
		"/hosts/host/:id/",
		"/hosts/host/:id/edit/",
		"/groups/",
		"/groups/group/",
		"/groups/group/:id/",
		"/groups/group/:id/edit/",
		"/scripts/",
		"/scripts/script/",
		"/scripts/script/:id/",
		"/scripts/script/:id/edit/",
		"/scripts/script/:id/execute/",
		"/options/",
		"/options/users/user/",
		"/options/users/user/:username",
	}
	for _, route := range ngRoutes {
		server.HTTP.GET(route, v.AngularJS, authenticatedOptions)
	}

	server.Start()
}

func unauthorizedHandle(w http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.Header.Get("Accept"), "text/html") {
		w.Header().Add("Location", "/login?unauthorized")
		cookie := http.Cookie{
			Name:    ottoSessionCookie,
			Value:   "",
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, -1),
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(307)
		return
	}

	w.WriteHeader(403)
	w.Write([]byte("{\"error\":{\"code\":403,\"message\":\"unauthorized\"}}"))
	return
}
