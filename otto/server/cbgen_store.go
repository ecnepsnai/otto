package server

// This file is was generated automatically by Codegen v1.12.3
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
func storeSetup() {
	IdentityStore.Store = cbgenStoreNewStore("identity", "")
	ShadowStore.Store = cbgenStoreNewStore("shadow", "")
	cbgenStoreRegisterGobTypes()
}
func cbgenStoreRegisterGobTypes() {
}

// storeTeardown tears down all stores
func storeTeardown() {
	IdentityStore.Store.Close()
	ShadowStore.Store.Close()
}

func cbgenStoreNewStore(storeName string, bucketName string) *store.Store {
	s, err := store.New(Directories.Data, storeName, &store.Options{BucketName: bucketName})
	if err != nil {
		log.Fatal("Error opening %s store: %s", storeName, err.Error())
	}
	return s
}
