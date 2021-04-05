package server

import (
	"io"
	"os"
	"time"

	"github.com/ecnepsnai/limits"
)

func (s attachmentStoreObject) AllAttachments() []Attachment {
	objects, err := s.Table.GetAll(nil)
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

func (s attachmentStoreObject) AllAttachmentsForScript(scriptID string) []Attachment {
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
		attachment := s.AttachmentWithID(id)
		if attachment == nil {
			log.Error("Script references non-existant attachment: script_id='%s' attachment_id='%s'", scriptID, id)
			continue
		}
		attachments[i] = *attachment
	}

	return attachments
}

func (s attachmentStoreObject) AttachmentWithID(id string) *Attachment {
	object, err := s.Table.Get(id)
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
	Data     io.Reader
	Path     string `min:"1"`
	Name     string `min:"1"`
	MimeType string `min:"1"`
	UID      int
	GID      int
	Mode     uint32
	Size     uint64
}

func (s attachmentStoreObject) NewAttachment(params newAttachmentParameters) (*Attachment, *Error) {
	if err := limits.Check(params); err != nil {
		return nil, ErrorUser(err.Error())
	}

	attachment := Attachment{
		ID:       newPlainID(),
		Path:     params.Path,
		Name:     params.Name,
		MimeType: params.MimeType,
		UID:      params.UID,
		GID:      params.GID,
		Mode:     params.Mode,
		Created:  time.Now(),
		Modified: time.Now(),
		Size:     params.Size,
	}

	f, err := os.OpenFile(attachment.FilePath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Error("Error opening attachment file '%s': %s", attachment.FilePath(), err.Error())
		return nil, ErrorFrom(err)
	}
	defer f.Close()

	if _, err := io.Copy(f, params.Data); err != nil {
		log.Error("Error writing attachment file '%s': %s", attachment.FilePath(), err.Error())
		return nil, ErrorFrom(err)
	}

	if err := AttachmentStore.Table.Add(attachment); err != nil {
		log.Error("Error saving script attachment '%s': %s", attachment.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	return &attachment, nil
}

type editAttachmentParams struct {
	Path string
	UID  int
	GID  int
	Mode uint32
}

func (s attachmentStoreObject) EditAttachment(id string, params editAttachmentParams) (*Attachment, *Error) {
	attachment := s.AttachmentWithID(id)
	if attachment == nil {
		return nil, ErrorUser("No script with ID")
	}

	attachment.Path = params.Path
	attachment.UID = params.UID
	attachment.GID = params.GID
	attachment.Mode = params.Mode
	attachment.Modified = time.Now()

	if err := s.Table.Update(*attachment); err != nil {
		log.Error("Error updating script attachment '%s': %s", attachment.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	return attachment, nil
}

func (s attachmentStoreObject) DeleteAttachment(id string) *Error {
	attachment := s.AttachmentWithID(id)
	if attachment == nil {
		return ErrorUser("No script with ID")
	}

	attachmentPath := attachment.FilePath()

	if err := s.Table.Delete(*attachment); err != nil {
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

func (s attachmentStoreObject) Cleanup() *Error {
	attachmentsWithScripts := map[string]bool{}
	attachments := s.AllAttachments()
	scripts := ScriptStore.AllScripts()

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
				Enabled:          script.Enabled,
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
			if err := s.DeleteAttachment(attachment.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
