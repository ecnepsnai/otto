package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ecnepsnai/web"
)

func (h *handle) AttachmentList(request web.Request) (interface{}, *web.Error) {
	return AttachmentStore.AllAttachments(), nil
}

func (h *handle) AttachmentUpload(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	pathStr := request.HTTP.FormValue("Path")
	inheritStr := request.HTTP.FormValue("Inherit")
	inherit := inheritStr == "true"
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

	req := newAttachmentParameters{
		Data:     fileUpload,
		Path:     pathStr,
		Name:     info.Filename,
		MimeType: info.Header.Get("Content-Type"),
		Owner: RunAs{
			Inherit: inherit,
			UID:     uint32(uid),
			GID:     uint32(gid),
		},
		Mode: uint32(mode),
		Size: uint64(info.Size),
	}

	attachment, erro := AttachmentStore.NewAttachment(req)
	if erro != nil {
		if erro.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(erro.Message)
	}
	EventStore.AttachmentAdded(attachment, session.Username)

	return attachment, nil
}

func (h *handle) AttachmentGet(request web.Request) (interface{}, *web.Error) {
	attachmentID := request.Parameters["id"]
	return AttachmentStore.AttachmentWithID(attachmentID), nil
}

func (v *view) AttachmentDownload(request web.Request, writer web.Writer) (response web.Response) {
	attachmentID := request.Parameters["id"]

	attachment := AttachmentStore.AttachmentWithID(attachmentID)
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

	attachmentID := request.Parameters["id"]

	pathStr := request.HTTP.FormValue("Path")
	inheritStr := request.HTTP.FormValue("Inherit")
	inherit := inheritStr == "true"
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
	if err != nil && err.Error() != "http: no such file" {
		log.Error("Error getting form file: %s", err.Error())
		return nil, web.CommonErrors.BadRequest
	}

	req := editAttachmentParams{
		Path: pathStr,
		Owner: RunAs{
			Inherit: inherit,
			UID:     uint32(uid),
			GID:     uint32(gid),
		},
		Mode: uint32(mode),
	}

	if info != nil {
		req.Data = fileUpload
		req.Name = info.Filename
		req.MimeType = info.Header.Get("Content-Type")
		req.Size = uint64(info.Size)
	}

	attachment, erro := AttachmentStore.EditAttachment(attachmentID, req)
	if erro != nil {
		if erro.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(erro.Message)
	}

	EventStore.AttachmentModified(attachment, session.Username)

	return attachment, nil
}

func (h *handle) AttachmentDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	attachmentID := request.Parameters["id"]

	if err := AttachmentStore.DeleteAttachment(attachmentID); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.AttachmentDeleted(attachmentID, session.Username)

	return true, nil
}
