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
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/ecnepsnai/otto"
)

func tryAutoRegister() {
	env := envMap()

	// Host
	host, present := env["REGISTER_HOST"]
	if !present {
		return
	}

	// PSK
	psk, present := env["REGISTER_PSK"]
	if !present {
		return
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

	// Get the uname
	unameB, err := exec.Command("uname", "-a").CombinedOutput()
	if err != nil {
		panic("Unable to get uname of this host: " + err.Error())
	}
	uname := string(unameB)

	// Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		panic("Unable to get hostname of this host: " + err.Error())
	}

	// Get the current uid and gid for the defaults
	uid, gid := getUIDandGID()
	// Get the local IP
	localIP := getOutboundIP().String()

	// Make the request
	request := otto.RegisterRequest{
		Address:  localIP,
		PSK:      psk,
		Uname:    uname,
		Hostname: hostname,
		Port:     port,
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

	// Save the config
	conf := clientConfig{
		PSK:        registerResponse.Data.PSK,
		LogPath:    ".",
		DefaultUID: uid,
		DefaultGID: gid,
		Path:       env["PATH"],
	}
	f, err := os.OpenFile("otto_client.conf", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		rlog.Fatalf("Error opening config file: %s", err.Error())
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(conf); err != nil {
		rlog.Fatalf("Error encoding options: %s", err.Error())
	}
	rlog.Printf("Successfully registered with otto server '%s'", host)
	if exitWhenFinished {
		os.Exit(0)
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

// getOutboundIP get the IP address used to connect to a remote destination
func getOutboundIP() net.IP {
	// Because we're using udp there's no handshake or connection actually happening here
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		rlog.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
