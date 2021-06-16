package server

import (
	"time"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/otto/server/environ"
)

// Host describes a otto host
type Host struct {
	ID            string `ds:"primary"`
	Name          string `ds:"unique" min:"1" max:"140"`
	Address       string `ds:"unique" min:"1"`
	Port          uint32
	PSK           string `min:"1" max:"512"`
	LastPSKRotate time.Time
	Enabled       bool `ds:"index"`
	GroupIDs      []string
	Environment   []environ.Variable
}

// Groups return all groups for this host
func (h Host) Groups() ([]Group, *Error) {
	groups := make([]Group, len(h.GroupIDs))
	for i, groupID := range h.GroupIDs {
		group := GroupStore.GroupWithID(groupID)
		groups[i] = *group
	}
	return groups, nil
}

// ScriptEnabledGroup describes a host where a script is enabled on it by a group
type ScriptEnabledGroup struct {
	ScriptID   string
	ScriptName string
	GroupID    string
	GroupName  string
}

// Scripts return all scripts for this host
func (h Host) Scripts() []ScriptEnabledGroup {
	hostScripts := []ScriptEnabledGroup{}
	groups, err := h.Groups()
	if err != nil {
		return nil
	}
	for _, group := range groups {
		scripts, err := group.Scripts()
		if err != nil {
			return nil
		}
		ehabledGroups := make([]ScriptEnabledGroup, len(scripts))
		for i, script := range scripts {
			ehabledGroups[i] = ScriptEnabledGroup{
				ScriptID:   script.ID,
				ScriptName: script.Name,
				GroupID:    group.ID,
				GroupName:  group.Name,
			}
		}
		hostScripts = append(hostScripts, ehabledGroups...)
	}
	return hostScripts
}

func (h *Host) RotatePSKIfNeeded() *Error {
	days := uint(time.Since(h.LastPSKRotate).Hours() / 24)
	if days < Options.Security.RotatePSK.FrequencyDays {
		return nil
	}

	_, err := h.RotatePSKNow()
	return err
}

func (h *Host) RotatePSKNow() (string, *Error) {
	newPSK := newHostPSK()

	log.PDebug("Rotating host PSK", map[string]interface{}{
		"host_id":   h.ID,
		"host_name": h.Name,
	})

	_, err := h.TriggerAction(otto.MessageTriggerAction{
		Action: otto.ActionUpdatePSK,
		NewPSK: newPSK,
	}, nil, nil)
	if err != nil {
		log.PError("Error rotating host PSK", map[string]interface{}{
			"host_id":   h.ID,
			"host_name": h.Name,
			"error":     err.Message,
		})
		return "", err
	}

	_, err = HostStore.EditHost(h, editHostParameters{
		Name:          h.Name,
		Address:       h.Address,
		Port:          h.Port,
		PSK:           newPSK,
		LastPSKRotate: time.Now(),
		Enabled:       h.Enabled,
		GroupIDs:      h.GroupIDs,
		Environment:   h.Environment,
	})
	if err != nil {
		log.PError("Error rotating host PSK", map[string]interface{}{
			"host_id":   h.ID,
			"host_name": h.Name,
			"error":     err.Message,
		})
		return "", err
	}

	log.PInfo("Rotated host PSK", map[string]interface{}{
		"host_id":   h.ID,
		"host_name": h.Name,
	})

	return newPSK, nil
}
