package server

import (
	"io"
	"os"
	"time"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
)

func (s attachmentStoreObject) AllAttachments() (attachments []Attachment) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		attachments = s.allAttachments(tx)
		return nil
	})
	return
}

func (s attachmentStoreObject) allAttachments(tx ds.IReadTransaction) []Attachment {
	objects, err := tx.GetAll(nil)
	if err != nil {
		log.Error("Error listing attachments: error='%s'", err.Error())
		return []Attachment{}
	}
	if len(objects) == 0 {
		return []Attachment{}
	}

	files := make([]Attachment, len(objects))
	for i, objects := range objects {
		file, k := objects.(Attachment)
		if !k {
			log.Fatal("Error listing attachments: error='%s'", "invalid type")
		}
		files[i] = file
	}

	return files
}

func (s attachmentStoreObject) AllAttachmentsForScript(scriptID string) (attachments []Attachment) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		attachments = s.allAttachmentsForScript(tx, scriptID)
		return nil
	})
	return
}

func (s attachmentStoreObject) allAttachmentsForScript(tx ds.IReadTransaction, scriptID string) []Attachment {
	script := ScriptStore.ScriptWithID(scriptID)
	if script == nil {
		return []Attachment{}
	}
	count := len(script.AttachmentIDs)
	if count == 0 {
		return []Attachment{}
	}

	attachments := make([]Attachment, count)
	for i, id := range script.AttachmentIDs {
		attachment := s.attachmentWithID(tx, id)
		if attachment == nil {
			log.Error("Script references non-existant attachment: script_id='%s' attachment_id='%s'", scriptID, id)
			continue
		}
		attachments[i] = *attachment
	}

	return attachments
}

func (s attachmentStoreObject) AttachmentWithID(id string) (attachment *Attachment) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		attachment = s.attachmentWithID(tx, id)
		return nil
	})
	return
}

func (s attachmentStoreObject) attachmentWithID(tx ds.IReadTransaction, id string) *Attachment {
	object, err := tx.Get(id)
	if err != nil {
		log.Error("Error getting attachment: id='%s' error='%s'", id, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	file, k := object.(Attachment)
	if !k {
		log.Fatal("Error getting attachment: id='%s' error='%s'", id, "invalid type")
	}

	return &file
}

type newAttachmentParameters struct {
	Data        io.Reader
	Path        string `min:"1"`
	Name        string `min:"1"`
	MimeType    string `min:"1"`
	Owner       RunAs
	Mode        uint32
	Size        uint64
	AfterScript bool
}

func (s attachmentStoreObject) NewAttachment(params newAttachmentParameters) (attachment *Attachment, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		attachment, err = s.newAttachment(tx, params)
		return nil
	})
	return
}

