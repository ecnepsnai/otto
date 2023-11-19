package server

type Runbook struct {
	ID            string `ds:"primary"`
	Name          string `ds:"unique" min:"3" max:"128"`
	GroupIDs      []string
	ScriptIDs     []string
	HaltOnFailure bool
	RunLevel      int
}
