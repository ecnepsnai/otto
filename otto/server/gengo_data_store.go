package server

// This file is was generated automatically by GenGo v1.13.0
// Do not make changes to this file as they will be lost

import (
	"path"

	"github.com/ecnepsnai/ds"
)

type attachmentStoreObject struct{ Table *ds.Table }

// AttachmentStore the global attachment store
var AttachmentStore = attachmentStoreObject{}

func gengoDataStoreRegisterAttachmentStore(storageDir string) {
	table, err := ds.Register(Attachment{}, path.Join(storageDir, "attachment.db"), &ds.Options{DisableSorting: true})
	if err != nil {
		log.Fatal("Error registering attachment store: %s", err.Error())
	}
	AttachmentStore.Table = table
}

type eventStoreObject struct{ Table *ds.Table }

// EventStore the global event store
var EventStore = eventStoreObject{}

func gengoDataStoreRegisterEventStore(storageDir string) {
	table, err := ds.Register(Event{}, path.Join(storageDir, "event.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering event store: %s", err.Error())
	}
	EventStore.Table = table
}

type groupStoreObject struct{ Table *ds.Table }

// GroupStore the global group store
var GroupStore = groupStoreObject{}

func gengoDataStoreRegisterGroupStore(storageDir string) {
	table, err := ds.Register(Group{}, path.Join(storageDir, "group.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering group store: %s", err.Error())
	}
	GroupStore.Table = table
}

type hostStoreObject struct{ Table *ds.Table }

// HostStore the global host store
var HostStore = hostStoreObject{}

func gengoDataStoreRegisterHostStore(storageDir string) {
	table, err := ds.Register(Host{}, path.Join(storageDir, "host.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering host store: %s", err.Error())
	}
	HostStore.Table = table
}

type registerruleStoreObject struct{ Table *ds.Table }

// RegisterRuleStore the global registerrule store
var RegisterRuleStore = registerruleStoreObject{}

func gengoDataStoreRegisterRegisterRuleStore(storageDir string) {
	table, err := ds.Register(RegisterRule{}, path.Join(storageDir, "registerrule.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering registerrule store: %s", err.Error())
	}
	RegisterRuleStore.Table = table
}

type runbookStoreObject struct{ Table *ds.Table }

// RunbookStore the global runbook store
var RunbookStore = runbookStoreObject{}

func gengoDataStoreRegisterRunbookStore(storageDir string) {
	table, err := ds.Register(Runbook{}, path.Join(storageDir, "runbook.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering runbook store: %s", err.Error())
	}
	RunbookStore.Table = table
}

type runbookreportStoreObject struct{ Table *ds.Table }

// RunbookReportStore the global runbookreport store
var RunbookReportStore = runbookreportStoreObject{}

func gengoDataStoreRegisterRunbookReportStore(storageDir string) {
	table, err := ds.Register(RunbookReport{}, path.Join(storageDir, "runbookreport.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering runbookreport store: %s", err.Error())
	}
	RunbookReportStore.Table = table
}

type scheduleStoreObject struct{ Table *ds.Table }

// ScheduleStore the global schedule store
var ScheduleStore = scheduleStoreObject{}

func gengoDataStoreRegisterScheduleStore(storageDir string) {
	table, err := ds.Register(Schedule{}, path.Join(storageDir, "schedule.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering schedule store: %s", err.Error())
	}
	ScheduleStore.Table = table
}

type schedulereportStoreObject struct{ Table *ds.Table }

// ScheduleReportStore the global schedulereport store
var ScheduleReportStore = schedulereportStoreObject{}

func gengoDataStoreRegisterScheduleReportStore(storageDir string) {
	table, err := ds.Register(ScheduleReport{}, path.Join(storageDir, "schedulereport.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering schedulereport store: %s", err.Error())
	}
	ScheduleReportStore.Table = table
}

type scriptStoreObject struct{ Table *ds.Table }

// ScriptStore the global script store
var ScriptStore = scriptStoreObject{}

func gengoDataStoreRegisterScriptStore(storageDir string) {
	table, err := ds.Register(Script{}, path.Join(storageDir, "script.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering script store: %s", err.Error())
	}
	ScriptStore.Table = table
}

type userStoreObject struct{ Table *ds.Table }

// UserStore the global user store
var UserStore = userStoreObject{}

func gengoDataStoreRegisterUserStore(storageDir string) {
	table, err := ds.Register(User{}, path.Join(storageDir, "user.db"), &ds.Options{})
	if err != nil {
		log.Fatal("Error registering user store: %s", err.Error())
	}
	UserStore.Table = table
}

// dataStoreSetup set up the data store
func dataStoreSetup(storageDir string) {
	gengoDataStoreRegisterAttachmentStore(storageDir)
	gengoDataStoreRegisterEventStore(storageDir)
	gengoDataStoreRegisterGroupStore(storageDir)
	gengoDataStoreRegisterHostStore(storageDir)
	gengoDataStoreRegisterRegisterRuleStore(storageDir)
	gengoDataStoreRegisterRunbookStore(storageDir)
	gengoDataStoreRegisterRunbookReportStore(storageDir)
	gengoDataStoreRegisterScheduleStore(storageDir)
	gengoDataStoreRegisterScheduleReportStore(storageDir)
	gengoDataStoreRegisterScriptStore(storageDir)
	gengoDataStoreRegisterUserStore(storageDir)
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
	if RunbookStore.Table != nil {
		RunbookStore.Table.Close()
	}
	if RunbookReportStore.Table != nil {
		RunbookReportStore.Table.Close()
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
