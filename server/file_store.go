package server

import "time"

func (s fileStoreObject) AllFiles() ([]File, *Error) {
	objects, err := s.Table.GetAll(nil)
	if err != nil {
		log.Error("Error getting all script files: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objects == nil || len(objects) == 0 {
		return []File{}, nil
	}

	files := make([]File, len(objects))
	for i, obj := range objects {
		file, k := obj.(File)
		if !k {
			log.Error("Object is not of type 'File'")
			return []File{}, ErrorServer("incorrect type")
		}
		files[i] = file
	}

	return files, nil
}

func (s fileStoreObject) AllFilesForScript(scriptID string) ([]File, *Error) {
	objects, err := s.Table.GetIndex("ScriptID", scriptID, nil)
	if err != nil {
		log.Error("Error getting all script files for script '%s': %s", scriptID, err.Error())
		return nil, ErrorFrom(err)
	}
	if objects == nil || len(objects) == 0 {
		return []File{}, nil
	}

	files := make([]File, len(objects))
	for i, obj := range objects {
		file, k := obj.(File)
		if !k {
			log.Error("Object is not of type 'File'")
			return []File{}, ErrorServer("incorrect type")
		}
		files[i] = file
	}

	return files, nil
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
	file.Modified = time.Now()

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

func (s fileStoreObject) Cleanup() *Error {
	filesWithScripts := map[string]bool{}
	files, err := s.AllFiles()
	if err != nil {
		return err
	}
	scripts, err := ScriptStore.AllScripts()
	if err != nil {
		return err
	}
	for _, file := range files {
		for _, script := range scripts {
			if StringSliceContains(file.ID, script.FileIDs) {
				filesWithScripts[file.ID] = true
				break
			}
		}
	}

	for _, file := range files {
		if filesWithScripts[file.ID] {
			continue
		}

		if time.Since(file.Modified) > 1*time.Hour {
			log.Warn("Removing orphaned script file '%s'", file.ID)
			if err := s.DeleteFile(file.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
