package ora

import (
	"archive/zip"
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestNewOra(t *testing.T) {
	ora := NewOra()
	if ora == nil {
		t.Fatal("NewOra() returned nil")
	}
	if len(ora.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(ora.Children))
	}
	if len(ora.ChildrenUUIDs) != 0 {
		t.Errorf("Expected 0 ChildrenUUIDs, got %d", len(ora.ChildrenUUIDs))
	}
}

func TestOraLoad(t *testing.T) {
	// Create a mock ORA file
	mockORA := createMockORA()

	ora := NewOra()
	err := ora.Load(bytes.NewReader(mockORA), int64(len(mockORA)))
	if err != nil {
		t.Fatalf("Failed to load mock ORA: %v", err)
	}

	if ora.RootGroup == nil {
		t.Fatal("RootGroup is nil after loading")
	}

	if len(ora.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(ora.Children))
	}
}

func TestGetByUUID(t *testing.T) {
	mockORA := createMockORA()

	ora := NewOra()
	err := ora.Load(bytes.NewReader(mockORA), int64(len(mockORA)))
	if err != nil {
		t.Fatalf("Failed to load mock ORA: %v", err)
	}

	item, err := ora.GetByUUID("layer1")
	if err != nil {
		t.Fatalf("Failed to get item by UUID: %v", err)
	}

	if item.Name() != "Layer 1" {
		t.Errorf("Expected layer name 'Layer 1', got '%s'", item.Name())
	}

	_, err = ora.GetByUUID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent UUID, got nil")
	}
}

func TestItemProperties(t *testing.T) {
	mockORA := createMockORA()

	ora := NewOra()
	err := ora.Load(bytes.NewReader(mockORA), int64(len(mockORA)))
	if err != nil {
		t.Fatalf("Failed to load mock ORA: %v", err)
	}

	item, _ := ora.GetByUUID("layer1")

	if item.Name() != "Layer 1" {
		t.Errorf("Expected name 'Layer 1', got '%s'", item.Name())
	}

	if item.Opacity() != 0.5 {
		t.Errorf("Expected opacity 0.5, got %f", item.Opacity())
	}

	if !item.Visible() {
		t.Error("Expected layer to be visible")
	}
}

func createMockORA() []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Create stack.xml
	stackXML := `<?xml version='1.0' encoding='UTF-8'?>
<image version="0.0.1">
  <stack>
    <layer name="Layer 1" src="data/layer1.png" uuid="layer1" opacity="0.5" visibility="visible"/>
    <layer name="Layer 2" src="data/layer2.png" uuid="layer2" opacity="1.0" visibility="visible"/>
  </stack>
</image>`

	xmlFile, _ := zipWriter.Create("stack.xml")
	xmlFile.Write([]byte(stackXML))

	// Create mock PNG files
	createMockPNG(zipWriter, "data/layer1.png")
	createMockPNG(zipWriter, "data/layer2.png")

	zipWriter.Close()
	return buf.Bytes()
}

func createMockPNG(zipWriter *zip.Writer, filename string) {
	pngFile, _ := zipWriter.Create(filename)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	png.Encode(pngFile, img)
}
