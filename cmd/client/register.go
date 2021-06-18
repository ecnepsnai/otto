package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	rlog "log"
	"net"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ecnepsnai/osquery"
	"github.com/ecnepsnai/otto"
)

var registerProperties otto.RegisterRequestProperties

func tryAutoRegister() {
	env := envMap()

	host, present := env["REGISTER_HOST"]
	if !present {
		return
	}

	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		panic("Invalid value for variable REGISTER_HOST")
	}
	// Trim trailing slash if there is one
	host = strings.TrimSuffix(host, "/")
	hostNoProto := strings.ReplaceAll(strings.ReplaceAll(host, "http://", ""), "https://", "")

	key, present := env["REGISTER_KEY"]
	if !present {
		return
	}
	if key == "" {
		panic("Invalid value for variable REGISTER_KEY")
	}

	// Disable TLS Verification
	noTLSVerify := env["REGISTER_NO_TLS_VERIFY"] == "1"

	// Exit when Finished
	exitWhenFinished := true
	if _, p := env["REGISTER_DONT_EXIT_ON_FINISH"]; p {
		exitWhenFinished = false
	}

	// Client Port
	port := uint32(12444)
	portStr, present := env["OTTO_CLIENT_PORT"]
	if present {
		p, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			panic("Invalid value for variable OTTO_CLIENT_PORT")
		}
		port = uint32(p)
	}

	rlog.Printf("Will attempt to register client with otto server '%s'", host)
	if noTLSVerify {
		rlog.Printf("WARNING - TLS Verification disabled, this is insecure and should not be used in production")
	}
	if strings.HasSuffix(host, "http://") {
		rlog.Printf("WARNING - Not using TLS, registration key will be sent in plain-text")
	}

	// Get the current uid and gid for the defaults
	uid, gid := getUIDandGID()

	// Make the request
	request := otto.RegisterRequest{
		Key:        key,
		Port:       port,
		Properties: registerProperties,
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic("Error forming JSON request")
	}

	// Prepare the request
	url := fmt.Sprintf("%s/api/register", host)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		panic("Error forming HTTP request")
	}

	req.Header.Add("X-OTTO-PROTO-VERSION", fmt.Sprintf("%d", otto.ProtocolVersion))

	tr := &http.Transport{}
	if noTLSVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	httpClient := &http.Client{Transport: tr}

	// Request registration
	rlog.Printf("HTTP PUT -> %s", url)
	response, err := httpClient.Do(req)
	if err != nil {
		rlog.Fatalf("Error connecting to otto server: %s", err.Error())
	}
	rlog.Printf("HTTP %d", response.StatusCode)
	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Request.Body)
		rlog.Fatalf("HTTP Error %d from otto server. Response: %s", response.StatusCode, buf.String())
	}

	// Parse the response
	type responseType struct {
		Data otto.RegisterResponse `json:"data"`
	}
	registerResponse := responseType{}
	if err := json.NewDecoder(response.Body).Decode(&registerResponse); err != nil {
		rlog.Fatalf("Error decoding response body: %s", err.Error())
	}
	if registerResponse.Data.PSK == "" {
		rlog.Fatalf("No PSK returned from otto server")
	}

	var listenAddr = fmt.Sprintf("0.0.0.0:%d", port)
	var allowFrom = "0.0.0.0/0"
	// Check if we connected to the server using IPv6
	if len(getOutboundIP(hostNoProto)) == 16 {
		listenAddr = fmt.Sprintf("[::]:%d", port)
		allowFrom = "::/0"
	}

	// Save the config
	conf := clientConfig{
		PSK:        registerResponse.Data.PSK,
		LogPath:    ".",
		DefaultUID: uid,
		DefaultGID: gid,
		Path:       env["PATH"],
		ListenAddr: listenAddr,
		AllowFrom:  allowFrom,
	}
	config = &conf
	if err := saveNewConfig(conf); err != nil {
		rlog.Fatalf("Error saving config file: %s", err.Error())
	}
	rlog.Printf("Successfully registered with otto server '%s', configuration: %+v", host, conf)

	for _, script := range registerResponse.Data.Scripts {
		rlog.Printf("Executing first-run script: %s", script.Name)
		devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, os.ModePerm)
		runScript(devNull, script, nil)
		devNull.Close()
	}

	if exitWhenFinished {
		os.Exit(0)
	}
}

func loadRegisterProperties() {
	// Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		panic("Error getting system hostname: " + err.Error())
	}

	info, err := osquery.Get()
	if err != nil {
		panic("Error getting system information: " + err.Error())
	}

	registerProperties = otto.RegisterRequestProperties{
		Hostname:            hostname,
		KernelName:          info.Kernel,
		KernelVersion:       info.KernelVersion,
		DistributionName:    info.Variant,
		DistributionVersion: info.VariantVersion,
	}
}

// envMap return a map of all environment variables
func envMap() map[string]string {
	results := map[string]string{}
	for _, env := range os.Environ() {
		components := strings.SplitN(env, "=", 2)
		results[components[0]] = components[1]
	}
	return results
}

// getUIDandGID return the current users UID and primary group GID
func getUIDandGID() (uid, gid uint32) {
	me, err := user.Current()
	if err != nil {
		panic("Unable to get current user")
	}

	u, err := strconv.ParseUint(me.Uid, 10, 32)
	if err != nil {
		panic("Uid is not a number: " + me.Uid)
	}
	g, err := strconv.ParseUint(me.Gid, 10, 32)
	if err != nil {
		panic("Gid is not a number: " + me.Gid)
	}

	uid = uint32(u)
	gid = uint32(g)
	return
}

// compareUIDandGID compare the given UID and GID to that of the running user
func compareUIDandGID(uid, gid uint32) bool {
	u, g := getUIDandGID()
	return uid == u && gid == g
}

// getOutboundIP get the IP address used to connect to a remote destination
func getOutboundIP(host string) net.IP {
	// Because we're using udp there's no handshake or connection actually happening here
	conn, err := net.Dial("udp", host)
	if err != nil {
		rlog.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
