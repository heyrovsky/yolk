package main

import (
	"fmt"
	"log"
	"time"

	"github.com/heyrovsky/yolk/pkg/imaged"
)

// func main() {
// 	d := downloader.NewDownloader("https://sourceforge.net/projects/osboxes/files/v/vm/59-Uu--svr/14.04.4/1404464.7z/download", "ubuntu.7z")

// 	done := make(chan error)
// 	go func() {
// 		done <- d.Download()
// 	}()

// 	for {
// 		select {
// 		case err := <-done:
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Printf("\nDownload completed: %.2f%%\n", d.DownloadPercentage)
// 			return
// 		default:
// 			fmt.Printf("\rProgress: %.2f%%", d.DownloadPercentage)
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}
// }

func main() {
	d := imaged.NewQcow2ImageDaemon(
		"ubuntu_server",
		"25.04",
		"https://sourceforge.net/projects/osboxes/files/v/vm/59-Uu--svr/25.04/64bit.7z/download",
	)

	done := make(chan error)
	go func() {
		done <- d.Exec()
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
			status, err := d.Status()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(status))
			return
		default:
			status, err := d.Status()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(status))
			time.Sleep(500 * time.Millisecond)
		}
	}

}
