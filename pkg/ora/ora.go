package ora

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"strconv"
	"sync"
)

const (
	TypeLayer = iota
	TypeGroup
)

type Ora struct {
	Children      []*Item
	ChildrenUUIDs map[string]*Item
	RootGroup     *Group
	ElemRoot      *XMLElement
	mutex         sync.RWMutex
}

type XMLElement struct {
	XMLName xml.Name
	Attr    []xml.Attr   `xml:",any,attr"`
	Content []byte       `xml:",innerxml"`
	Nodes   []XMLElement `xml:",any"`
}

type Item struct {
	Project *Ora
	Elem    *XMLElement
	Type    int
	Image   image.Image
}

type Group struct {
	Item
	Children []*Item
}

func NewOra() *Ora {
	return &Ora{
		Children:      make([]*Item, 0),
		ChildrenUUIDs: make(map[string]*Item),
		mutex:         sync.RWMutex{},
	}
}

func (j *Ora) Load(reader io.ReaderAt, size int64) error {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	var stackXML *zip.File
	for _, file := range zipReader.File {
		if file.Name == "stack.xml" {
			stackXML = file
			break
		}
	}

	if stackXML == nil {
		return errors.New("stack.xml not found in ORA file")
	}

	xmlFile, err := stackXML.Open()
	if err != nil {
		return fmt.Errorf("failed to open stack.xml: %w", err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	j.ElemRoot = &XMLElement{}
	err = decoder.Decode(j.ElemRoot)
	if err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	j.RootGroup = &Group{
		Item: Item{
			Project: j,
			Elem:    &j.ElemRoot.Nodes[0],
			Type:    TypeGroup,
		},
	}

	err = j.buildTree(j.RootGroup, zipReader)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	return nil
}

func (j *Ora) buildTree(parent *Group, zipReader *zip.Reader) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(parent.Elem.Nodes))

	for _, child := range parent.Elem.Nodes {
		wg.Add(1)
		go func(child XMLElement) {
			defer wg.Done()
			var item *Item
			switch child.XMLName.Local {
			case "stack":
				group := &Group{
					Item: Item{
						Project: j,
						Elem:    &child,
						Type:    TypeGroup,
					},
				}
				err := j.buildTree(group, zipReader)
				if err != nil {
					errChan <- err
				}
				item = &group.Item
			case "layer":
				layer, err := j.createLayer(&child, zipReader)
				if err != nil {
					errChan <- err
					return
				}
				item = layer
			default:
				return
			}

			j.mutex.Lock()
			j.Children = append(j.Children, item)
			uuid := getAttrValue(child.Attr, "uuid")
			j.ChildrenUUIDs[uuid] = item
			parent.Children = append(parent.Children, item)
			j.mutex.Unlock()
		}(child)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *Ora) createLayer(elem *XMLElement, zipReader *zip.Reader) (*Item, error) {
	src := getAttrValue(elem.Attr, "src")
	var imgFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == src {
			imgFile = file
			break
		}
	}

	if imgFile == nil {
		return nil, fmt.Errorf("image file not found: %s", src)
	}

	imgReader, err := imgFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer imgReader.Close()

	img, err := png.Decode(imgReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	return &Item{
		Project: j,
		Elem:    elem,
		Type:    TypeLayer,
		Image:   img,
	}, nil

}

func (j *Ora) GetByUUID(uuid string) (*Item, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	item, ok := j.ChildrenUUIDs[uuid]
	if !ok {
		return nil, fmt.Errorf("item with UUID %s not found", uuid)
	}
	return item, nil
}

func (i *Item) Name() string {
	return getAttrValue(i.Elem.Attr, "name")
}

func (i *Item) Opacity() float64 {
	opacityStr := getAttrValue(i.Elem.Attr, "opacity")
	opacity, _ := strconv.ParseFloat(opacityStr, 64)
	return opacity
}

func (i *Item) Visible() bool {
	return getAttrValue(i.Elem.Attr, "visibility") == "visible"
}

func getAttrValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
