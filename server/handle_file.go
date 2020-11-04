package server

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ecnepsnai/web"
)

func (h *handle) FileList(request web.Request) (interface{}, *web.Error) {
	files, err := FileStore.AllFiles()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return files, nil
}

func (h *handle) FileUpload(request web.Request) (interface{}, *web.Error) {
	pathStr := request.HTTP.FormValue("path")
	uidStr := request.HTTP.FormValue("uid")
	gidStr := request.HTTP.FormValue("gid")
	modeStr := request.HTTP.FormValue("mode")

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, web.ValidationError("Invalid uid")
	}
	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		return nil, web.ValidationError("Invalid gid")
	}
	mode, err := strconv.ParseUint(modeStr, 10, 32)
	if err != nil {
		return nil, web.ValidationError("Invalid mode")
	}

	fileUpload, _, err := request.HTTP.FormFile("file")
	if err != nil {
		log.Error("Error getting form file: %s", err.Error())
		return nil, web.CommonErrors.BadRequest
	}

	file := File{
		ID:       NewID(),
		Path:     pathStr,
		UID:      uid,
		GID:      gid,
		Mode:     uint32(mode),
		Created:  time.Now(),
		Modified: time.Now(),
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

	if err := FileStore.Table.Add(file); err != nil {
		log.Error("Error saving script file '%s': %s", file.ID, err.Error())
		return nil, web.CommonErrors.ServerError
	}

	return file, nil
}

func (h *handle) FileGet(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	file, err := FileStore.FileWithID(fileID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return file, nil
}

func (h *handle) FileEdit(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	file, err := FileStore.FileWithID(fileID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	req := editFileParams{}
	if err := request.Decode(&req); err != nil {
		return nil, web.CommonErrors.BadRequest
	}

	file, err = FileStore.EditFile(fileID, req)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return file, nil
}

func (h *handle) FileDelete(request web.Request) (interface{}, *web.Error) {
	fileID := request.Params.ByName("id")

	if err := FileStore.DeleteFile(fileID); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return true, nil
}
