package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/ecnepsnai/otto/server"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/web"
)

func startServer(stop chan uint8) error {
	if err := os.MkdirAll(path.Join(dataDir, "server"), 7644); err != nil {
		return err
	}

	out, err := os.OpenFile(path.Join(dataDir, "server", "stdout.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command(path.Join(testDir, "server"), "-v", "--data-dir", path.Join(dataDir, "server"))
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Dir = path.Join(dataDir, "server")
	cmd.Env = append(cmd.Env, "UNSAFE_ALLOW_DEFAULT_PASSWORD=1")
	if err := cmd.Start(); err != nil {
		return err
	}

	process := cmd.Process
	_ = <-stop
	process.Signal(syscall.SIGKILL)
	fmt.Println("Stopping server")
	out.Sync()
	out.Close()
	return nil
}

func getApiKey() (*string, error) {
	var apiKey string

	err := logTask("Generate API key", func() error {
		// Login
		type credentialsType struct {
			Username string
			Password string
		}
		credentials := credentialsType{
			Username: "admin",
			Password: "admin",
		}
		body, err := json.Marshal(credentials)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/login", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		req, err = http.NewRequest("POST", "http://127.0.0.1:8080/api/users/user/admin/apikey", nil)
		if err != nil {
			return err
		}
		for _, cookie := range resp.Cookies() {
			req.AddCookie(cookie)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		type tResponse struct {
			Data string `json:"data"`
		}

		apiKeyResponse := tResponse{}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &apiKeyResponse); err != nil {
			return err
		}

		apiKey = apiKeyResponse.Data
		if apiKey == "" {
			return fmt.Errorf("No API key generated")
		}
		return nil
	})
	return &apiKey, err
}

func newServerRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, "http://127.0.0.1:8080"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-OTTO-USERNAME", "admin")
	req.Header.Add("X-OTTO-API-KEY", serverApiKey)
	return req, nil
}

func getGroup() (*string, error) {
	// Get the default group
	req, err := newServerRequest("GET", "/api/groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	type tGroupResponse struct {
		Data  []server.Group
		Error *web.Error
	}
	groupResponse := tGroupResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&groupResponse); err != nil {
		return nil, err
	}
	if groupResponse.Error != nil {
		return nil, fmt.Errorf("%s", groupResponse.Error.Message)
	}
	id := groupResponse.Data[0].ID
	return &id, nil
}

func getScript() (*string, error) {
	// Get the default script
	req, err := newServerRequest("GET", "/api/scripts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	type tScriptResponse struct {
		Data  []server.Script
		Error *web.Error
	}
	scriptResponse := tScriptResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&scriptResponse); err != nil {
		return nil, err
	}
	if scriptResponse.Error != nil {
		return nil, fmt.Errorf("%s", scriptResponse.Error.Message)
	}
	id := scriptResponse.Data[0].ID
	return &id, nil
}

func getHost() (*string, error) {
	// Get the default host
	req, err := newServerRequest("GET", "/api/hosts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	type tHostResponse struct {
		Data  []server.Host
		Error *web.Error
	}
	hostResponse := tHostResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&hostResponse); err != nil {
		return nil, err
	}
	if hostResponse.Error != nil {
		return nil, fmt.Errorf("%s", hostResponse.Error.Message)
	}
	id := hostResponse.Data[0].ID
	return &id, nil
}

func enableAutoRegister() error {
	return logTask("Configure automatic registration", func() error {
		groupId, err := getGroup()
		if err != nil {
			return err
		}

		// Configure automatic registration
		registerOptions := server.OptionsRegister{
			Enabled:        true,
			Key:            registerKey,
			DefaultGroupID: *groupId,
		}
		body, err := json.Marshal(registerOptions)
		if err != nil {
			return err
		}

		req, err := newServerRequest("POST", "/api/register/options", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		type tOptionsResponse struct {
			Data  bool
			Error *web.Error
		}
		optionsResponse := tOptionsResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&optionsResponse); err != nil {
			return err
		}
		if optionsResponse.Error != nil {
			return fmt.Errorf("%s", optionsResponse.Error.Message)
		}
		if !optionsResponse.Data {
			return fmt.Errorf("unable to set register options")
		}
		return nil
	})
}

func createAttachment(name string, afterScript bool) (*string, error) {
	multipartBuffer := &bytes.Buffer{}
	m := multipart.NewWriter(multipartBuffer)
	m.WriteField("Path", path.Join(dataDir, "server", name))
	m.WriteField("Inherit", "true")
	m.WriteField("Mode", "0644")
	m.WriteField("UID", "0")
	m.WriteField("GID", "0")
	if afterScript {
		m.WriteField("AfterScript", "true")
	}
	if w, err := m.CreateFormFile("file", name); err != nil {
		w.Write([]byte("1\n"))
	}
	m.Close()

	req, err := newServerRequest("PUT", "/api/attachments", multipartBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	type tAttachmentResponse struct {
		Data  *server.Attachment
		Error *web.Error
	}
	attachmentResponse := tAttachmentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&attachmentResponse); err != nil {
		return nil, err
	}
	if attachmentResponse.Error != nil {
		return nil, fmt.Errorf("%s", attachmentResponse.Error.Message)
	}
	if attachmentResponse.Data == nil {
		return nil, fmt.Errorf("no attachment results returnred")
	}
	return &attachmentResponse.Data.ID, nil
}

