package server

func staticEnvironment() map[string]string {
	return map[string]string{
		"OTTO_VERSION": ServerVersion,
		"OTTO_URL":     Options.ServerURL,
	}
}
