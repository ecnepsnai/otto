package server

import (
	"sync"
	"time"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/shared/otto"
)

// Heartbeat describes a heartbeat to a host
type Heartbeat struct {
	Address     string
	IsReachable bool
	LastReply   time.Time
	LastAttempt time.Time
	Version     string
	Properties  map[string]string
}

type heartbeatStoreType struct {
	Heartbeats map[string]Heartbeat
	Lock       *sync.Mutex
}

var heartbeatStore = &heartbeatStoreType{map[string]Heartbeat{}, &sync.Mutex{}}

func (s *heartbeatStoreType) AllHeartbeats() []Heartbeat {
	var heartbeats = make([]Heartbeat, len(s.Heartbeats))

	s.Lock.Lock()
	defer s.Lock.Unlock()
	i := 0
	for _, heartbeat := range s.Heartbeats {
		heartbeats[i] = heartbeat
		i++
	}

	return heartbeats
}

func (s *heartbeatStoreType) LastHeartbeat(host *Host) *Heartbeat {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	heartbeat, present := s.Heartbeats[host.Address]
	if !present {
		return nil
	}

	return &heartbeat
}

// StartHeartbeatMonitor starts the heartbeat monitor
func StartHeartbeatMonitor() {
	for {
		HostStore.PingAll()
		time.Sleep(time.Duration(Options.Network.HeartbeatFrequency) * time.Minute)
	}
}

func (s *hostStoreObject) PingAll() error {
	hosts := s.AllHosts()
	for _, h := range hosts {
		go func(host Host) {
			host.Ping()
			host.RotateIdentityIfNeeded()
		}(h)
	}
	return nil
}

func (s *heartbeatStoreType) RegisterHeartbeatReply(host *Host, reply otto.MessageHeartbeatResponse) (*Heartbeat, *Error) {
	heartbeat := Heartbeat{
		Address:     host.Address,
		IsReachable: true,
		LastReply:   time.Now(),
		LastAttempt: time.Now(),
		Version:     reply.AgentVersion,
		Properties:  reply.Properties,
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()

	wasUnreachable := true
	if hb, p := s.Heartbeats[host.Address]; p {
		if hb.IsReachable {
			wasUnreachable = false
		}
	}

	s.Heartbeats[host.Address] = heartbeat

	if wasUnreachable {
		log.PInfo("Host became reachable", map[string]interface{}{
			"host_id":   host.ID,
			"host_name": host.Name,
		})
		EventStore.HostBecameReachable(host)
	}

	return &heartbeat, nil
}

func (s *heartbeatStoreType) UpdateHostReachability(host *Host, isReachable bool) (*Heartbeat, *Error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	log.PDebug("Update reachability", map[string]interface{}{
		"host_id":      host.ID,
		"host_name":    host.Name,
		"is_reachable": isReachable,
	})

	heartbeat := s.Heartbeats[host.Address]

	becameUnreachable := heartbeat.IsReachable && !isReachable
	becameReachable := !heartbeat.IsReachable && isReachable

	heartbeat.Address = host.Address
	heartbeat.IsReachable = isReachable
	heartbeat.LastAttempt = time.Now()

	s.Heartbeats[host.Address] = heartbeat

	if becameUnreachable {
		log.PWarn("Host became unreachable", map[string]interface{}{
			"host_id":   host.ID,
			"host_name": host.Name,
		})
		EventStore.HostBecameUnreachable(host, heartbeat.LastReply)
	}
	if becameReachable {
		log.PInfo("Host became reachable", map[string]interface{}{
			"host_id":   host.ID,
			"host_name": host.Name,
		})
		EventStore.HostBecameReachable(host)
	}

	s.UpdateReachabilityStats()
	return &heartbeat, nil
}

func (s *heartbeatStoreType) UpdateReachabilityStats() {
	reachable := uint64(0)
	unreachable := uint64(0)
	for _, hb := range s.Heartbeats {
		if hb.IsReachable {
			reachable++
		} else {
			unreachable++
		}
	}
	Stats.Counters.ReachableHosts.Set(reachable)
	Stats.Counters.UnreachableHosts.Set(unreachable)
}

func (s *heartbeatStoreType) CleanupHeartbeats(hostTx ds.IReadTransaction) *Error {
	heartbeats := s.AllHeartbeats()
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for _, heartbeat := range heartbeats {
		host := HostStore.hostWithAddress(hostTx, heartbeat.Address)
		if host == nil {
			delete(s.Heartbeats, heartbeat.Address)
		}
	}
	s.UpdateReachabilityStats()
	return nil
}
