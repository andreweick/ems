package cloudflare

//noinspection GoUnhandledErrorResult

import (
	"bytes"
	"github.com/missionfocus/ems/internal"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Add(path string) {
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
				AddFile(path+fi.Name(), &config)
			}
		}
	} else {
		AddFile(path, &config)
	}
}

func AddFile(path string, config *internal.Config) {
	dir, filename := filepath.Split(path)
	extension := filepath.Ext(filename)
	cloudflareId := strings.TrimSuffix(filename, extension)

	// Write to Cloudflare
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open file: %s\n", path)
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		log.Fatalf("Could not create multipart FormFile writer")
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("Could not copy the file into the part of multipart\n")
	}

	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Disposition", "form-data; name=\"id\"")

	part2, err := writer.CreatePart(metadataHeader)
	if err != nil {
		log.Fatalf("Could not create part2 of multipart\n")
	}

	part2.Write([]byte(cloudflareId))

	err = writer.Close()
	if err != nil {
		log.Fatalf("Could not close writer\n")
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.cloudflare.com/client/v4/accounts/"+config.CfAccountID+"/images/v1",
		body,
	)
	if err != nil {
		log.Fatalf("Could not create http.NewRequest object\n")
	}

	bearer := "Bearer " + config.CfImagesStreamReadWrite

	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		fileResult, err := os.Create(filepath.Join(dir, cloudflareId+"-cloudflareId.json"))
		if err != nil {
			log.Fatalf("Could not create file %s\n%v\n", filepath.Join(dir, cloudflareId+"-cloudflareId.json"), err.Error())

		}

		fileResult.ReadFrom(bodyJson)
		fileResult.Close()
	}

	// Write photograph metadata
	
}
