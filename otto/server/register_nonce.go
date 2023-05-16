package server

import "sync"

type registerNonceT struct {
	lock   *sync.Mutex
	nonces []string
}

var usedRegisterNonces = &registerNonceT{&sync.Mutex{}, []string{}}

func (t *registerNonceT) PreviouslyUsed(nonce string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return sliceContains(nonce, t.nonces)
}

func (t *registerNonceT) Add(nonce string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.nonces = append(t.nonces, nonce)
}
