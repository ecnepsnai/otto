package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	rlog "log"
	"net"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ecnepsnai/osquery"
	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
)

var registerProperties otto.RegisterRequestProperties

func tryAutoRegister() {
	host := os.Getenv("REGISTER_HOST")
	if host == "" {
		return
	}

	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		panic("Invalid value for variable REGISTER_HOST")
	}
	// Trim trailing slash if there is one
	host = strings.TrimSuffix(host, "/")
	hostNoProto := strings.ReplaceAll(strings.ReplaceAll(host, "http://", ""), "https://", "")

	key := os.Getenv("REGISTER_KEY")
	if key == "" {
		return
	}
	if key == "" {
		panic("Invalid value for variable REGISTER_KEY")
	}

	// Disable TLS Verification
	noTLSVerify := os.Getenv("REGISTER_NO_TLS_VERIFY") == "1"

	// Exit when Finished
	exitWhenFinished := os.Getenv("REGISTER_DONT_EXIT_ON_FINISH") == ""

	// Client Port
	port := uint32(12444)
	portStr := os.Getenv("OTTO_CLIENT_PORT")
	if portStr != "" {
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

	if err := generateIdentity(); err != nil {
		panic("Error generating client identity: " + err.Error())
	}
	signer, err := loadClientIdentity()
	if err != nil {
		panic("Error reading client identity: " + err.Error())
	}
	rlog.Printf("Client identity: %s", base64.RawURLEncoding.EncodeToString(signer.PublicKey().Marshal()))

	// Make the request
	request := otto.RegisterRequest{
		ClientIdentity: base64.StdEncoding.EncodeToString(signer.PublicKey().Marshal()),
		Port:           port,
		Properties:     registerProperties,
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic("Error forming JSON request")
	}
	encryptedData, err := secutil.Encryption.AES_256_GCM.Encrypt(data, key)
	if err != nil {
		panic("Error encrypting request")
	}

	// Prepare the request
	url := fmt.Sprintf("%s/api/register", host)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(encryptedData))
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

	registerResponse := otto.RegisterResponse{}

	encryptedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		panic("Error reading reply")
	}
	decryptedResponse, err := secutil.Encryption.AES_256_GCM.Decrypt(encryptedResponse, key)
	if err != nil {
		panic("Error decrypting reply")
	}
	if err := json.Unmarshal(decryptedResponse, &registerResponse); err != nil {
		rlog.Fatalf("Error decoding response body: %s", err.Error())
	}
	if registerResponse.ServerIdentity == "" {
		rlog.Fatalf("No server identity returned from otto server")
	}
	rlog.Printf("Server identity: %s", registerResponse.ServerIdentity)

	var listenAddr = fmt.Sprintf("0.0.0.0:%d", port)
	// Check if we connected to the server using IPv6
	if len(getOutboundIP(hostNoProto)) == 16 {
		listenAddr = fmt.Sprintf("[::]:%d", port)
	}

	// Save the config
	conf := clientConfig{
		IdentityPath:   otto_IDENTITY_FILE_NAME,
		ServerIdentity: registerResponse.ServerIdentity,
		LogPath:        ".",
		DefaultUID:     uid,
		DefaultGID:     gid,
		Path:           os.Getenv("PATH"),
		ListenAddr:     listenAddr,
		AllowFrom:      defaultConfig().AllowFrom,
	}
	config = &conf
	if err := saveNewConfig(conf); err != nil {
		rlog.Fatalf("Error saving config file: %s", err.Error())
	}
	rlog.Printf("Successfully registered with otto server '%s', configuration: %+v", host, conf)

	for _, script := range registerResponse.Scripts {
		rlog.Printf("Executing first-run script: %s", script.Name)
		runScript(nil, script, nil)
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
