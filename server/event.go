package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
)

var eventLog = logtic.Connect("event")

// Event describes an Otto server event
type Event struct {
	ID      string `ds:"primary"`
	Event   string `ds:"index"`
	Time    time.Time
	Details map[string]string
}

// Save save the event
func (e Event) Save() {
	if !IsEventType(e.Event) {
		panic("Attempt to add event with unknown type")
	}

	if err := EventStore.Table.Add(e); err != nil {
		log.Error("Error saving event: %s", err.Error())
	}
	if logtic.Log.Level >= logtic.LevelInfo {
		details := []string{}
		for k, v := range e.Details {
			details = append(details, fmt.Sprintf("'%s=%s'", k, v))
		}
		eventLog.Info("%s: %s", e.Event, strings.Join(details, " "))
	}
}

func newEvent(eventType string, details map[string]string) Event {
	return Event{
		ID:      newID(),
		Event:   eventType,
		Time:    time.Now(),
		Details: details,
	}
}

func (s *eventStoreObject) UserLoggedIn(username string, remoteAddr string) {
	event := newEvent(EventTypeUserLoggedIn, map[string]string{
		"username":   username,
		"remoteAddr": remoteAddr,
	})

	event.Save()
}

func (s *eventStoreObject) UserIncorrectPassword(username string, remoteAddr string) {
	event := newEvent(EventTypeUserIncorrectPassword, map[string]string{
		"username":   username,
		"remoteAddr": remoteAddr,
	})

	event.Save()
}

func (s *eventStoreObject) UserLoggedOut(username string) {
	event := newEvent(EventTypeUserLoggedOut, map[string]string{
		"username": username,
	})

	event.Save()
}

