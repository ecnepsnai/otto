package server

import (
	"sync"
	"time"

	"github.com/ecnepsnai/otto"
)

// Heartbeat describes a heartbeat to a host
type Heartbeat struct {
	Address     string
	IsReachable bool
	LastReply   time.Time
	LastAttempt time.Time
	LastVersion string
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
	for true {
		HostStore.PingAll()
		time.Sleep(time.Duration(Options.Network.HeartbeatFrequency) * time.Minute)
	}
}

func (s *hostStoreObject) PingAll() error {
	hosts, err := s.AllHosts()
	if err != nil {
		return err.Error
	}
	for _, h := range hosts {
		go func(host Host) {
			host.Ping()
		}(h)
	}
	return nil
}

func (s *heartbeatStoreType) MarkHostReachable(host *Host, reply *otto.Reply) (*Heartbeat, *Error) {
	heartbeat := Heartbeat{
		Address:     host.Address,
		IsReachable: true,
		LastReply:   time.Now(),
		LastAttempt: time.Now(),
		LastVersion: reply.Version,
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Heartbeats[host.Address] = heartbeat
	return &heartbeat, nil
}

func (s *heartbeatStoreType) MarkHostUnreachable(host *Host) (*Heartbeat, *Error) {
	heartbeat := Heartbeat{
		Address:     host.Address,
		IsReachable: false,
		LastAttempt: time.Now(),
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Heartbeats[host.Address] = heartbeat
	return &heartbeat, nil
}

func (s *heartbeatStoreType) CleanupHeartbeats() *Error {
	heartbeats := s.AllHeartbeats()
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for _, heartbeat := range heartbeats {
		host, err := HostStore.HostWithAddress(heartbeat.Address)
		if err != nil {
			return err
		}
		if host == nil {
			delete(s.Heartbeats, heartbeat.Address)
		}
	}
	return nil
}
