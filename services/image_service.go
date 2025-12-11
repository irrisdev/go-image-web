package services

import (
	"bytes"
	"fmt"
	"go-image-web/models"
	"go-image-web/store"
	"image"
	"io"
	"log"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const maxUploadBytes = 15 << 20

var imageWidths = []int{600, 800, 1200, 1600}
var allowedFormats = map[string]struct{}{"jpeg": {}, "png": {}, "gif": {}, "jpg": {}}

// returns 4 types of errors(fileSize, decoding/format, whitelisted format, save original image, save varient image)
func SaveImage(file multipart.File, filename string) (string, error) {

	// limit filesize, copy buffer for several reads
	buf, err := readUpload(file)
	if err != nil {
		return "", err
	}

	// decode image from byte buffer
	img, format, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}

	// check if format is trusted based on magic bit
	if _, ok := allowedFormats[format]; !ok {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	// generate new uuid for file
	id := uuid.New().String()

	// save original image
	if err := store.SaveOriginalImage(buf, id, format); err != nil {
		return "", err
	}

	// re encode new varients for original image
	for _, wpx := range imageWidths {
		if err := store.SaveVarientImage(id, img, wpx, format); err != nil {
			log.Printf("error while saving varient image: %s_%d.%s", id, wpx, format)
			continue
		}
	}

	return id, nil
}

func GetImage(id string) (*models.ImageVarient, error) {

	// split id into parts
	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid image id format: %s", id)
	}

	uuid := parts[0]
	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid width in id: %s", id)
	}

	// return exact match if exists
	meta := store.GetGuidImageMetadata(uuid)
	if varient, ok := meta.GetVariant(width); ok {
		return &varient, nil
	}

	var (
		closest *models.ImageVarient
		minDiff int = math.MaxInt
	)

	// find min diff value using linear search
	meta.VarientsMu.RLock()
	for w, v := range meta.Varients {
		diff := abs(w - width)
		if diff < minDiff {
			minDiff = diff
			vCopy := v
			closest = &vCopy
		}
	}
	meta.VarientsMu.RUnlock()

	if closest != nil {
		return closest, nil
	}

	return nil, fmt.Errorf("no varients found for: %s", uuid)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func readUpload(file multipart.File) ([]byte, error) {
	limited := &io.LimitedReader{R: file, N: maxUploadBytes + 1}
	buf, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read upload: %w", err)
	}
	if len(buf) > maxUploadBytes {
		return nil, fmt.Errorf("file too large: max %d bytes", maxUploadBytes)
	}
	return buf, err
}

func ImageExists(id string) bool {
	return false
}

func LoadImageFile(id string) (string, error) {
	return "", nil
}
