package server

// Schedule describes a recurring task
type Schedule struct {
	ID       string `ds:"primary"`
	ScriptID string `ds:"index"`
	HostIDs  []string
	GroupIDs []string
}
