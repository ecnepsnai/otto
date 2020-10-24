package server

func (s fileStoreObject) AllFiles() ([]File, *Error) {
	objs, err := s.Table.GetAll(nil)
	if err != nil {
		log.Error("Error getting all script files: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []File{}, nil
	}

	filess := make([]File, len(objs))
	for i, obj := range objs {
		files, k := obj.(File)
		if !k {
			log.Error("Object is not of type 'File'")
			return []File{}, ErrorServer("incorrect type")
		}
		filess[i] = files
	}

	return filess, nil
}

func (s fileStoreObject) AllFilesForScript(scriptID string) ([]File, *Error) {
	objs, err := s.Table.GetIndex("ScriptID", scriptID, nil)
	if err != nil {
		log.Error("Error getting all script files for script '%s': %s", scriptID, err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []File{}, nil
	}

	filess := make([]File, len(objs))
	for i, obj := range objs {
		files, k := obj.(File)
		if !k {
			log.Error("Object is not of type 'File'")
			return []File{}, ErrorServer("incorrect type")
		}
		filess[i] = files
	}

	return filess, nil
}

func (s fileStoreObject) FileWithID(id string) (*File, *Error) {
	obj, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting script file with ID '%s': %s", id, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	file, k := obj.(File)
	if !k {
		log.Error("Object is not of type 'File'")
		return nil, ErrorServer("incorrect type")
	}

	return &file, nil
}

type editFileParams struct {
	Path string
	UID  int
	GID  int
	Mode uint32
}

func (s fileStoreObject) EditFile(id string, params editFileParams) (*File, *Error) {
	file, err := s.FileWithID(id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrorUser("No script with ID")
	}

	file.Path = params.Path
	file.UID = params.UID
	file.GID = params.GID
	file.Mode = params.Mode

	if err := s.Table.Update(*file); err != nil {
		log.Error("Error updating script file '%s': %s", file.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	return file, nil
}

func (s fileStoreObject) DeleteFile(id string) *Error {
	file, err := s.FileWithID(id)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrorUser("No script with ID")
	}

	if err := s.Table.Delete(*file); err != nil {
		log.Error("Error deleting script file '%s': %s", file.ID, err.Error())
		return ErrorFrom(err)
	}

	return nil
}
