package interface_init

import (
	"fmt"
	"log"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/repo"
	"github.com/kokdot/go-musthave-devops/internal/store"
)

var (
	m  repo.Repo
	// storeInterval time.Duration = onboarding_server.GetStoreInterval()
	// storeFile = onboarding_server.GetStoreFile()
	// restore = onboarding_server.GetRestore()
	// url = onboarding_server.GetURL()
	// key = onboarding_server.GetKey()
)

func InterfaceInit(storeInterval time.Duration, storeFile string, restore bool, url string, key string) repo.Repo {
	// fmt.Println("------------------------handler---------InterfaceInit---------------------------------")
	


	if storeInterval > 0 {
		fmt.Println("-------------------------------------if storeInterval > 0--------------------------------------")
		// if storeFile == "" {
		fmt.Println("-------------------------------------if storeFile == 0--------------------------------------")

			m, err := store.NewMemStorage(storeInterval, storeFile , restore , url , key)
            if err != nil {
				log.Fatalf("failed to create MemStorage, err: %s", err)
			}
			// fmt.Printf("!!!!!!!------------1----------!!!!!!!m:    %#v\n", m)
			return m
		// } else {
		// 	fmt.Println("------------------------------------else-if storeFile == 0--------------------------------------")

		// 	ms, err := store.NewMemStorageWithFile(storeFile)
		// 	if err != nil {
		// 		log.Fatalf("failed to create MemStorage, err: %s", err)
		// 	}
		// 	m = ms
		// 	fmt.Printf("!!!!!!!-----------2-----------!!!!!!!m:    %#v\n", m)
		// 	if restore {
		// 		fmt.Println("-------------------------------------if restore--------------------------------------")

		// 		_, err := m.ReadStorage()
		// 		if err != nil {
		// 			log.Printf("Can't to read data froM file, err: %s", err)
		// 		}
		// 		sm, _ := m.GetAllValues()
		// 		fmt.Println("---------------------------- sm -----------------!restore!----------------:    \n", sm)
		// 	}
		// 	fmt.Println("-------------------------------------DownloadingToFile--------------------------------------")
		// 	downloading_to_file.DownloadingToFile(m, )
		// 	// m = ms
		// 	fmt.Printf("!!!!!!!----------3------------!!!!!!!m:    %#v\n", m)
		// 	return m
		// }
	} else {
		fmt.Println("-----------------------------------------else----if storeInterval > 0------------------------------------")
		// if storeFile == "" {
			fmt.Println("---------------------------------------------if storeFile == 0------------------------------------")
			m, err := store.NewFileStorage(storeInterval, storeFile , restore , url , key)
			if err != nil {
				log.Fatalf("failed to create FileStorage, err: %s", err)
			}
		// 	m = ms
		// } else {
		// 	fmt.Println("------------------------------------------else---if storeFile == 0------------------------------------")

		// 	ms, err := store.NewFileStorageWithFile(storeFile)
		// 	if err != nil {
		// 		log.Fatalf("failed to create FileStorage, err: %s", err)
		// 	}
		// 	m = ms
		// 	if restore {
		// 		fmt.Println("-------------------------------------if restore--------------------------------------")

		// 		_, err := m.ReadStorage()
		// 		if err != nil {
		// 			log.Printf("Can't to read data from file, err: %s", err)
		// 		}
		// 	}
		// }
		return m
	}
}

