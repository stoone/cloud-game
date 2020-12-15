package downloader

import (
	"log"
	"os"

	"github.com/cavaliercoder/grab"
)

type GrabDownloader struct {
	conf Config
}

func NewGrabDownloader(conf Config) Downloader {
	return &GrabDownloader{
		conf: conf,
	}
}

func (d *GrabDownloader) Download(dest string, urls ...string) {
	reqs := make([]*grab.Request, 0)
	for _, url := range urls {
		req, err := grab.NewRequest(dest, url)
		if err != nil {
			panic(err)
		}
		reqs = append(reqs, req)
	}

	client := grab.NewClient()
	respch := client.DoBatch(4, reqs...)

	// check each response
	for resp := range respch {
		if err := resp.Err(); err != nil {
			log.SetOutput(os.Stderr)
			log.Printf("Download failed: %v\n", err)
			log.SetOutput(os.Stdout)
			panic(err)
		}

		log.Printf("  %v\n", resp.HTTPResponse.Status)
		log.Printf("Downloaded %s to %s\n", resp.Request.URL(), resp.Filename)
	}
}
