package server

import (
	"bytes"
	"os"
	"testing"

	"github.com/ecnepsnai/security"
)

func TestAddGetAttachment(t *testing.T) {
	data := bytes.NewReader(security.RandomBytes(16))

	attachment, err := AttachmentStore.NewAttachment(newAttachmentParameters{
		Data:     data,
		Path:     randomString(6),
		Name:     randomString(6),
		MimeType: "text",
		UID:      os.Getuid(),
		GID:      os.Getgid(),
		Mode:     0644,
		Size:     16,
	})
	if err != nil {
		t.Fatalf("Error making new attachment: %s", err.Message)
	}
	if attachment == nil {
		t.Fatalf("Should return a attachment")
	}

	if AttachmentStore.AttachmentWithID(attachment.ID) == nil {
		t.Fatalf("Should return a attachment with an ID")
	}

	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:          randomString(6),
		Executable:    "/bin/bash",
		Script:        "#!/bin/bash\necho hello\n",
		AttachmentIDs: []string{attachment.ID},
		UID:           0,
		GID:           0,
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	if len(AttachmentStore.AllAttachmentsForScript(script.ID)) == 0 {
		t.Fatalf("Should return an attachment")
	}
}

func TestEditAttachment(t *testing.T) {
	data := bytes.NewReader(security.RandomBytes(16))

	attachment, err := AttachmentStore.NewAttachment(newAttachmentParameters{
		Data:     data,
		Path:     randomString(6),
		Name:     randomString(6),
		MimeType: "text",
		UID:      os.Getuid(),
		GID:      os.Getgid(),
		Mode:     0644,
		Size:     16,
	})
	if err != nil {
		t.Fatalf("Error making new attachment: %s", err.Message)
	}
	if attachment == nil {
		t.Fatalf("Should return a attachment")
	}

	_, err = AttachmentStore.EditAttachment(attachment.ID, editAttachmentParams{
		Path: attachment.Path,
		UID:  attachment.UID,
		GID:  attachment.GID,
		Mode: 0777,
	})
	if err != nil {
		t.Fatalf("Error editing attachment: %s", err.Message)
	}

	if AttachmentStore.AttachmentWithID(attachment.ID).Mode == attachment.Mode {
		t.Fatalf("Should update attachment mode")
	}
}

func TestDeleteAttachment(t *testing.T) {
	data := bytes.NewReader(security.RandomBytes(16))

	attachment, err := AttachmentStore.NewAttachment(newAttachmentParameters{
		Data:     data,
		Path:     randomString(6),
		Name:     randomString(6),
		MimeType: "text",
		UID:      os.Getuid(),
		GID:      os.Getgid(),
		Mode:     0644,
		Size:     16,
	})
	if err != nil {
		t.Fatalf("Error making new attachment: %s", err.Message)
	}
	if attachment == nil {
		t.Fatalf("Should return a attachment")
	}

	if err := AttachmentStore.DeleteAttachment(attachment.ID); err != nil {
		t.Fatalf("Error deleting attachment: %s", err.Message)
	}

	if FileExists(attachment.FilePath()) {
		t.Fatalf("Attachment file should not exist")
	}
}
