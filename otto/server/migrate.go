package server

import (
	"path"

	"github.com/ecnepsnai/ds"
)

var neededTableVersion = 12

func migrateIfNeeded() {
	currentVersion := State.GetTableVersion()

	if currentVersion == 0 {
		State.SetTableVersion(neededTableVersion + 1)
		log.Debug("Setting default table version to %d", neededTableVersion+1)
		return
	}

	if neededTableVersion-currentVersion > 1 {
		log.Fatal("Refusing to migrate datastore that is too old - follow the supported upgrade path and don't skip versions. Table version %d, required version %d", currentVersion, neededTableVersion)
	}

	i := currentVersion
	for i <= neededTableVersion {
		i++

		results := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "attachment.db"),
			NewPath:   path.Join(Directories.Data, "attachment.db"),
			OldType:   Attachment{},
			NewType:   Attachment{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				attachment, ok := old.(Attachment)
				if !ok {
					panic("Invalid type")
				}
				if attachment.Checksum != "" {
					return attachment, nil
				}

				checksum, err := attachment.GetChecksum()
				if err != nil {
					log.PError("Error calculating checksum of attachment file", map[string]interface{}{
						"path":  attachment.FilePath(),
						"error": err.Error(),
					})
					return nil, nil
				}

				attachment.Checksum = checksum
				return attachment, nil
			},
		})

		if results.Error != nil {
			log.Fatal("Error migrating attachment database: %s", results.Error.Error())
		}
	}

	State.SetTableVersion(i)
}
