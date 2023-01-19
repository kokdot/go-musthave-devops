package interfaceinit

import (
	"fmt"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/repo"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"github.com/kokdot/go-musthave-devops/internal/downloadingtofile"
)

var m  repo.Repo

func InterfaceInit(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) (repo.Repo, error) {
	if dataBaseDSN != "" {
		fmt.Println("----------if dataBaseDSN != \"\" {-------")
		d, err := store.NewDBStorage(storeInterval, storeFile , restore , url , key, dataBaseDSN)
		fmt.Println("----------d, err := -------", d, "------", err)
		if restore {
			err := d.ReadStorage()
			if err != nil {
				return nil, fmt.Errorf("failed to restore DBStorage, err: %s", err)
			}
		}
		return d, nil
	} else {
		if storeInterval > 0 {
			m, err := store.NewMemStorage(storeInterval, storeFile , restore , url , key, dataBaseDSN)
			if err != nil {
				return nil, fmt.Errorf("failed to create MemStorage, err: %s", err)
			}
			if storeFile != "" {
				downloadingtofile.DownloadingToFile(m)
			}
			if restore {
				err := m.ReadStorage()
				if err != nil {
					return m, nil//, fmt.Errorf("failed to restore MemStorage, err: %s", err)
				}
			}
			return m, nil
		} else {
			f, err := store.NewFileStorage(storeInterval, storeFile , restore , url , key, dataBaseDSN)
			if err != nil {
				return nil, fmt.Errorf("failed to create FileStorage, err: %s", err)
			}
			if storeFile != "" {
				downloadingtofile.DownloadingToFile(f)
			}
			if restore {
				err = m.ReadStorage()
				if err != nil {
					return nil, fmt.Errorf("failed to restore FileStorage, err: %s", err)
				}
			}
			return f, nil
		}
	}
}

