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
	server.MaxRequestsPerSecond = 15

	maxBodyLength := uint64(104857600)

	authenticatedOptions := func(allowPartial bool) web.HandleOptions {
		return web.HandleOptions{
			AuthenticateMethod: func(request *http.Request) interface{} {
				return sessionForHTTPRequest(request, allowPartial)
			},
			MaxBodyLength:      maxBodyLength,
			UnauthorizedMethod: unauthorizedHandle,
		}
	}
	unauthenticatedOptions := web.HandleOptions{
		AuthenticateMethod: func(request *http.Request) interface{} {
			return 1
		},
		MaxBodyLength: maxBodyLength,
	}

	h := handle{}
	v := view{}

	server.HTTP.Static("/static/*filepath", Directories.Build)
	server.HTTP.Static(fmt.Sprintf("/otto%s/*filepath", ServerVersion), Directories.Build)
	server.HTTP.Static("/clients/*filepath", Directories.Clients)

	// Authentication
	server.HTTP.GET("/login", v.Login, unauthenticatedOptions)
	server.API.POST("/api/login", h.Login, unauthenticatedOptions)
	server.API.POST("/api/logout", h.Logout, authenticatedOptions(true))

	// Hosts
	server.API.GET("/api/hosts", h.HostList, authenticatedOptions(false))
	server.API.PUT("/api/hosts/host", h.HostNew, authenticatedOptions(false))
	server.API.GET("/api/hosts/host/:id", h.HostGet, authenticatedOptions(false))
	server.API.GET("/api/hosts/host/:id/scripts", h.HostGetScripts, authenticatedOptions(false))
	server.API.GET("/api/hosts/host/:id/groups", h.HostGetGroups, authenticatedOptions(false))
	server.API.GET("/api/hosts/host/:id/schedules", h.HostGetSchedules, authenticatedOptions(false))
	server.API.POST("/api/hosts/host/:id/psk", h.HostRotatePSK, authenticatedOptions(false))
	server.API.POST("/api/hosts/host/:id", h.HostEdit, authenticatedOptions(false))
	server.API.DELETE("/api/hosts/host/:id", h.HostDelete, authenticatedOptions(false))

	// Register
	server.API.PUT("/api/register", h.Register, unauthenticatedOptions)
	// Register Rules
	server.API.GET("/api/register/rules", h.RegisterRuleList, authenticatedOptions(false))
	server.API.PUT("/api/register/rules/rule", h.RegisterRuleNew, authenticatedOptions(false))
	server.API.GET("/api/register/rules/rule/:id", h.RegisterRuleGet, authenticatedOptions(false))
	server.API.POST("/api/register/rules/rule/:id", h.RegisterRuleEdit, authenticatedOptions(false))
	server.API.DELETE("/api/register/rules/rule/:id", h.RegisterRuleDelete, authenticatedOptions(false))

	// Groups
	server.API.GET("/api/groups", h.GroupList, authenticatedOptions(false))
	server.API.GET("/api/groups/membership", h.GroupGetMembership, authenticatedOptions(false))
	server.API.PUT("/api/groups/group", h.GroupNew, authenticatedOptions(false))
	server.API.GET("/api/groups/group/:id", h.GroupGet, authenticatedOptions(false))
	server.API.GET("/api/groups/group/:id/scripts", h.GroupGetScripts, authenticatedOptions(false))
	server.API.GET("/api/groups/group/:id/hosts", h.GroupGetHosts, authenticatedOptions(false))
	server.API.GET("/api/groups/group/:id/schedules", h.GroupGetSchedules, authenticatedOptions(false))
	server.API.POST("/api/groups/group/:id/hosts", h.GroupSetHosts, authenticatedOptions(false))
	server.API.POST("/api/groups/group/:id", h.GroupEdit, authenticatedOptions(false))
	server.API.DELETE("/api/groups/group/:id", h.GroupDelete, authenticatedOptions(false))

	// Schedules
	server.API.GET("/api/schedules", h.ScheduleList, authenticatedOptions(false))
	server.API.PUT("/api/schedules/schedule", h.ScheduleNew, authenticatedOptions(false))
	server.API.GET("/api/schedules/schedule/:id", h.ScheduleGet, authenticatedOptions(false))
	server.API.GET("/api/schedules/schedule/:id/reports", h.ScheduleGetReports, authenticatedOptions(false))
	server.API.GET("/api/schedules/schedule/:id/hosts", h.ScheduleGetHosts, authenticatedOptions(false))
	server.API.GET("/api/schedules/schedule/:id/groups", h.ScheduleGetGroups, authenticatedOptions(false))
	server.API.GET("/api/schedules/schedule/:id/script", h.ScheduleGetScript, authenticatedOptions(false))
	server.API.POST("/api/schedules/schedule/:id", h.ScheduleEdit, authenticatedOptions(false))
	server.API.DELETE("/api/schedules/schedule/:id", h.ScheduleDelete, authenticatedOptions(false))

	// Heartbeats
	server.API.GET("/api/heartbeat", h.HeartbeatLast, authenticatedOptions(false))

	// Scripts
	server.API.GET("/api/scripts", h.ScriptList, authenticatedOptions(false))
	server.API.PUT("/api/scripts/script", h.ScriptNew, authenticatedOptions(false))
	server.API.GET("/api/scripts/script/:id", h.ScriptGet, authenticatedOptions(false))
	server.API.GET("/api/scripts/script/:id/hosts", h.ScriptGetHosts, authenticatedOptions(false))
	server.API.GET("/api/scripts/script/:id/groups", h.ScriptGetGroups, authenticatedOptions(false))
	server.API.GET("/api/scripts/script/:id/schedules", h.ScriptGetSchedules, authenticatedOptions(false))
	server.API.GET("/api/scripts/script/:id/attachments", h.ScriptGetAttachments, authenticatedOptions(false))
	server.API.POST("/api/scripts/script/:id/groups", h.ScriptSetGroups, authenticatedOptions(false))
	server.API.POST("/api/scripts/script/:id", h.ScriptEdit, authenticatedOptions(false))
	server.API.DELETE("/api/scripts/script/:id", h.ScriptDelete, authenticatedOptions(false))

	// Attachments
	server.API.GET("/api/attachments", h.AttachmentList, authenticatedOptions(false))
	server.API.PUT("/api/attachments", h.AttachmentUpload, authenticatedOptions(false))
	server.API.GET("/api/attachments/attachment/:id", h.AttachmentGet, authenticatedOptions(false))
	server.HTTP.GET("/api/attachments/attachment/:id/download", v.AttachmentDownload, authenticatedOptions(false))
	server.API.POST("/api/attachments/attachment/:id", h.AttachmentEdit, authenticatedOptions(false))
	server.API.DELETE("/api/attachments/attachment/:id", h.AttachmentDelete, authenticatedOptions(false))

	// Request
	server.API.PUT("/api/action/sync", h.RequestNew, authenticatedOptions(false))
	server.Socket("/api/action/async", h.RequestStream, authenticatedOptions(false))

	// State
	server.API.GET("/api/state", h.State, authenticatedOptions(false))

	// Users
	server.API.GET("/api/users", h.UserList, authenticatedOptions(false))
	server.API.PUT("/api/users/user", h.UserNew, authenticatedOptions(false))
	server.API.GET("/api/users/user/:username", h.UserGet, authenticatedOptions(false))
	server.API.POST("/api/users/user/:username", h.UserEdit, authenticatedOptions(false))
	server.API.POST("/api/users/user/:username/apikey", h.UserResetAPIKey, authenticatedOptions(false))
	server.API.POST("/api/users/reset_password", h.UserResetPassword, authenticatedOptions(true))
	server.API.DELETE("/api/users/user/:username", h.UserDelete, authenticatedOptions(false))

	// Options
	server.API.GET("/api/options", h.OptionsGet, authenticatedOptions(false))
	server.API.POST("/api/options", h.OptionsSet, authenticatedOptions(false))

	// Events
	server.API.GET("/api/events", h.EventsGet, authenticatedOptions(false))

	// System Search
	server.API.POST("/api/search/system", h.SystemSearch, authenticatedOptions(false))

	// Redirect
	server.HTTP.GET("/", v.Redirect, unauthenticatedOptions)

	server.HTTP.GET("/favicon.ico", v.Favicon, unauthenticatedOptions)

	server.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		notFoundFile := path.Join(Directories.Build, "404.html")
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "application/json") {
			json.NewEncoder(w).Encode(web.CommonErrors.NotFound)
		} else if FileExists(notFoundFile) {
			file, err := os.OpenFile(path.Join(Directories.Build, "404.html"), os.O_RDONLY, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			io.CopyBuffer(w, file, nil)
		} else {
			w.Write([]byte("not found"))
		}
	}

	ngRoutes := []string{
		"/hosts",
		"/hosts/host",
		"/hosts/host/:id",
		"/hosts/host/:id/edit",
		"/groups",
		"/groups/group",
		"/groups/group/:id",
		"/groups/group/:id/edit",
		"/scripts",
		"/scripts/script",
		"/scripts/script/:id",
		"/scripts/script/:id/edit",
		"/scripts/script/:id/execute",
		"/schedules",
		"/schedules/schedule",
		"/schedules/schedule/:id",
		"/schedules/schedule/:id/edit",
		"/system/options",
		"/system/users",
		"/system/register",
		"/events",
	}
	for _, route := range ngRoutes {
		server.HTTP.GET(route, v.JavaScript, authenticatedOptions(false))
	}

	server.Start()
}

func unauthorizedHandle(w http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.Header.Get("Accept"), "text/html") {
		w.Header().Add("Location", "/login?unauthorized&redirect="+request.URL.Path)
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
}
