package photo

import (
	"crypto/sha256"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type MetaData struct {
	Name           string
	CaptureDate    time.Time
	Headline       string
	Description    string
	Sha256         string
	PerceptualHash string
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

func NewPMD(filepath string) *MetaData {
	fileBytes, err := os.Open(filepath)

	if err != nil {
		log.Printf("err: %x", err)
	}

	x, err := exif.Decode(fileBytes)

	if err != nil {
		fmt.Print("should not get an error")
	}

	var pmd = new(MetaData)

	pmd.CaptureDate, err = x.DateTime()

	if err != nil {
		fmt.Print("error reading the time")
	}

	// `Format` and `Parse` use example-based layouts. Usually
	// you'll use a constant from `time` for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time `Mon Jan 2 15:04:05 MST 2006` to show the
	// pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	// pmd.CaptureYear = pmd.CaptureDate.Format("2006")
	// pmd.CaptureYearMonth = pmd.CaptureDate.Format("2006-01")
	// pmd.CaptureYearMonthDay = pmd.CaptureDate.Format("2006-01-02")

	exifValueDescription, _ := x.Get(exif.ImageDescription)

	pmd.Description = GetCleanExifValue(exifValueDescription)

	// Need to rewind the file again to get the sha256
	_, err = fileBytes.Seek(0, io.SeekStart)
	if err != nil {
		return nil
	}

	h := sha256.New()
	if _, err := io.Copy(h, fileBytes); err != nil {
		log.Printf("err: %x", err)
	}
	pmd.Sha256 = fmt.Sprintf("%x", h.Sum(nil))

	// Perceptual hash (and yet another rewind of the file)
	_, err = fileBytes.Seek(0, io.SeekStart)
	if err != nil {
		return nil
	}
	img1, _ := jpeg.Decode(fileBytes)

	phash, _ := goimagehash.PerceptionHash(img1)
	pmd.PerceptualHash = fmt.Sprintf("%x", phash.GetHash())

	err = fileBytes.Close()
	if err != nil {
		log.Fatal(err)
	}

	return pmd
}
