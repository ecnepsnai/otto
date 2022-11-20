package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/otto/shared/otto"
)

func (s *hostStoreObject) HostWithID(id string) (host *Host) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		host = s.hostWithID(tx, id)
		return nil
	})
	return
}

func (s *hostStoreObject) hostWithID(tx ds.IReadTransaction, id string) *Host {
	object, err := tx.Get(id)
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

func (s *hostStoreObject) HostWithAddress(address string) (host *Host) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		host = s.hostWithAddress(tx, address)
		return nil
	})
	return
}

func (s *hostStoreObject) hostWithAddress(tx ds.IReadTransaction, address string) *Host {
	object, err := tx.GetUnique("Address", address)
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

func (s *hostStoreObject) HostWithName(name string) (host *Host) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		host = s.hostWithName(tx, name)
		return nil
	})
	return
}

func (s *hostStoreObject) hostWithName(tx ds.IReadTransaction, name string) *Host {
	object, err := tx.GetUnique("Name", name)
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

func (s *hostStoreObject) findDuplicate(tx ds.IReadTransaction, name, address string) string {
	nameHost := s.hostWithName(tx, name)
	if nameHost != nil {
		return nameHost.ID
	}
	addressHost := s.hostWithAddress(tx, address)
	if addressHost != nil {
		return addressHost.ID
	}

	return ""
}

func (s *hostStoreObject) AllHosts() (hosts []Host) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		hosts = s.allHosts(tx)
		return nil
	})
	return
}

func (s *hostStoreObject) allHosts(tx ds.IReadTransaction) []Host {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
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
			panic("invalid object found in host store at insertindex " + fmt.Sprintf("%d", i) + ": " + fmt.Sprintf("%#v", obj))
		}
		hosts[i] = host
	}

	return hosts
}

type newHostParameters struct {
	Name          string
	Address       string
	Port          uint32
	AgentIdentity string
	GroupIDs      []string
	Environment   []environ.Variable
}

func (s *hostStoreObject) NewHost(params newHostParameters) (host *Host, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		host, err = s.newHost(tx, params)
		return nil
	})
	return
}

func (s *hostStoreObject) newHost(tx ds.IReadWriteTransaction, params newHostParameters) (*Host, *Error) {
	if s.findDuplicate(tx, params.Name, params.Address) != "" {
		log.Warn("Host with name '%s' or address '%s' already exists", params.Name, params.Address)
		return nil, ErrorUser("Name or Address already in use")
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, groupID := range params.GroupIDs {
		group := GroupCache.ByID(groupID)
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
		Trust:       HostTrust{},
		Enabled:     true,
		GroupIDs:    groupIDs,
		Environment: params.Environment,
	}
	if err := limits.Check(host); err != nil {
		return nil, ErrorUser(err.Error())
	}
	if params.AgentIdentity != "" {
		host.Trust.TrustedIdentity = params.AgentIdentity
		host.Trust.LastTrustUpdate = time.Now()
	}

	serverId, err := otto.NewIdentity()
	if err != nil {
		log.PError("Error generating identity for host", map[string]interface{}{
			"host_name": params.Name,
			"error":     err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	if err := tx.Add(host); err != nil {
		log.Error("Error adding new host '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}
	IdentityStore.Set(host.ID, serverId)

	HostCache.Update(tx)
	GroupStore.Table.StartRead(func(groupTx ds.IReadTransaction) error {
		GroupCache.Update(groupTx)
		return nil
	})
	log.PInfo("Created new host", map[string]interface{}{
		"id":      host.ID,
		"name":    host.Name,
		"address": host.Address,
		"port":    host.Port,
	})
	return &host, nil
}

type editHostParameters struct {
	Name        string
	Address     string
	Port        uint32
	Enabled     bool
	GroupIDs    []string
	Environment []environ.Variable
}

func (s *hostStoreObject) EditHost(host *Host, params editHostParameters) (newHost *Host, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		newHost, err = s.editHost(tx, host, params)
		return nil
	})
	return
}

func (s *hostStoreObject) editHost(tx ds.IReadWriteTransaction, host *Host, params editHostParameters) (*Host, *Error) {
	dupID := s.findDuplicate(tx, params.Name, params.Address)
	if dupID != "" && dupID != host.ID {
		log.Warn("Host with name '%s' or address '%s' already exists", params.Name, params.Address)
		return nil, ErrorUser("Name or Address already in use")
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var groupIDs = make([]string, len(params.GroupIDs))
	for i, groupID := range params.GroupIDs {
		group := GroupCache.ByID(groupID)
		if s == nil {
			log.Warn("No group with ID '%s'", groupID)
			return nil, ErrorUser("No group with ID '%s'", groupID)
		}
		groupIDs[i] = group.ID
	}

	host.Name = params.Name
	host.Address = params.Address
	host.Port = params.Port
	host.Enabled = params.Enabled
	host.GroupIDs = groupIDs
	host.Environment = params.Environment
	if err := limits.Check(host); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Update(*host); err != nil {
		log.Error("Error updating host '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating host '%s'", params.Name)
	HostCache.Update(tx)
	GroupStore.Table.StartRead(func(groupTx ds.IReadTransaction) error {
		GroupCache.Update(groupTx)
		return nil
	})
	return host, nil
}

func (s *hostStoreObject) UpdateHostTrust(hostID string, trust HostTrust) (rerr *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		host := s.hostWithID(tx, hostID)
		if host == nil {
			rerr = ErrorServer("no host with ID %s", hostID)
			return nil
		}

		host.Trust = trust
		if err := tx.Update(*host); err != nil {
			log.Error("Error updating host '%s': %s", host.Name, err.Error())
			rerr = ErrorFrom(err)
			return nil
		}

		log.Info("Updating host trust '%s'", host.Name)
		HostCache.Update(tx)
		GroupStore.Table.StartRead(func(groupTx ds.IReadTransaction) error {
			GroupCache.Update(groupTx)
			return nil
		})
		return nil
	})
	return
}

func (s *hostStoreObject) DeleteHost(host *Host) (rerr *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		for _, schedule := range ScheduleCache.All() {
			if len(schedule.Scope.HostIDs) == 0 {
				continue
			}
			if stringSliceContainsFold(host.ID, schedule.Scope.HostIDs) {
				rerr = ErrorUser("Host belongs to schedule %s", schedule.Name)
				return nil
			}
		}

		if err := tx.Delete(*host); err != nil {
			log.Error("Error deleting host '%s': %s", host.Name, err.Error())
			rerr = ErrorFrom(err)
			return nil
		}

		heartbeatStore.CleanupHeartbeats(tx)
		HostCache.Update(tx)
		GroupStore.Table.StartRead(func(groupTx ds.IReadTransaction) error {
			GroupCache.Update(groupTx)
			return nil
		})
		log.Info("Deleted host '%s'", host.Name)
		return nil
	})
	return
}
