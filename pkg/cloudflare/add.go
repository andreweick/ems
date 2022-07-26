package cloudflare

//noinspection GoUnhandledErrorResult

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/missionfocus/ems/internal"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"image/jpeg"
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

type MetaData struct {
	Id             string    `json:"id"`
	CaptureDate    time.Time `json:"CaptureDate"`
	Headline       string    `json:"Headline"`
	Description    string    `json:"Description"`
	Sha256         string    `json:"Sha256"`
	PerceptualHash string    `json:"PerceptualHash"`
}

func GetCleanExifValue(md *tiff.Tag) string {
	if md == nil {
		return ""
	}
	s := fmt.Sprintf("%v", md)

	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
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
				AddCloudflare(path+fi.Name(), &config)
				AddMetadata(path + fi.Name())
			}
		}
	} else {
		AddCloudflare(path, &config)
		AddMetadata(path)
	}
}

func AddCloudflare(path string, config *internal.Config) {
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
}

func AddMetadata(path string) {
	dir, filename := filepath.Split(path)
	extension := filepath.Ext(filename)
	cloudflareId := strings.TrimSuffix(filename, extension)

	// Write photograph metadata
	fileExifBytes, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open for Exif reading: %s\n", path)
	}

	x, err := exif.Decode(fileExifBytes)
	if err != nil {
		log.Fatalf("Cannot decode Exif bytes: %s\n", path)
	}

	var pmd = new(MetaData)

	pmd.Id = cloudflareId

	pmd.CaptureDate, err = x.DateTime()
	if err != nil {
		log.Printf("Cannot read DateTime from %s\n", path)
	}

	exifValueDescription, _ := x.Get(exif.ImageDescription)
	pmd.Description = GetCleanExifValue(exifValueDescription)

	// Need to rewind the file again to get the sha256
	_, err = fileExifBytes.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalf("Cannot rewind file %s\n", path)
	}

	h := sha256.New()
	if _, err := io.Copy(h, fileExifBytes); err != nil {
		log.Fatalf("cannot get sha256 from %s\n", path)
	}
	pmd.Sha256 = fmt.Sprintf("%x", h.Sum(nil))

	// Perceptual hash (and yet another rewind of the file)
	_, err = fileExifBytes.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalf("Cannot rewind file (perceptual hash) %s\n", path)
	}
	img1, _ := jpeg.Decode(fileExifBytes)

	phash, _ := goimagehash.PerceptionHash(img1)
	pmd.PerceptualHash = fmt.Sprintf("%x", phash.GetHash())

	err = fileExifBytes.Close()
	if err != nil {
		log.Fatalf("Cannot close file after perceptual hash: %s\n", path)
	}

	fileResult, err := os.Create(filepath.Join(dir, cloudflareId+"-metadata.json"))
	if err != nil {
		log.Fatalf("Could not create file %s\n%v\n", filepath.Join(dir, cloudflareId+"-metadata.json"), err.Error())
	}

	jsonString, _ := json.MarshalIndent(pmd, "", "  ")
	fileResult.Write(jsonString)

	fileResult.Close()
}
