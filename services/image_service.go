package services

import (
	"fmt"
	"go-image-web/models"
	"go-image-web/store"
	"image"
	"io"
	"log"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var imageWidths = []int{600, 800, 1200, 1600}
var allowedFormats = map[string]struct{}{"jpeg": {}, "png": {}, "gif": {}, "jpg": {}}

const (
	MaxImageHeight = 2560
	MaxImageWidth  = 2560
)

// returns 4 types of errors(fileSize, decoding/format, whitelisted format, save original image, save varient image)
func SaveImage(file multipart.File, filename string) (string, error) {

	// generate new uuid for file
	id := uuid.New().String()

	// create temp file to process
	tmpPath, err := store.CreateTmpFile(id, file)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	// re-open tmp save
	srcFile, err := os.Open(tmpPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// cheap load of image config
	cfg, format, err := image.DecodeConfig(srcFile)
	if err != nil {
		return "", err
	}

	// check if format is trusted based on magic bit
	if _, ok := allowedFormats[format]; !ok {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	log.Printf("format: %s, %dx%d", format, cfg.Width, cfg.Height)

	// // validate image config
	// if cfg.Width > MaxImageWidth || cfg.Height > MaxImageHeight {
	// 	return "", fmt.Errorf("dimensions exeed limits of: %dx%d", MaxImageWidth, MaxImageHeight)
	// }

	// Rewind file for full decode
	if _, err := srcFile.Seek(0, 0); err != nil {
		return "", err
	}

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	// save original image
	if err := store.SaveOriginalImage(img, cfg, id, format); err != nil {
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

	originalMeta := store.GetGuidImageMetadata(id)
	if originalMeta != nil {
		return &models.ImageVarient{
			Path: originalMeta.OriginalPath,
		}, nil
	}

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
	limited := &io.LimitedReader{R: file, N: store.MaxUploadBytes + 1}
	buf, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read upload: %w", err)
	}
	if len(buf) > store.MaxUploadBytes {
		return nil, fmt.Errorf("file too large: max %d bytes", store.MaxUploadBytes)
	}
	return buf, err
}

func ImageExists(id string) bool {
	return false
}

func LoadImageFile(id string) (string, error) {
	return "", nil
}
