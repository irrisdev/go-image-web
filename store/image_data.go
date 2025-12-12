package store

import (
	"bytes"
	"errors"
	"fmt"
	"go-image-web/models"
	"image"
	"image/gif"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

const MaxUploadBytes = 15 << 20 // 10MB

var (
	VarientImageDir  string = "data/img/varient"
	OriginalImageDir string = "data/img/original"
	TmpImageDir      string = "data/img/tmp"
)

// in memory storage for image metadata
var (
	ImageIndex   = map[string]*models.ImageMetadata{}
	ImageIndexMu sync.RWMutex
)

func init() {
	CheckCreateDir(VarientImageDir)
	CheckCreateDir(OriginalImageDir)
	CheckCreateDir(TmpImageDir)
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

// Create a temp file before processing
func CreateTmpFile(uuid string, file multipart.File) (string, error) {
	path := filepath.Join(TmpImageDir, uuid)
	tmpFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		os.Remove(path)
		return "", err
	}

	return path, nil
}

func loadImages() {
	// load images from original directory
	originalImages, err := os.ReadDir(OriginalImageDir)
	if err != nil {
		log.Fatal(err)
	}

	// create original image metadata on server init
	for _, i := range originalImages {

		// check file info
		fi, err := i.Info()
		if err != nil {
			log.Println(err)
			continue
		}

		// check file naming convention
		parts := strings.Split(fi.Name(), "_")
		if len(parts) != 2 {
			log.Println("invalid filename")
			continue
		}

		// open file to decode config
		srcPath := filepath.Join(OriginalImageDir, fi.Name())
		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Print(err)
			continue
		}
		defer srcFile.Close()

		// check config and format is valid
		cfg, format, err := image.DecodeConfig(srcFile)
		if err != nil {
			log.Print(err)
			continue
		}

		uuid := parts[0]
		meta := &models.ImageMetadata{
			UUID:         uuid,
			OriginalExt:  format,
			OriginalPath: srcPath,
			ModifiedTime: fi.ModTime(),
			OriginalSize: fi.Size(),

			OriginalWidth:  cfg.Width,
			OriginalHeight: cfg.Height,

			Varients: make(map[int]models.ImageVarient),
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

		// open varient file
		srcPath := filepath.Join(VarientImageDir, fn)
		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Print(err)
			continue
		}
		defer srcFile.Close()

		// load image config for valid extension/format
		_, format, err := image.DecodeConfig(srcFile)
		if err != nil {
			log.Print(err)
			continue
		}

		// strip .jpg etc extension from filename
		base := StripExtension(fn)
		parts := strings.Split(base, "_")
		if len(parts) != 2 {
			log.Println("invalid filename")
			continue
		}

		// collect the 2 parts
		uuid := parts[0]
		width, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Println("invalid filename")
			continue
		}

		// set index[uuid] to the varient metadata
		if meta, ok := ImageIndex[uuid]; ok {
			meta.Varients[width] = models.ImageVarient{
				Width: width,
				Path:  srcPath,
				Ext:   format,
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

func SaveOriginalImage(img image.Image, cfg image.Config, uuid string, format string) error {
	// create new filename
	fn := fmt.Sprintf("%s_original.%s", uuid, format)
	// create filepath
	savePath := filepath.Join(OriginalImageDir, fn)
	// create system file
	saveFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer saveFile.Close()

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

	// encode image and remove if fail
	if err := imaging.Encode(saveFile, img, encFmt); err != nil {
		os.Remove(savePath)
		return err
	}

	// create image metadata
	meta := &models.ImageMetadata{
		UUID:           uuid,
		OriginalExt:    format,
		OriginalPath:   savePath,
		OriginalWidth:  cfg.Width,
		OriginalHeight: cfg.Height,
		ModifiedTime:   time.Now(),

		Varients: make(map[int]models.ImageVarient),
	}

	if err := saveFile.Sync(); err == nil {
		fi, err := saveFile.Stat()
		if err == nil {
			meta.OriginalSize = fi.Size()
		}
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
		os.Remove(savePath)
		return err
	}

	// save metadata for varient image
	AddVarientMetadata(uuid, &models.ImageVarient{
		Width: wpx,
		Path:  savePath,
		Ext:   format,
	})

	return nil
}

func SaveVarientGIF(uuid string, buf []byte, wpx int) error {
	// Decode the full GIF with all frames
	g, err := gif.DecodeAll(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("decode gif: %w", err)
	}

	// Resize each frame
	for i, frame := range g.Image {
		resized := imaging.Resize(frame, wpx, 0, imaging.Lanczos)
		// Convert back to paletted image
		bounds := resized.Bounds()
		palettedImg := image.NewPaletted(bounds, frame.Palette)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				palettedImg.Set(x, y, resized.At(x, y))
			}
		}
		g.Image[i] = palettedImg
	}

	// Update config dimensions
	g.Config.Width = wpx
	g.Config.Height = 0 // Will be set by aspect ratio

	fn := fmt.Sprintf("%s_%d.gif", uuid, wpx)
	savePath := filepath.Join(VarientImageDir, fn)

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, g)
}

func CheckCreateDir(path string) {
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
