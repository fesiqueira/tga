package tga

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

type decoder struct {
	rs     io.ReadSeeker
	header Header
	image  Image
	footer Footer
}

const (
	headerLen = 18
	footerLen = 26
)

func (d *decoder) decode(r io.Reader) (image.Image, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d.rs = bytes.NewReader(buf)

	err = d.read(newSection(footerLen, -footerLen, io.SeekEnd), &d.footer)
	if err != nil {
		return nil, err
	}

	err = d.read(newSection(headerLen, 0, io.SeekStart), &d.header)
	if err != nil {
		return nil, err
	}

	d.image = Image{
		ID:       make([]byte, d.header.IDLength),
		ColorMap: []byte{},
		Data:     make([]byte, d.header.ImageBytes()),
	}

	err = d.read(newSection(int(d.header.IDLength), headerLen, io.SeekStart), d.image.ID)
	if err != nil {
		return nil, err
	}

	err = d.read(newSection(d.header.ImageBytes(), headerLen+int(d.header.IDLength), io.SeekStart), d.image.Data)
	if err != nil {
		return nil, err
	}

	switch d.header.ImageType {
	case UncompressedRGBImage:
		return image.NewRGBA(image.Rect(0, 0, 0, 0)), nil
	default:
		return nil, fmt.Errorf("image type '%d' not supported", d.header.ImageType)
	}
}

func (d *decoder) read(config sectionConfig, data any) error {
	r := bytes.NewBuffer(nil)

	_, err := d.rs.Seek(config.offset, int(config.whence))
	if err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	// ensure reader is always in the beginning of the file
	defer func() {
		if config.offset != 0 && config.whence != io.SeekStart {
			d.rs.Seek(0, io.SeekStart)
		}
	}()

	_, err = io.CopyN(r, d.rs, config.length)
	if err != nil {
		return fmt.Errorf("failed to io.CopyN bytes: %v", err)
	}

	return binary.Read(r, binary.LittleEndian, data)
}

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	return d.decode(r)
}
