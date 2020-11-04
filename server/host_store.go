package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *hostStoreObject) HostWithID(id string) (*Host, *Error) {
	obj, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting host with ID '%s': %s", id, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	host, k := obj.(Host)
	if !k {
		log.Error("Object is not of type 'Host'")
		return nil, ErrorServer("incorrect type")
	}

	return &host, nil
}

func (s *hostStoreObject) HostWithAddress(address string) (*Host, *Error) {
	obj, err := s.Table.GetUnique("Address", address)
	if err != nil {
		log.Error("Error getting host with address '%s': %s", address, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	host, k := obj.(Host)
	if !k {
		log.Error("Object is not of type 'Host'")
		return nil, ErrorServer("incorrect type")
	}

	return &host, nil
}

func (s *hostStoreObject) HostWithName(name string) (*Host, *Error) {
	obj, err := s.Table.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting host with name '%s': %s", name, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	host, k := obj.(Host)
	if !k {
		log.Error("Object is not of type 'Host'")
		return nil, ErrorServer("incorrect type")
	}

	return &host, nil
}

func (s *hostStoreObject) findDuplicate(name, address string) string {
	nameHost, err := s.HostWithName(name)
	if err != nil {
		return ""
	}
	if nameHost != nil {
		return nameHost.ID
	}
	addressHost, err := s.HostWithAddress(address)
	if err != nil {
		return ""
	}
	if addressHost != nil {
		return addressHost.ID
	}

	return ""
}

func (s *hostStoreObject) AllHosts() ([]Host, *Error) {
	objs, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all hosts: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []Host{}, nil
	}

	hosts := make([]Host, len(objs))
	for i, obj := range objs {
		host, k := obj.(Host)
		if !k {
			log.Error("Object is not of type 'Host'")
			return []Host{}, ErrorServer("incorrect type")
		}
		hosts[i] = host
	}

	return hosts, nil
}

type newHostParameters struct {
	Name        string
	Address     string
	Port        uint32
	PSK         string
	GroupIDs    []string
	Environment []environ.Variable
}

func (s *hostStoreObject) NewHost(params newHostParameters) (*Host, *Error) {
	if s.findDuplicate(params.Name, params.Address) != "" {
		log.Warn("Host with name '%s' or address '%s' already exists", params.Name, params.Address)
		return nil, ErrorUser("Name or Address already in use")
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, group := range params.GroupIDs {
		s, err := GroupStore.GroupWithID(group)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.Warn("No group with ID '%s'", group)
			return nil, ErrorUser("No group with ID '%s'", group)
		}
		groupIDs[i] = s.ID
	}

	host := Host{
		ID:          newID(),
		Name:        params.Name,
		Address:     params.Address,
		Port:        params.Port,
		PSK:         params.PSK,
		Enabled:     true,
		GroupIDs:    groupIDs,
		Environment: params.Environment,
	}

	if err := s.Table.Add(host); err != nil {
		log.Error("Error adding new host '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new host '%s'", params.Name)
	UpdateGroupCache()
	return &host, nil
}

type editHostParameters struct {
	Name        string
	Address     string
	Port        uint32
	PSK         string
	Enabled     bool
	GroupIDs    []string
	Environment []environ.Variable
}

func (s *hostStoreObject) EditHost(host *Host, params editHostParameters) (*Host, *Error) {
	dupID := s.findDuplicate(params.Name, params.Address)
	if dupID != "" && dupID != host.ID {
		log.Warn("Host with name '%s' or address '%s' already exists", params.Name, params.Address)
		return nil, ErrorUser("Name or Address already in use")
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, group := range params.GroupIDs {
		s, err := GroupStore.GroupWithID(group)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.Warn("No group with ID '%s'", group)
			return nil, ErrorUser("No group with ID '%s'", group)
		}
		groupIDs[i] = s.ID
	}

	host.Name = params.Name
	host.Address = params.Address
	host.Port = params.Port
	host.PSK = params.PSK
	host.Enabled = params.Enabled
	host.GroupIDs = groupIDs
	host.Environment = params.Environment

	if err := s.Table.Update(*host); err != nil {
		log.Error("Error updating host '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating host '%s'", params.Name)
	UpdateGroupCache()
	return host, nil
}

func (s *hostStoreObject) DeleteHost(host *Host) *Error {
	if err := s.Table.Delete(*host); err != nil {
		log.Error("Error deleting host '%s': %s", host.Name, err.Error())
		return ErrorFrom(err)
	}

	heartbeatStore.CleanupHeartbeats()
	UpdateGroupCache()
	log.Info("Deleting host '%s'", host.Name)
	return nil
}
