package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *hostStoreObject) HostWithID(id string) *Host {
	object, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting host: id='%s' error='%s'", id, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	host, k := object.(Host)
	if !k {
		log.Error("Error getting host: id='%s' error='%s'", id, "invalid type")
	}

	return &host
}

func (s *hostStoreObject) HostWithAddress(address string) *Host {
	object, err := s.Table.GetUnique("Address", address)
	if err != nil {
		log.Error("Error getting host: address='%s' error='%s'", address, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	host, k := object.(Host)
	if !k {
		log.Error("Error getting host: address='%s' error='%s'", address, "invalid type")
	}

	return &host
}

func (s *hostStoreObject) HostWithName(name string) *Host {
	object, err := s.Table.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting host: name='%s' error='%s'", name, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	host, k := object.(Host)
	if !k {
		log.Error("Error getting host: name='%s' error='%s'", name, "invalid type")
	}

	return &host
}

func (s *hostStoreObject) findDuplicate(name, address string) string {
	nameHost := s.HostWithName(name)
	if nameHost != nil {
		return nameHost.ID
	}
	addressHost := s.HostWithAddress(address)
	if addressHost != nil {
		return addressHost.ID
	}

	return ""
}

func (s *hostStoreObject) AllHosts() []Host {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all hosts: error='%s'", err.Error())
		return []Host{}
	}
	if len(objects) == 0 {
		return []Host{}
	}

	hosts := make([]Host, len(objects))
	for i, obj := range objects {
		host, k := obj.(Host)
		if !k {
			log.Fatal("Error listing all hosts: error='%s'", "invalid type")
		}
		hosts[i] = host
	}

	return hosts
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

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, groupID := range params.GroupIDs {
		group := GroupStore.GroupWithID(groupID)
		if s == nil {
			log.Warn("No group with ID '%s'", groupID)
			return nil, ErrorUser("No group with ID '%s'", groupID)
		}
		groupIDs[i] = group.ID
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
	if err := limits.Check(host); err != nil {
		return nil, ErrorUser(err.Error())
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

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, groupID := range params.GroupIDs {
		group := GroupStore.GroupWithID(groupID)
		if s == nil {
			log.Warn("No group with ID '%s'", groupID)
			return nil, ErrorUser("No group with ID '%s'", groupID)
		}
		groupIDs[i] = group.ID
	}

	host.Name = params.Name
	host.Address = params.Address
	host.Port = params.Port
	host.PSK = params.PSK
	host.Enabled = params.Enabled
	host.GroupIDs = groupIDs
	host.Environment = params.Environment
	if err := limits.Check(host); err != nil {
		return nil, ErrorUser(err.Error())
	}

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
