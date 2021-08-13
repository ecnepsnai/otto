package server

import (
	"strings"

	"github.com/ecnepsnai/search"
)

// SystemSearchResult describes a search result for a system search
type SystemSearchResult struct {
	Type  string
	Label string
	URL   string
}

// SystemSearch search the otto system
func SystemSearch(q string) []SystemSearchResult {
	if q == "" {
		return []SystemSearchResult{}
	}
	query := strings.ToLower(q)

	s := search.Search{}
	for _, host := range HostCache.All() {
		s.Feed(host, "Name", "Address")
	}
	for _, group := range GroupCache.All() {
		s.Feed(group, "Name")
	}
	for _, script := range ScriptCache.All() {
		s.Feed(script, "Name")
	}
	for _, schedule := range ScheduleCache.All() {
		s.Feed(schedule, "Name")
	}
	for _, user := range UserCache.All() {
		s.Feed(user, "Username", "Email")
	}

	objects := s.Search(query)
	results := make([]SystemSearchResult, len(objects))
	for i, object := range objects {
		if host, isHost := object.(Host); isHost {
			results[i] = SystemSearchResult{
				Type:  "Host",
				Label: host.Name,
				URL:   "/hosts/host/" + host.ID,
			}
		}
		if group, isGroup := object.(Group); isGroup {
			results[i] = SystemSearchResult{
				Type:  "Group",
				Label: group.Name,
				URL:   "/groups/group/" + group.ID,
			}
		}
		if script, isScript := object.(Script); isScript {
			results[i] = SystemSearchResult{
				Type:  "Script",
				Label: script.Name,
				URL:   "/scripts/script/" + script.ID,
			}
		}
		if schedule, isSchedule := object.(Schedule); isSchedule {
			results[i] = SystemSearchResult{
				Type:  "Schedule",
				Label: schedule.Name,
				URL:   "/schedules/schedule/" + schedule.ID,
			}
		}
		if user, isUser := object.(User); isUser {
			results[i] = SystemSearchResult{
				Type:  "User",
				Label: user.Username,
				URL:   "/users/user/" + user.Username,
			}
		}
	}

	log.PInfo("System search", map[string]interface{}{
		"query":       query,
		"num_results": len(results),
	})

	return results
}
