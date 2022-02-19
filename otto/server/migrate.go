package server

import "github.com/ecnepsnai/otto"

var neededTableVersion = 11

func migrateIfNeeded() {
	currentVersion := State.GetTableVersion()

	if currentVersion == 0 {
		State.SetTableVersion(neededTableVersion + 1)
		log.Debug("Setting default table version to %d", neededTableVersion+1)
		return
	}

	if neededTableVersion-currentVersion > 1 {
		log.Fatal("Refusing to migrate datastore that is too old - follow the supported upgrade path and don't skip versions. Table version %d, required version %d", currentVersion, neededTableVersion)
	}

	i := currentVersion
	for i <= neededTableVersion {
		i++

		cbgenDataStoreRegisterHostStore()
		storeSetup()

		hosts := HostStore.AllHosts()
		for _, host := range hosts {
			if IdentityStore.Get(host.ID) != nil {
				continue
			}

			id, err := otto.NewIdentity()
			if err != nil {
				log.PFatal("Error generating identity for host", map[string]interface{}{
					"host_id": host.ID,
					"error":   err.Error(),
				})
			}

			IdentityStore.Set(host.ID, id)
			log.PInfo("Generated server identity for host", map[string]interface{}{
				"host_id":           host.ID,
				"server_public_key": id.PublicKeyString(),
			})
		}

		HostStore.Table.Close()
		HostStore.Table = nil
		storeTeardown()
	}

	State.SetTableVersion(i)
}
