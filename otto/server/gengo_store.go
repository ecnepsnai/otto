package server

// This file is was generated automatically by GenGo v1.13.0
// Do not make changes to this file as they will be lost

import (
	"sync"

	"github.com/ecnepsnai/store"
)

type identityStoreObject struct {
	Store *store.Store
	Lock  *sync.Mutex
}

type shadowStoreObject struct {
	Store *store.Store
	Lock  *sync.Mutex
}

// IdentityStore the global identity store
var IdentityStore = identityStoreObject{Lock: &sync.Mutex{}}

// ShadowStore the global shadow store
var ShadowStore = shadowStoreObject{Lock: &sync.Mutex{}}

// storeSetup sets up all stores
func storeSetup(storageDir string) {
	IdentityStore.Store = gengoStoreNewStore(storageDir, "identity", "")
	ShadowStore.Store = gengoStoreNewStore(storageDir, "shadow", "")
	gengoStoreRegisterGobTypes()
}
func gengoStoreRegisterGobTypes() {
}

// storeTeardown tears down all stores
func storeTeardown() {
	IdentityStore.Store.Close()
	ShadowStore.Store.Close()
}

func gengoStoreNewStore(storageDir, storeName, bucketName string) *store.Store {
	s, err := store.New(storageDir, storeName, &store.Options{BucketName: bucketName})
	if err != nil {
		log.Fatal("Error opening %s store: %s", storeName, err.Error())
	}
	return s
}
