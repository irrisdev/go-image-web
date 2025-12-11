package store

import (
	"errors"
	"fmt"
	"go-image-web/models"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

var (
	VarientImageDir  string = "data/img/varient"
	OriginalImageDir string = "data/img/original"
)

// in memory storage for image metadata
var (
	ImageIndex   = map[string]*models.ImageMetadata{}
	ImageIndexMu sync.RWMutex
)

func init() {
	checkCreateDir(VarientImageDir)
	checkCreateDir(OriginalImageDir)
	// load images from filesystem
	loadImages()
}

// add new metadata to index
func AddImageMetadata(meta *models.ImageMetadata) {
	ImageIndexMu.Lock()
	defer ImageIndexMu.Unlock()
	ImageIndex[meta.UUID] = meta
}

// retrieve all metadata
func GetImageMetadata() map[string]*models.ImageMetadata {
	ImageIndexMu.RLock()
	defer ImageIndexMu.RUnlock()
	return ImageIndex
}

// retrieve specific image metadata based on
func GetGuidImageMetadata(uuid string) *models.ImageMetadata {
	ImageIndexMu.RLock()
	defer ImageIndexMu.RUnlock()
	return ImageIndex[uuid]
}

// append varient to existing image metadata
func AddVarientMetadata(uuid string, varient *models.ImageVarient) {
	ImageIndexMu.Lock()
	defer ImageIndexMu.Unlock()

	if v, ok := ImageIndex[uuid]; ok {
		v.SetVariant(varient.Width, *varient)
	}
}

func loadImages() {
	// load images from original directory
	originalImages, err := os.ReadDir(OriginalImageDir)
	if err != nil {
		log.Fatal(err)
	}

	// create original image metadata on server init
	for _, i := range originalImages {
		fn := i.Name()

		parts := strings.Split(fn, "_")
		if len(parts) != 2 {
			continue
		}

		uuid := parts[0]
		ext := filepath.Ext(fn)

		meta := &models.ImageMetadata{
			UUID:        uuid,
			OriginalExt: ext,
			Original:    filepath.Join(OriginalImageDir, fn),
			Varients:    make(map[int]models.ImageVarient),
		}

		ImageIndex[uuid] = meta
	}
	log.Printf("loaded: %d original images from: %s", len(ImageIndex), OriginalImageDir)

	varientImages, err := os.ReadDir(VarientImageDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, varient := range varientImages {
		fn := varient.Name()

		// strip .jpg etc extension from filename
		base := StripExtension(fn)
		parts := strings.Split(base, "_")
		if len(parts) != 2 {
			continue
		}

		// collect the 2 parts
		uuid := parts[0]
		width, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		// set index[uuid] to the varient metadata
		if meta, ok := ImageIndex[uuid]; ok {
			meta.Varients[width] = models.ImageVarient{
				Width: width,
				Path:  filepath.Join(VarientImageDir, fn),
			}
		}
	}

	// conut varients created
	var varientCount int = 0
	for _, v := range ImageIndex {
		varientCount += len(v.Varients)
	}

	log.Printf("loaded: %d varient images from: %s", varientCount, VarientImageDir)

}

func SaveOriginalImage(buf []byte, uuid string, ext string) error {
	fn := fmt.Sprintf("%s_original.%s", uuid, ext)
	savePath := filepath.Join(OriginalImageDir, fn)
	if err := os.WriteFile(savePath, buf, 0644); err != nil {
		return err
	}

	meta := &models.ImageMetadata{
		UUID:        uuid,
		OriginalExt: ext,
		Original:    savePath,
		Varients:    make(map[int]models.ImageVarient),
	}

	AddImageMetadata(meta)

	return nil
}

func SaveVarientImage(uuid string, img image.Image, wpx int, format string) error {

	// craft filename {uuid}_{width}.{format}
	fn := fmt.Sprintf("%s_%d.%s", uuid, wpx, format)

	// resize image with imaging
	resized := imaging.Resize(img, wpx, 0, imaging.Lanczos)

	// craft os path
	savePath := filepath.Join(VarientImageDir, fn)

	// create the file
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var encFmt imaging.Format
	switch format {
	case "jpeg", "jpg":
		encFmt = imaging.JPEG
	case "png":
		encFmt = imaging.PNG
	case "gif":
		encFmt = imaging.GIF
	default:
		encFmt = imaging.PNG
	}

	// save the new image file after encoding, returns error if fail
	if err := imaging.Encode(f, resized, encFmt); err != nil {
		return err
	}

	// save metadata for varient image
	AddVarientMetadata(uuid, &models.ImageVarient{
		Width: wpx,
		Path:  savePath,
	})

	return nil
}

func checkCreateDir(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func StripExtension(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
