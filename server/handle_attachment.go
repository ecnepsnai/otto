package server

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ecnepsnai/web"
)

func (h *handle) AttachmentList(request web.Request) (interface{}, *web.Error) {
	files, err := AttachmentStore.AllAttachments()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return files, nil
}

func (h *handle) AttachmentUpload(request web.Request) (interface{}, *web.Error) {
	pathStr := request.HTTP.FormValue("Path")
	uidStr := request.HTTP.FormValue("UID")
	gidStr := request.HTTP.FormValue("GID")
	modeStr := request.HTTP.FormValue("Mode")

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, web.ValidationError("Invalid uid '%s'", uidStr)
	}
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		return nil, web.ValidationError("Invalid gid '%s'", gidStr)
	}
	mode, err := strconv.ParseUint(modeStr, 10, 32)
	if err != nil {
		return nil, web.ValidationError("Invalid mode '%s'", modeStr)
	}

	fileUpload, info, err := request.HTTP.FormFile("file")
	if err != nil {
		log.Error("Error getting form file: %s", err.Error())
		return nil, web.CommonErrors.BadRequest
	}

	file := Attachment{
		ID:       newPlainID(),
		Path:     pathStr,
		Name:     info.Filename,
		MimeType: info.Header.Get("Content-Type"),
		UID:      uid,
		GID:      gid,
		Mode:     uint32(mode),
		Created:  time.Now(),
		Modified: time.Now(),
		Size:     uint64(info.Size),
	}

	f, err := os.OpenFile(file.FilePath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Error("Error opening script file '%s': %s", file.FilePath(), err.Error())
		return nil, web.CommonErrors.ServerError
	}
	defer f.Close()

	if _, err := io.Copy(f, fileUpload); err != nil {
		log.Error("Error writing script file '%s': %s", file.FilePath(), err.Error())
		return nil, web.CommonErrors.ServerError
	}

	if err := AttachmentStore.Table.Add(file); err != nil {
		log.Error("Error saving script file '%s': %s", file.ID, err.Error())
		return nil, web.CommonErrors.ServerError
	}

	return file, nil
}

func (h *handle) AttachmentGet(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	file, err := AttachmentStore.AttachmentWithID(fileID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return file, nil
}

func (v *view) AttachmentDownload(request web.Request, writer web.Writer) (response web.Response) {
	fileID := request.Params.ByName("id")

	file, erro := AttachmentStore.AttachmentWithID(fileID)
	if erro != nil {
		if erro.Server {
			response.Status = 500
			return
		}
		response.Status = 400
		return
	}

	f, err := os.OpenFile(file.FilePath(), os.O_RDONLY, 0644)
	if err != nil {
		response.Status = 500
		return
	}
	response.ContentType = file.MimeType
	response.Headers = map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", file.Name),
	}
	response.Reader = f
	return
}

func (h *handle) AttachmentEdit(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	file, err := AttachmentStore.AttachmentWithID(fileID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	req := editAttachmentParams{}
	if err := request.Decode(&req); err != nil {
		return nil, web.CommonErrors.BadRequest
	}

	file, err = AttachmentStore.EditAttachment(fileID, req)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return file, nil
}

func (h *handle) AttachmentDelete(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	if err := AttachmentStore.DeleteAttachment(fileID); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return true, nil
}
