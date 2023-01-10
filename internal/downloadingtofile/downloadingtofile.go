package downloading_to_file

import(
	"fmt"
	"time"
	"log"
	"github.com/kokdot/go-musthave-devops/internal/repo"
)

func DownloadingToFile(m repo.Repo) {

	go func() {
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