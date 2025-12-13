package store

import (
	"go-image-web/internal/models"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	return img
}

func setupTestDirs(t *testing.T) func() {
	t.Helper()

	origOriginalDir := OriginalImageDir
	origVarientDir := VarientImageDir
	origTmpDir := TmpImageDir

	tmpDir := t.TempDir()
	OriginalImageDir = filepath.Join(tmpDir, "original")
	VarientImageDir = filepath.Join(tmpDir, "varient")
	TmpImageDir = filepath.Join(tmpDir, "tmp")

	CheckCreateDir(OriginalImageDir)
	CheckCreateDir(VarientImageDir)
	CheckCreateDir(TmpImageDir)

	// Clear ImageIndex
	ImageIndexMu.Lock()
	ImageIndex = make(map[string]*models.ImageMetadata)
	ImageIndexMu.Unlock()

	return func() {
		OriginalImageDir = origOriginalDir
		VarientImageDir = origVarientDir
		TmpImageDir = origTmpDir
	}
}

func TestSaveOriginalImage_JPEG(t *testing.T) {
	cleanup := setupTestDirs(t)
	defer cleanup()

	img := createTestImage(100, 100)
	cfg := image.Config{Width: 100, Height: 100}
	uuid := "test-uuid-jpeg"

	err := SaveOriginalImage(img, cfg, uuid, "jpeg")
	if err != nil {
		t.Fatalf("SaveOriginalImage failed: %v", err)
	}

	expectedPath := filepath.Join(OriginalImageDir, uuid+"_original.jpeg")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}

	meta := GetGuidImageMetadata(uuid)
	if meta == nil {
		t.Fatal("Metadata not added to ImageIndex")
	}
	if meta.OriginalWidth != 100 || meta.OriginalHeight != 100 {
		t.Errorf("Unexpected dimensions: %dx%d", meta.OriginalWidth, meta.OriginalHeight)
	}
}

func TestSaveOriginalImage_PNG(t *testing.T) {
	cleanup := setupTestDirs(t)
	defer cleanup()

	img := createTestImage(50, 50)
	cfg := image.Config{Width: 50, Height: 50}
	uuid := "test-uuid-png"

	err := SaveOriginalImage(img, cfg, uuid, "png")
	if err != nil {
		t.Fatalf("SaveOriginalImage failed: %v", err)
	}

	expectedPath := filepath.Join(OriginalImageDir, uuid+"_original.png")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}

	meta := GetGuidImageMetadata(uuid)
	if meta == nil {
		t.Fatal("Metadata not added to ImageIndex")
	}
	if meta.OriginalExt != "png" {
		t.Errorf("Expected extension 'png', got '%s'", meta.OriginalExt)
	}
}

func TestSaveOriginalImage_GIF(t *testing.T) {
	cleanup := setupTestDirs(t)
	defer cleanup()

	img := createTestImage(30, 30)
	cfg := image.Config{Width: 30, Height: 30}
	uuid := "test-uuid-gif"

	err := SaveOriginalImage(img, cfg, uuid, "gif")
	if err != nil {
		t.Fatalf("SaveOriginalImage failed: %v", err)
	}

	expectedPath := filepath.Join(OriginalImageDir, uuid+"_original.gif")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}
}

func TestSaveOriginalImage_DefaultFormat(t *testing.T) {
	cleanup := setupTestDirs(t)
	defer cleanup()

	img := createTestImage(20, 20)
	cfg := image.Config{Width: 20, Height: 20}
	uuid := "test-uuid-default"

	err := SaveOriginalImage(img, cfg, uuid, "unknown")
	if err != nil {
		t.Fatalf("SaveOriginalImage failed: %v", err)
	}

	expectedPath := filepath.Join(OriginalImageDir, uuid+"_original.unknown")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}
}

func TestSaveOriginalImage_MetadataFields(t *testing.T) {
	cleanup := setupTestDirs(t)
	defer cleanup()

	img := createTestImage(200, 150)
	cfg := image.Config{Width: 200, Height: 150}
	uuid := "test-uuid-meta"

	err := SaveOriginalImage(img, cfg, uuid, "jpeg")
	if err != nil {
		t.Fatalf("SaveOriginalImage failed: %v", err)
	}

	meta := GetGuidImageMetadata(uuid)
	if meta == nil {
		t.Fatal("Metadata not found")
	}

	if meta.UUID != uuid {
		t.Errorf("Expected UUID '%s', got '%s'", uuid, meta.UUID)
	}
	if meta.OriginalWidth != 200 {
		t.Errorf("Expected width 200, got %d", meta.OriginalWidth)
	}
	if meta.OriginalHeight != 150 {
		t.Errorf("Expected height 150, got %d", meta.OriginalHeight)
	}
	if meta.OriginalSize <= 0 {
		t.Errorf("Expected positive file size, got %d", meta.OriginalSize)
	}
	if meta.Varients == nil {
		t.Error("Varients map not initialized")
	}
}
