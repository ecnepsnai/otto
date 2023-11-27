package server

// This file is was generated automatically by Codegen v1.12.2
// Do not make changes to this file as they will be lost

import (
	"path"

	"github.com/ecnepsnai/ds"
)

type attachmentStoreObject struct{ Table *ds.Table }

// AttachmentStore the global attachment store
var AttachmentStore = attachmentStoreObject{}

func cbgenDataStoreRegisterAttachmentStore() {
	table, err := ds.Register(Attachment{}, path.Join(Directories.Data, "attachment.db"), &ds.Options{DisableSorting: true})
	if err != nil {
		log.Fatal("Error registering attachment store: %s", err.Error())
	}
	AttachmentStore.Table = table
}

type eventStoreObject struct{ Table *ds.Table }

// EventStore the global event store
var EventStore = eventStoreObject{}

func cbgenDataStoreRegisterEventStore() {
	table, err := ds.Register(Event{}, path.Join(Directories.Data, "event.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering event store: %s", err.Error())
	}
	EventStore.Table = table
}

type groupStoreObject struct{ Table *ds.Table }

// GroupStore the global group store
var GroupStore = groupStoreObject{}

func cbgenDataStoreRegisterGroupStore() {
	table, err := ds.Register(Group{}, path.Join(Directories.Data, "group.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering group store: %s", err.Error())
	}
	GroupStore.Table = table
}

type hostStoreObject struct{ Table *ds.Table }

// HostStore the global host store
var HostStore = hostStoreObject{}

func cbgenDataStoreRegisterHostStore() {
	table, err := ds.Register(Host{}, path.Join(Directories.Data, "host.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering host store: %s", err.Error())
	}
	HostStore.Table = table
}

type registerruleStoreObject struct{ Table *ds.Table }

// RegisterRuleStore the global registerrule store
var RegisterRuleStore = registerruleStoreObject{}

func cbgenDataStoreRegisterRegisterRuleStore() {
	table, err := ds.Register(RegisterRule{}, path.Join(Directories.Data, "registerrule.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering registerrule store: %s", err.Error())
	}
	RegisterRuleStore.Table = table
}

type scheduleStoreObject struct{ Table *ds.Table }

// ScheduleStore the global schedule store
var ScheduleStore = scheduleStoreObject{}

func cbgenDataStoreRegisterScheduleStore() {
	table, err := ds.Register(Schedule{}, path.Join(Directories.Data, "schedule.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering schedule store: %s", err.Error())
	}
	ScheduleStore.Table = table
}

type schedulereportStoreObject struct{ Table *ds.Table }

// ScheduleReportStore the global schedulereport store
var ScheduleReportStore = schedulereportStoreObject{}

func cbgenDataStoreRegisterScheduleReportStore() {
	table, err := ds.Register(ScheduleReport{}, path.Join(Directories.Data, "schedulereport.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering schedulereport store: %s", err.Error())
	}
	ScheduleReportStore.Table = table
}

type scriptStoreObject struct{ Table *ds.Table }

// ScriptStore the global script store
var ScriptStore = scriptStoreObject{}

func cbgenDataStoreRegisterScriptStore() {
	table, err := ds.Register(Script{}, path.Join(Directories.Data, "script.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering script store: %s", err.Error())
	}
	ScriptStore.Table = table
}

type userStoreObject struct{ Table *ds.Table }

// UserStore the global user store
var UserStore = userStoreObject{}

func cbgenDataStoreRegisterUserStore() {
	table, err := ds.Register(User{}, path.Join(Directories.Data, "user.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering user store: %s", err.Error())
	}
	UserStore.Table = table
}

// dataStoreSetup set up the data store
func dataStoreSetup() {
	cbgenDataStoreRegisterAttachmentStore()
	cbgenDataStoreRegisterEventStore()
	cbgenDataStoreRegisterGroupStore()
	cbgenDataStoreRegisterHostStore()
	cbgenDataStoreRegisterRegisterRuleStore()
	cbgenDataStoreRegisterScheduleStore()
	cbgenDataStoreRegisterScheduleReportStore()
	cbgenDataStoreRegisterScriptStore()
	cbgenDataStoreRegisterUserStore()
}

// dataStoreTeardown tear down the data store
func dataStoreTeardown() {
	if AttachmentStore.Table != nil {
		AttachmentStore.Table.Close()
	}
	if EventStore.Table != nil {
		EventStore.Table.Close()
	}
	if GroupStore.Table != nil {
		GroupStore.Table.Close()
	}
	if HostStore.Table != nil {
		HostStore.Table.Close()
	}
	if RegisterRuleStore.Table != nil {
		RegisterRuleStore.Table.Close()
	}
	if ScheduleStore.Table != nil {
		ScheduleStore.Table.Close()
	}
	if ScheduleReportStore.Table != nil {
		ScheduleReportStore.Table.Close()
	}
	if ScriptStore.Table != nil {
		ScriptStore.Table.Close()
	}
	if UserStore.Table != nil {
		UserStore.Table.Close()
	}
}
