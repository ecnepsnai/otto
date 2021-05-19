package server

import (
	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/otto/server/environ"
)

// Script describes an otto script
type Script struct {
	ID               string `ds:"primary"`
	Name             string `ds:"unique" min:"1" max:"140"`
	Enabled          bool   `ds:"index"`
	Executable       string `min:"1"`
	Script           string `min:"1"`
	Environment      []environ.Variable
	RunAs            RunAs
	WorkingDirectory string
	AfterExecution   string
	AttachmentIDs    []string
}

// RunAs describes the properties of which user runs a script
type RunAs struct {
	Inherit bool
	UID     uint32
	GID     uint32
}

// Groups all groups with this script enabled
func (s *Script) Groups() []Group {
	enabledGroups := []Group{}

	groups := GroupStore.AllGroups()
	for _, group := range groups {
		hasScript := false
		for _, scriptID := range group.ScriptIDs {
			if scriptID == s.ID {
				hasScript = true
				break
			}
		}
		if !hasScript {
			continue
		}
		enabledGroups = append(enabledGroups, group)
	}

	return enabledGroups
}

// ScriptEnabledHost describes a host where a script is eanbled on it by a group
type ScriptEnabledHost struct {
	ScriptID   string
	ScriptName string
	GroupID    string
	GroupName  string
	HostID     string
	HostName   string
}

// Hosts all hosts with this script enabled
func (s *Script) Hosts() []ScriptEnabledHost {
	enabledHosts := []ScriptEnabledHost{}

	for _, group := range s.Groups() {
		hosts, err := group.Hosts()
		if err != nil {
			return []ScriptEnabledHost{}
		}
		ehs := make([]ScriptEnabledHost, len(hosts))
		for i, host := range hosts {
			ehs[i] = ScriptEnabledHost{
				ScriptID:   s.ID,
				ScriptName: s.Name,
				GroupID:    group.ID,
				GroupName:  group.Name,
				HostID:     host.ID,
				HostName:   host.Name,
			}
		}
		enabledHosts = append(enabledHosts, ehs...)
	}

	return enabledHosts
}

func (s *scriptStoreObject) SetGroups(script *Script, groupIDs []string) *Error {
	groups := map[string]bool{}
	allGroups := GroupStore.AllGroups()
	for _, group := range allGroups {
		var i = -1
		for y, groupID := range groupIDs {
			if groupID == group.ID {
				i = y
				break
			}
		}
		groups[group.ID] = i != -1
	}

	for groupID, enable := range groups {
		group := GroupStore.GroupWithID(groupID)
		if group == nil {
			return ErrorUser("No group with ID %s", groupID)
		}

		var i = -1
		for y, scriptID := range group.ScriptIDs {
			if scriptID == script.ID {
				i = y
				break
			}
		}

		if i == -1 && enable {
			group.ScriptIDs = append(group.ScriptIDs, script.ID)
			log.Debug("Enabling script '%s' on group '%s'", script.Name, group.Name)
		} else if i != -1 && !enable {
			group.ScriptIDs = append(group.ScriptIDs[:i], group.ScriptIDs[i+1:]...)
			log.Debug("Disabling script '%s' on group '%s'", script.Name, group.Name)
		} else {
			continue
		}

		if err := GroupStore.Table.Update(*group); err != nil {
			log.Error("Error updating group '%s': %s", group.Name, err.Error())
			return ErrorFrom(err)
		}
	}

	return nil
}

// Attachments all files for this script
func (s *Script) Attachments() ([]Attachment, *Error) {
	if len(s.AttachmentIDs) == 0 {
		return []Attachment{}, nil
	}

	attachments := make([]Attachment, len(s.AttachmentIDs))
	for i, id := range s.AttachmentIDs {
		attachment := AttachmentStore.AttachmentWithID(id)
		if attachment == nil {
			log.Error("Unknown attachment found on script: attachment_id='%s' script_id='%s'", id, s.ID)
			log.Warn("Triggering cleanup of attachments")
			AttachmentStore.Cleanup()
			return nil, ErrorUser("A temporary error occurred, reload the page to try again")
		}
		attachments[i] = *attachment
	}

	return attachments, nil
}

func (s Script) OttoScript() (*otto.Script, *Error) {
	fileIDs, err := s.Attachments()
	if err != nil {
		return nil, err
	}
	files := make([]otto.File, len(s.AttachmentIDs))
	for i, file := range fileIDs {
		file, erro := file.OttoFile()
		if erro != nil {
			return nil, ErrorFrom(erro)
		}
		files[i] = *file
	}

	return &otto.Script{
		Name: s.Name,
		RunAs: otto.RunAs{
			UID:     s.RunAs.UID,
			GID:     s.RunAs.GID,
			Inherit: s.RunAs.Inherit,
		},
		Executable:       s.Executable,
		Data:             []byte(s.Script),
		WorkingDirectory: s.WorkingDirectory,
		Environment:      map[string]string{},
		Files:            files,
	}, nil
}
