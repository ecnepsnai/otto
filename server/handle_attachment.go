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
	attachments, err := AttachmentStore.AllAttachments()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return attachments, nil
}

func (h *handle) AttachmentUpload(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

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

	attachment := Attachment{
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

	f, err := os.OpenFile(attachment.FilePath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Error("Error opening attachment file '%s': %s", attachment.FilePath(), err.Error())
		return nil, web.CommonErrors.ServerError
	}
	defer f.Close()

	if _, err := io.Copy(f, fileUpload); err != nil {
		log.Error("Error writing attachment file '%s': %s", attachment.FilePath(), err.Error())
		return nil, web.CommonErrors.ServerError
	}

	if err := AttachmentStore.Table.Add(attachment); err != nil {
		log.Error("Error saving script attachment '%s': %s", attachment.ID, err.Error())
		return nil, web.CommonErrors.ServerError
	}

	EventStore.AttachmentAdded(&attachment, session.Username)

	return attachment, nil
}

func (h *handle) AttachmentGet(request web.Request) (interface{}, *web.Error) {
	attachmentID := request.Params.ByName("id")

	attachment, err := AttachmentStore.AttachmentWithID(attachmentID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return attachment, nil
}

func (v *view) AttachmentDownload(request web.Request, writer web.Writer) (response web.Response) {
	attachmentID := request.Params.ByName("id")

	attachment, erro := AttachmentStore.AttachmentWithID(attachmentID)
	if erro != nil {
		if erro.Server {
			response.Status = 500
			return
		}
		response.Status = 400
		return
	}

	f, err := os.OpenFile(attachment.FilePath(), os.O_RDONLY, 0644)
	if err != nil {
		response.Status = 500
		return
	}
	response.ContentType = attachment.MimeType
	response.Headers = map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", attachment.Name),
	}
	response.Reader = f
	return
}

func (h *handle) AttachmentEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	attachmentID := request.Params.ByName("id")

	attachment, err := AttachmentStore.AttachmentWithID(attachmentID)
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

	attachment, err = AttachmentStore.EditAttachment(attachmentID, req)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.AttachmentModified(attachment, session.Username)

	return attachment, nil
}

func (h *handle) AttachmentDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	attachmentID := request.Params.ByName("id")

	if err := AttachmentStore.DeleteAttachment(attachmentID); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.AttachmentDeleted(attachmentID, session.Username)

	return true, nil
}
