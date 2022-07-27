package cloudflare

import (
	"bytes"
	"github.com/missionfocus/ems/internal"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Delete(path string) {
	config, err := internal.LoadConfig()

	if err != nil {
		log.Panicln("cannot read environment variable")
	}

	// Check if Directory or File
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatalf("could not os.STAT the name: %s\n", path)
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			log.Printf("err: %x", err)
		}
		for _, fi := range files {
			time.Sleep(1 * time.Second) // give cloudflare a second

			if !fi.IsDir() && strings.HasSuffix(fi.Name(), "jpg") {
				log.Printf("Deleting %s\n", filepath.Join(path, fi.Name()))
				DeleteCloudflare(filepath.Join(path, fi.Name()), &config)
			}
		}
	} else {
		DeleteCloudflare(path, &config)
	}
}

func DeleteCloudflare(path string, config *internal.Config) {
	dir, filename := filepath.Split(path)
	extension := filepath.Ext(filename)
	cloudflareId := strings.TrimSuffix(filename, extension)

	req, err := http.NewRequest(
		"DELETE",
		"https://api.cloudflare.com/client/v4/accounts/"+config.CfAccountID+"/images/v1/"+cloudflareId,
		nil,
	)
	if err != nil {
		log.Fatalf("Could not create http.NewRequest object\n")
	}

	bearer := "Bearer " + config.CfImagesStreamReadWrite
	req.Header.Set("Authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Could not open http.DefaultClient.Do\n")
	} else {
		defer resp.Body.Close()

		bodyJson := &bytes.Buffer{}
		_, err := bodyJson.ReadFrom(resp.Body)
		if err != nil {
			log.Fatalf("Could not get response from Cloudflare\n")
		}
		// Write out the response
		fileResult, err := os.Create(filepath.Join(dir, cloudflareId+"-deletedId.json"))
		if err != nil {
			log.Fatalf("Could not create file %s\n%v\n", filepath.Join(dir, cloudflareId+"-cloudflareId.json"), err.Error())

		}

		fileResult.ReadFrom(bodyJson)
		fileResult.Close()

	}
}
