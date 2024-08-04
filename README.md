# openraster-go

openraster-go is a Go library for reading OpenRaster (ORA) files. OpenRaster is an open file format for storing layered raster graphics[1].

## Features

- 📂 Load and parse ORA files
- 🌲 Build and navigate the layer tree
- 🖼️ Access layer images and properties
- 🔍 Retrieve items by UUID
- 💫 Zero-dependencies

## Installation

To install the openraster-go library, use the following command:

```sh
go get github.com/dmytrogajewski/openraster-go
```

## Usage

Here's a quick example to get you started:

```go
package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/dmytrogajewski/openraster-go/pkg/ora"
)

func colorEquals(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func main() {
	f, err := os.Open("assets/map.ora")

	if err != nil {
		panic(err)
	}

	stats, err := f.Stat()

	if err != nil {
		panic(err)
	}

	o := ora.NewOra()
	err = o.Load(f, stats.Size())

	if err != nil {
		panic(err)
	}

	for _, oi := range o.Children {
		bounds := oi.Image.Bounds()

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := oi.Image.At(x, y)
				if colorEquals(c, color.Black) {
					fmt.Print("█")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
		fmt.Printf("%v\n", oi.Name())
	}

}
```

## API

### Types

- **Ora**: Main struct representing the ORA file
- **XMLElement**: Struct representing an XML element
- **Item**: Represents a layer or group
- **Group**: Represents a group of layers

### Methods

#### Ora

- **NewOra() *Ora**: Creates a new `Ora` instance
- **(o *Ora) Load(reader io.ReaderAt, size int64) error**: Loads and parses an ORA file
- **(o *Ora) GetByUUID(uuid string) (*Item, error)**: Retrieves an item by its UUID

#### Item

- **(i *Item) Name() string**: Returns the name of the item
- **(i *Item) Opacity() float64**: Returns the opacity of the item
- **(i *Item) Visible() bool**: Returns the visibility status of the item

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

---

Happy coding! 🎨🖌️
