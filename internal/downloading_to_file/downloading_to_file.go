package downloading_to_file

import(
	"fmt"
	"time"
	"log"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	// "github.com/kokdot/go-musthave-devops/internal/onboarding_server"
)
var (
	// m  store.Repo = interface_init.GetM()
	// storeInterval time.Duration = onboarding_server.GetStoreInterval()
	// storeFile = onboarding_server.GetStoreFile()
	// restore = onboarding_server.GetRestore()
)

func DownloadingToFile(m repo.Repo) {
	// fmt.Println("---------DownloadingToFile: -------------------   ", storeFile)
	// fmt.Println("---------StoreInterval: -------------------   ", storeInterval)
	// fmt.Println("---------DownloadingToFile-------------------", storeInterval)

	go func() {
		// var interval = StoreInterval
		// var interval = time.Duration(storeInterval) * time.Second
		for {
			<-time.After(m.GetStoreInterval())
			fmt.Println("main; line: 67; DownloadToFile", ";  file:  ")
			err := m.WriteStorage()
			if err != nil {
				log.Printf("StoreMap did not been saved in file, err: %s", err)
			}
		}
	}()
}