func (s *eventStoreObject) UserAdded(newUser *User, currentUser string) {
	event := newEvent(EventTypeUserAdded, map[string]string{
		"username": newUser.Username,
		"email":    newUser.Email,
		"added_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) UserModified(modifiedUsername string, currentUser string) {
	event := newEvent(EventTypeUserModified, map[string]string{
		"username":    modifiedUsername,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) UserResetPassword(username string) {
	event := newEvent(EventTypeUserResetPassword, map[string]string{
		"username": username,
	})

	event.Save()
}

func (s *eventStoreObject) UserResetAPIKey(modifiedUsername string, currentUser string) {
	event := newEvent(EventTypeUserResetAPIKey, map[string]string{
		"username": modifiedUsername,
		"reset_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) UserDeleted(deletedUsername string, currentUser string) {
	event := newEvent(EventTypeUserDeleted, map[string]string{
		"username":   deletedUsername,
		"deleted_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) HostAdded(host *Host, currentUser string) {
	event := newEvent(EventTypeHostAdded, map[string]string{
		"host_id":  host.ID,
		"name":     host.Name,
		"address":  host.Address,
		"added_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) HostModified(host *Host, currentUser string) {
	event := newEvent(EventTypeHostModified, map[string]string{
		"host_id":     host.ID,
		"name":        host.Name,
		"address":     host.Address,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) HostDeleted(host *Host, currentUser string) {
	event := newEvent(EventTypeHostDeleted, map[string]string{
		"host_id":    host.ID,
		"name":       host.Name,
		"address":    host.Address,
		"deleted_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) HostRegisterSuccess(host *Host, request otto.RegisterRequest, matchedRule *RegisterRule) {
	event := newEvent(EventTypeHostRegisterSuccess, map[string]string{
		"host_id":              host.ID,
		"name":                 host.Name,
		"address":              host.Address,
		"group_id":             host.GroupIDs[0],
		"hostname":             request.Properties.Hostname,
		"kernel_name":          request.Properties.KernelName,
		"kernel_version":       request.Properties.KernelVersion,
		"distribution_name":    request.Properties.DistributionName,
		"distribution_version": request.Properties.DistributionVersion,
	})
	if matchedRule != nil {
		for i, clause := range matchedRule.Clauses {
			keyProperty := fmt.Sprintf("matched_rule_clause%d_property", i)
			keyPattern := fmt.Sprintf("matched_rule_clause%d_pattern", i)
			event.Details[keyProperty] = clause.Property
			event.Details[keyPattern] = clause.Pattern
		}
		event.Details["matched_rule_group_id"] = matchedRule.GroupID
		event.Details["matched_rule_id"] = matchedRule.ID
		event.Details["matched_rule_name"] = matchedRule.Name
	}

	event.Save()
}

func (s *eventStoreObject) HostRegisterIncorrectPSK(request otto.RegisterRequest) {
	event := newEvent(EventTypeHostRegisterIncorrectPSK, map[string]string{
		"address":              request.Address,
		"hostname":             request.Properties.Hostname,
		"kernel_name":          request.Properties.KernelName,
		"kernel_version":       request.Properties.KernelVersion,
		"distribution_name":    request.Properties.DistributionName,
		"distribution_version": request.Properties.DistributionVersion,
	})

	event.Save()
}

func (s *eventStoreObject) GroupAdded(group *Group, currentUser string) {
	event := newEvent(EventTypeGroupAdded, map[string]string{
		"group_id": group.ID,
		"name":     group.Name,
		"added_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) GroupModified(group *Group, currentUser string) {
	event := newEvent(EventTypeGroupModified, map[string]string{
		"group_id":    group.ID,
		"name":        group.Name,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) GroupDeleted(group *Group, currentUser string) {
	event := newEvent(EventTypeGroupDeleted, map[string]string{
		"group_id":   group.ID,
		"name":       group.Name,
		"deleted_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScheduleAdded(schedule *Schedule, currentUser string) {
	event := newEvent(EventTypeScheduleAdded, map[string]string{
		"schedule_id": schedule.ID,
		"name":        schedule.Name,
		"script_id":   schedule.ScriptID,
		"pattern":     schedule.Pattern,
		"added_by":    currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScheduleModified(schedule *Schedule, currentUser string) {
	event := newEvent(EventTypeScheduleModified, map[string]string{
		"schedule_id": schedule.ID,
		"name":        schedule.Name,
		"script_id":   schedule.ScriptID,
		"pattern":     schedule.Pattern,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScheduleDeleted(schedule *Schedule, currentUser string) {
	event := newEvent(EventTypeScheduleDeleted, map[string]string{
		"schedule_id": schedule.ID,
		"name":        schedule.Name,
		"script_id":   schedule.ScriptID,
		"pattern":     schedule.Pattern,
		"deleted_by":  currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) AttachmentAdded(attachment *Attachment, currentUser string) {
	event := newEvent(EventTypeAttachmentAdded, map[string]string{
		"attachment_id": attachment.ID,
		"name":          attachment.Name,
		"file_path":     attachment.Path,
		"mimetype":      attachment.MimeType,
		"added_by":      currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) AttachmentModified(attachment *Attachment, currentUser string) {
	event := newEvent(EventTypeAttachmentModified, map[string]string{
		"attachment_id": attachment.ID,
		"name":          attachment.Name,
		"file_path":     attachment.Path,
		"mimetype":      attachment.MimeType,
		"modified_by":   currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) AttachmentDeleted(attachmentID string, currentUser string) {
	event := newEvent(EventTypeAttachmentDeleted, map[string]string{
		"attachment_id": attachmentID,
		"deleted_by":    currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScriptAdded(script *Script, currentUser string) {
	event := newEvent(EventTypeScriptAdded, map[string]string{
		"script_id": script.ID,
		"name":      script.Name,
		"added_by":  currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScriptModified(script *Script, currentUser string) {
	event := newEvent(EventTypeScriptModified, map[string]string{
		"script_id":   script.ID,
		"name":        script.Name,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScriptDeleted(script *Script, currentUser string) {
	event := newEvent(EventTypeScriptDeleted, map[string]string{
		"script_id":  script.ID,
		"name":       script.Name,
		"deleted_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) ScriptRun(script *Script, host *Host, result *otto.ScriptResult, schedule *Schedule, currentUser string) {
	event := newEvent(EventTypeScriptRun, map[string]string{
		"script_id": script.ID,
		"host_id":   host.ID,
		"exit_code": fmt.Sprintf("%d", result.Code),
	})
	if schedule != nil {
		event.Details["schedule_id"] = schedule.ID
	} else {
		event.Details["triggered_by"] = currentUser
	}

	event.Save()
}

func (s *eventStoreObject) ServerStarted(args []string) {
	a := strings.Join(args, " ")
	event := newEvent(EventTypeServerStarted, map[string]string{
		"args": a,
	})

	event.Save()
}

func (s *eventStoreObject) ServerOptionsModified(newHash string, currentUser string) {
	event := newEvent(EventTypeServerOptionsModified, map[string]string{
		"config_hash": newHash,
		"modified_by": currentUser,
	})

	event.Save()
}

func (s *eventStoreObject) RegisterRuleAdded(rule *RegisterRule, currentUser string) {
	details := map[string]string{
		"rule_id":   rule.ID,
		"rule_name": rule.Name,
		"group_id":  rule.GroupID,
		"added_by":  currentUser,
	}

	for i, clause := range rule.Clauses {
		keyProperty := fmt.Sprintf("clause%d_property", i)
		keyPattern := fmt.Sprintf("clause%d_pattern", i)
		details[keyProperty] = clause.Property
		details[keyPattern] = clause.Pattern
	}

	event := newEvent(EventTypeRegisterRuleAdded, details)
	event.Save()
}

func (s *eventStoreObject) RegisterRuleModified(rule *RegisterRule, currentUser string) {
	details := map[string]string{
		"rule_id":     rule.ID,
		"rule_name":   rule.Name,
		"group_id":    rule.GroupID,
		"modified_by": currentUser,
	}

	for i, clause := range rule.Clauses {
		keyProperty := fmt.Sprintf("clause%d_property", i)
		keyPattern := fmt.Sprintf("clause%d_pattern", i)
		details[keyProperty] = clause.Property
		details[keyPattern] = clause.Pattern
	}

	event := newEvent(EventTypeRegisterRuleAdded, details)
	event.Save()
}

func (s *eventStoreObject) RegisterRuleDeleted(rule *RegisterRule, currentUser string) {
	details := map[string]string{
		"rule_id":    rule.ID,
		"rule_name":  rule.Name,
		"group_id":   rule.GroupID,
		"deleted_by": currentUser,
	}

	for i, clause := range rule.Clauses {
		keyProperty := fmt.Sprintf("clause%d_property", i)
		keyPattern := fmt.Sprintf("clause%d_pattern", i)
		details[keyProperty] = clause.Property
		details[keyPattern] = clause.Pattern
	}

	event := newEvent(EventTypeRegisterRuleAdded, details)
	event.Save()
}
