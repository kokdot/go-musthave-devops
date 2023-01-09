package interface_init

import (
	"log"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/repo"
	"github.com/kokdot/go-musthave-devops/internal/store"
)

var m  repo.Repo

func InterfaceInit(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) repo.Repo {
	if storeInterval > 0 {
		m, err := store.NewMemStorage(storeInterval, storeFile , restore , url , key, dataBaseDSN)
		if err != nil {
			log.Fatalf("failed to create MemStorage, err: %s", err)
		}
		return m
	} else {
			m, err := store.NewFileStorage(storeInterval, storeFile , restore , url , key, dataBaseDSN)
			if err != nil {
				log.Fatalf("failed to create FileStorage, err: %s", err)
			}
		return m
	}
}

