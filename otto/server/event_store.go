package server

import "github.com/ecnepsnai/ds"

func (s *eventStoreObject) LastEvents(limit int) (events []Event, rerr *Error) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		objects, err := tx.GetAll(&ds.GetOptions{
			Sorted:    true,
			Ascending: true,
			Max:       limit,
		})
		if err != nil {
			rerr = ErrorFrom(err)
			return nil
		}
		count := len(objects)
		if count == 0 {
			return nil
		}
		events = make([]Event, count)
		for i, object := range objects {
			event, ok := object.(Event)
			if !ok {
				log.Panic("Incorrect object type, not Event: %#v", object)
			}
			events[i] = event
		}
		return nil
	})
	return
}
