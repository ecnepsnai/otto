package server

import "github.com/ecnepsnai/ds"

func (s *eventStoreObject) LastEvents(limit int) ([]Event, *Error) {
	objects, err := s.Table.GetAll(&ds.GetOptions{
		Sorted:    true,
		Ascending: true,
		Max:       limit,
	})
	if err != nil {
		return []Event{}, ErrorFrom(err)
	}
	count := len(objects)
	if count == 0 {
		return []Event{}, nil
	}
	events := make([]Event, count)
	for i, object := range objects {
		event, ok := object.(Event)
		if !ok {
			log.Panic("Incorrect object type, not Event: %#v", object)
		}
		events[i] = event
	}

	return events, nil
}