func (s attachmentStoreObject) newAttachment(tx ds.IReadWriteTransaction, params newAttachmentParameters) (*Attachment, *Error) {
	if err := limits.Check(params); err != nil {
		return nil, ErrorUser(err.Error())
	}

	attachment := Attachment{
		ID:          newPlainID(),
		Path:        params.Path,
		Name:        params.Name,
		MimeType:    params.MimeType,
		Owner:       params.Owner,
		Mode:        params.Mode,
		Created:     time.Now(),
		Modified:    time.Now(),
		Size:        params.Size,
		AfterScript: params.AfterScript,
	}

	f, err := os.OpenFile(attachment.FilePath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.PError("Error writing attachment file", map[string]interface{}{
			"attachment": attachment.FilePath(),
			"error":      err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	if _, err := io.Copy(f, params.Data); err != nil {
		log.PError("Error writing attachment file", map[string]interface{}{
			"attachment": attachment.FilePath(),
			"error":      err.Error(),
		})
		f.Close()
		return nil, ErrorFrom(err)
	}
	f.Close()

	checksum, err := getFileSHA256Checksum(attachment.FilePath())
	if err != nil {
		log.PError("Error calculating attachment checksum", map[string]interface{}{
			"attachment": attachment.FilePath(),
			"error":      err.Error(),
		})
		return nil, ErrorFrom(err)
	}
	attachment.Checksum = checksum

	if err := tx.Add(attachment); err != nil {
		log.PError("Error saving attachment", map[string]interface{}{
			"attachment": attachment.ID,
			"error":      err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	return &attachment, nil
}

type editAttachmentParams struct {
	Data        io.Reader
	Path        string `min:"1"`
	Name        string `min:"1"`
	MimeType    string `min:"1"`
	Owner       RunAs
	Mode        uint32
	Size        uint64
	AfterScript bool
}

func (s attachmentStoreObject) EditAttachment(id string, params editAttachmentParams) (attachment *Attachment, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		attachment, err = s.editAttachment(tx, id, params)
		return nil
	})
	return
}

func (s attachmentStoreObject) editAttachment(tx ds.IReadWriteTransaction, id string, params editAttachmentParams) (*Attachment, *Error) {
	attachment := s.attachmentWithID(tx, id)
	if attachment == nil {
		return nil, ErrorUser("No script with ID")
	}

	attachment.Path = params.Path
	attachment.Owner = params.Owner
	attachment.Mode = params.Mode
	attachment.Modified = time.Now()
	attachment.AfterScript = params.AfterScript

	if params.Data != nil {
		f, err := os.OpenFile(attachment.AtomicFilePath(), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.PError("Error writing updated attachment file", map[string]interface{}{
				"attachment": attachment.ID,
				"file_path":  attachment.AtomicFilePath(),
				"error":      err.Error(),
			})
			return nil, ErrorFrom(err)
		}

		if _, err := io.Copy(f, params.Data); err != nil {
			log.PError("Error writing updated attachment file", map[string]interface{}{
				"attachment": attachment.ID,
				"file_path":  attachment.AtomicFilePath(),
				"error":      err.Error(),
			})
			f.Close()
			return nil, ErrorFrom(err)
		}
		f.Close()

		checksum, err := getFileSHA256Checksum(attachment.AtomicFilePath())
		if err != nil {
			log.PError("Error calculating attachment checksum", map[string]interface{}{
				"attachment": attachment.ID,
				"file_path":  attachment.AtomicFilePath(),
				"error":      err.Error(),
			})
			return nil, ErrorFrom(err)
		}

		if err := os.Rename(attachment.AtomicFilePath(), attachment.FilePath()); err != nil {
			log.PError("Error writing updated attachment file", map[string]interface{}{
				"attachment": attachment.ID,
				"file_path":  attachment.AtomicFilePath(),
				"error":      err.Error(),
			})
			return nil, ErrorFrom(err)
		}

		attachment.Name = params.Name
		attachment.Size = params.Size
		attachment.Checksum = checksum
		attachment.MimeType = params.MimeType
	}

	if err := tx.Update(*attachment); err != nil {
		log.PError("Error updating attachment", map[string]interface{}{
			"attachment": attachment.ID,
			"error":      err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	return attachment, nil
}

func (s attachmentStoreObject) DeleteAttachment(id string) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.deleteAttachment(tx, id)
		return nil
	})
	return
}

func (s attachmentStoreObject) deleteAttachment(tx ds.IReadWriteTransaction, id string) *Error {
	attachment := s.attachmentWithID(tx, id)
	if attachment == nil {
		return ErrorUser("No script with ID")
	}

	attachmentPath := attachment.FilePath()

	if err := tx.Delete(*attachment); err != nil {
		log.Error("Error deleting script attachment '%s': %s", attachment.ID, err.Error())
		return ErrorFrom(err)
	}

	if err := os.Remove(attachmentPath); err != nil {
		log.Error("Error deleing attachment file: attachment_id='%s' file_path='%s' error='%s'", id, attachmentPath, err.Error())
		return ErrorFrom(err)
	}

	log.Info("Deleted attachment: attachment_id='%s' file_path='%s'", id, attachmentPath)
	return nil
}

func (s attachmentStoreObject) Cleanup() (rerr *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		attachmentsWithScripts := map[string]bool{}
		attachments := s.allAttachments(tx)
		scripts := ScriptStore.AllScripts()

		for i := len(attachments) - 1; i >= 0; i-- {
			attachment := attachments[i]

			if !FileExists(attachment.FilePath()) {
				log.PWarn("Dead attachment found", map[string]interface{}{
					"attachment_id": attachment.ID,
				})
				s.deleteAttachment(tx, attachment.ID)
				attachments = append(attachments[:i], attachments[i+1:]...)
				continue
			}

			checksum, err := getFileSHA256Checksum(attachment.FilePath())
			if err != nil {
				log.PError("Error calculating attachment checksum", map[string]interface{}{
					"attachment_id": attachment.ID,
					"error":         err.Error(),
				})
				continue
			}

			if checksum != attachment.Checksum {
				log.PError("Attachment checksum verification failed", map[string]interface{}{
					"attachment_id":     attachment.ID,
					"expected_checksum": attachment.Checksum,
					"actual_checksum":   checksum,
				})
				s.deleteAttachment(tx, attachment.ID)
				attachments = append(attachments[:i], attachments[i+1:]...)
				continue
			}
		}

		attachmentIDMap := map[string]bool{}
		for _, attachment := range attachments {
			attachmentIDMap[attachment.ID] = true
		}

		for _, attachment := range attachments {
			for _, script := range scripts {
				if stringSliceContains(attachment.ID, script.AttachmentIDs) {
					attachmentsWithScripts[attachment.ID] = true
					break
				}
			}
		}

		for _, script := range scripts {
			for idx, attachmentID := range script.AttachmentIDs {
				if attachmentIDMap[attachmentID] {
					continue
				}

				log.Warn("Unknown attachment found on script: attachment_id='%s' script_id='%s'", attachmentID, script.ID)
				attachmentIDs := append(script.AttachmentIDs[:idx], script.AttachmentIDs[idx+1:]...)
				ScriptStore.EditScript(&script, editScriptParameters{
					Name:             script.Name,
					Executable:       script.Executable,
					Script:           script.Script,
					Environment:      script.Environment,
					RunAs:            script.RunAs,
					WorkingDirectory: script.WorkingDirectory,
					AfterExecution:   script.AfterExecution,
					AttachmentIDs:    attachmentIDs,
				})
			}
		}

		for _, attachment := range attachments {
			if attachmentsWithScripts[attachment.ID] {
				continue
			}

			if time.Since(attachment.Modified) > 1*time.Hour {
				log.Warn("Orphaned attachment found: attachment_id='%s'", attachment.ID)
				if err := s.deleteAttachment(tx, attachment.ID); err != nil {
					rerr = err
					return nil
				}
			}
		}
		return nil
	})
	return
}