func createScript() error {
	return logTask("Create test script", func() error {
		beforeAttachmentId, err := createAttachment("upload1.txt", false)
		if err != nil {
			return err
		}
		afterAttachmentId, err := createAttachment("upload2.txt", true)
		if err != nil {
			return err
		}

		scriptId, err := getScript()
		if err != nil {
			return err
		}

		type editScriptParameters struct {
			Name             string
			Executable       string
			Script           string
			Environment      []environ.Variable
			RunAs            server.RunAs
			WorkingDirectory string
			AfterExecution   string
			AttachmentIDs    []string
			RunLevel         int
		}
		body, err := json.Marshal(editScriptParameters{
			Name:       "OttoTest",
			Executable: "/bin/bash",
			Script: `#!/bin/bash
set -e

if [[ ! -f "${BEFORE_UPLOAD_FILE}" ]]; then
    echo "${BEFORE_UPLOAD_FILE} does not exist"
	exit 1
fi
if [[ -f "${AFTER_UPLOAD_FILE}" ]]; then
    echo "${AFTER_UPLOAD_FILE} exists when it should not"
	exit 1
fi

echo 1 > ${TEST_FILE}
echo "HELLO!!!!!!!!!"
`,
			Environment: []environ.Variable{
				{
					Key:   "TEST_FILE",
					Value: path.Join(dataDir, "server", "test_passed.txt"),
				},
				{
					Key:   "BEFORE_UPLOAD_FILE",
					Value: path.Join(dataDir, "server", "upload1.txt"),
				},
				{
					Key:   "AFTER_UPLOAD_FILE",
					Value: path.Join(dataDir, "server", "upload2.txt"),
				},
			},
			RunAs: server.RunAs{
				Inherit: true,
			},
			RunLevel: server.ScriptRunLevelReadOnly,
			AttachmentIDs: []string{
				*beforeAttachmentId,
				*afterAttachmentId,
			},
		})
		if err != nil {
			return err
		}

		req, err := newServerRequest("POST", "/api/scripts/script/"+*scriptId, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		type tScriptResponse struct {
			Data  *server.Script
			Error *web.Error
		}
		scriptResponse := tScriptResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&scriptResponse); err != nil {
			return err
		}
		if scriptResponse.Error != nil {
			return fmt.Errorf("%s", scriptResponse.Error.Message)
		}
		if scriptResponse.Data == nil {
			return fmt.Errorf("no script results returnred")
		}
		return nil
	})
}

func runScript() error {
	return logTask("Run script", func() error {
		// Remove test files
		os.Remove(path.Join(dataDir, "server", "test_passed.txt"))
		os.Remove(path.Join(dataDir, "server", "upload1.txt"))
		os.Remove(path.Join(dataDir, "server", "upload2.txt"))

		// Get the default host
		req, err := newServerRequest("GET", "/api/hosts", nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		type tHostResponse struct {
			Data  []server.Host
			Error *web.Error
		}
		hostResponse := tHostResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&hostResponse); err != nil {
			return err
		}
		if hostResponse.Error != nil {
			return fmt.Errorf("%s", hostResponse.Error.Message)
		}
		defaultHostId := hostResponse.Data[0].ID

		// Get the default script
		req, err = newServerRequest("GET", "/api/scripts", nil)
		if err != nil {
			return err
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		type tScriptResponse struct {
			Data  []server.Script
			Error *web.Error
		}
		scriptResponse := tScriptResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&scriptResponse); err != nil {
			return err
		}
		if scriptResponse.Error != nil {
			return fmt.Errorf("%s", scriptResponse.Error.Message)
		}
		defaultScriptId := scriptResponse.Data[0].ID

		type requestParamsType struct {
			HostID   string
			Action   string
			ScriptID string
		}
		requestParams := requestParamsType{
			HostID:   defaultHostId,
			Action:   server.AgentActionRunScript,
			ScriptID: defaultScriptId,
		}
		body, err := json.Marshal(requestParams)
		if err != nil {
			return err
		}

		req, err = newServerRequest("PUT", "/api/action/sync", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		type tSyncResponse struct {
			Data  *server.ScriptResult
			Error *web.Error
		}
		syncResponse := tSyncResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&syncResponse); err != nil {
			return err
		}
		if syncResponse.Error != nil {
			return fmt.Errorf("%s", syncResponse.Error.Message)
		}
		if syncResponse.Data == nil {
			return fmt.Errorf("no script results returnred")
		}
		if !syncResponse.Data.Result.Success {
			return fmt.Errorf("unsuccessful script execution")
		}

		if _, err := os.Stat(path.Join(dataDir, "server", "test_passed.txt")); err != nil {
			return fmt.Errorf("test file not created")
		}
		if _, err := os.Stat(path.Join(dataDir, "server", "upload1.txt")); err != nil {
			return fmt.Errorf("attachment file not created")
		}
		if _, err := os.Stat(path.Join(dataDir, "server", "upload2.txt")); err != nil {
			return fmt.Errorf("attachment file not created")
		}

		return nil
	})
}

func rotateIdentity() error {
	return logTask("Rotating identity", func() error {
		hostId, err := getHost()
		if err != nil {
			return err
		}

		req, err := newServerRequest("POST", "/api/hosts/host/"+*hostId+"/id/rotate", nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		type tResponse struct {
			Data  *bool
			Error *web.Error
		}
		results := tResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			return err
		}
		if results.Error != nil {
			return fmt.Errorf("%s", results.Error.Message)
		}
		return nil
	})
}
