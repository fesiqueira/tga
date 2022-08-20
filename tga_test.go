package tga

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

type tgaFileMock struct {
	r       *bytes.Buffer
	content []byte
	offset  int64
}

func newTGAMock(header, image, footer []byte) *tgaFileMock {
	buf := make([]byte, 0)

	buf = append(buf, header...)
	buf = append(buf, image...)
	buf = append(buf, footer...)

	return &tgaFileMock{
		r:       bytes.NewBuffer(buf),
		content: buf,
	}
}

func (fm *tgaFileMock) Read(buf []byte) (int, error) {
	return fm.r.Read(buf)
}

func (fm *tgaFileMock) Seek(offset int64, whence int) (int64, error) {
	var err error

	switch whence {
	case io.SeekStart:
		if offset > 0 {
			offset--
		}

		fm.offset = offset
	case io.SeekCurrent:
		fm.offset += offset
	case io.SeekEnd:
		fm.offset = int64(len(fm.content)) + offset
	default:
		err = fmt.Errorf("unknown whence: `%v`", whence)
	}

	fm.r = bytes.NewBuffer(fm.content[fm.offset:len(fm.content)])

	return fm.offset, err
}

func TestRead(t *testing.T) {
	testCases := []struct {
		rs       io.ReadSeeker
		expected File
	}{
		{
			rs: newTGAMock([]byte{0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0}, []byte{12}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 'T', 'R', 'U', 'E', 'V', 'I', 'S', 'I', 'O', 'N', '-', 'X', 'F', 'I', 'L', 'E', '.', 0x00}),
			expected: File{
				Header: Header{
					IDLength:        0,
					ColorMapType:    0,
					ImageType:       UncompressedRGBImage,
					ColorMapOrigin:  0,
					ColorMapLength:  0,
					ColorMapDepth:   0,
					XOrigin:         0,
					YOrigin:         0,
					Width:           1,
					Height:          1,
					BitsPerPixel:    0,
					ImageDescriptor: 0,
				},
				Image: Image{
					ID:       []byte{},
					ColorMap: []byte{},
					Data:     []byte{12},
				},
				Footer: Footer{
					ExtensionAreaOffset:      0,
					DeveloperDirectoryOffset: 0,
					Signature:                [16]byte{'T', 'R', 'U', 'E', 'V', 'I', 'S', 'I', 'O', 'N', '-', 'X', 'F', 'I', 'L', 'E'},
					Point:                    '.',
					End:                      0x00,
				},
			},
		},
	}

	for i, tc := range testCases {
		got, err := Read(tc.rs)
		if err != nil {
			t.Fatalf("test %d: failed to Read file: %v", i+1, err)
		}

		if !reflect.DeepEqual(tc.expected, got) {
			t.Errorf("test %d:\nexpected %+v,\nbut got %+v", i+1, tc.expected, got)
		}
	}
}

func TestInternalRead(t *testing.T) {
	testCases := []struct {
		got      any
		filename string
		expected any
		config   sectionConfig
	}{
		{
			got:      Header{},
			filename: "./assets/test.tga",
			config:   headerSection,
			expected: Header{
				IDLength:        0,
				ColorMapType:    0,
				ImageType:       UncompressedRGBImage,
				ColorMapOrigin:  0,
				ColorMapLength:  0,
				ColorMapDepth:   0,
				XOrigin:         0,
				YOrigin:         0,
				Width:           256,
				Height:          256,
				BitsPerPixel:    32,
				ImageDescriptor: 8,
			},
		},
		{
			got:      Header{},
			filename: "./assets/test2.tga",
			config:   headerSection,
			expected: Header{
				IDLength:        0,
				ColorMapType:    0,
				ImageType:       UncompressedRGBImage,
				ColorMapOrigin:  0,
				ColorMapLength:  0,
				ColorMapDepth:   0,
				XOrigin:         0,
				YOrigin:         0,
				Width:           1280,
				Height:          853,
				BitsPerPixel:    24,
				ImageDescriptor: 0,
			},
		},
		{
			got:      Header{},
			filename: "./assets/flag_t16.tga",
			config:   headerSection,
			expected: Header{
				IDLength:        0,
				ColorMapType:    0,
				ImageType:       UncompressedRGBImage,
				ColorMapOrigin:  0,
				ColorMapLength:  0,
				ColorMapDepth:   0,
				XOrigin:         0,
				YOrigin:         0,
				Width:           124,
				Height:          124,
				BitsPerPixel:    16,
				ImageDescriptor: 32,
			},
		},
		{
			got:      Header{},
			filename: "./assets/xing_t24.tga",
			config:   headerSection,
			expected: Header{
				IDLength:        0,
				ColorMapType:    0,
				ImageType:       UncompressedRGBImage,
				ColorMapOrigin:  0,
				ColorMapLength:  0,
				ColorMapDepth:   0,
				XOrigin:         0,
				YOrigin:         0,
				Width:           240,
				Height:          164,
				BitsPerPixel:    24,
				ImageDescriptor: 32,
			},
		},
		{
			got:      Footer{},
			filename: "./assets/test.tga",
			config:   footerSection,
			expected: Footer{
				ExtensionAreaOffset:      0,
				DeveloperDirectoryOffset: 0,
				Signature:                [16]byte{'T', 'R', 'U', 'E', 'V', 'I', 'S', 'I', 'O', 'N', '-', 'X', 'F', 'I', 'L', 'E'},
				Point:                    '.',
				End:                      0x00,
			},
		},
		{
			got:      Footer{},
			filename: "./assets/test2.tga",
			config:   footerSection,
			expected: Footer{
				ExtensionAreaOffset:      1450344536,
				DeveloperDirectoryOffset: 2220323969,
				Signature:                [16]byte{115, 87, 133, 114, 87, 131, 112, 85, 129, 111, 85, 125, 112, 87, 125, 114},
				Point:                    90,
				End:                      125,
			},
		},
		{
			got:      Footer{},
			filename: "./assets/flag_t16.tga",
			config:   footerSection,
			expected: Footer{
				ExtensionAreaOffset:      2145419232,
				DeveloperDirectoryOffset: 2082438175,
				Signature:                [16]byte{31, 124, 31, 124, 31, 124, 31, 124, 31, 124, 31, 124, 31, 124, 31, 124},
				Point:                    31,
				End:                      124,
			},
		},
		{
			got:      Footer{},
			filename: "./assets/xing_t24.tga",
			config:   footerSection,
			expected: Footer{
				ExtensionAreaOffset:      1350363959,
				DeveloperDirectoryOffset: 3328740322,
				Signature:                [16]byte{128, 101, 146, 82, 56, 100, 94, 82, 86, 51, 49, 32, 51, 49, 32, 49},
				Point:                    52,
				End:                      34,
			},
		},
	}

	_ = testCases

	for i, tc := range testCases {
		f, err := os.Open(tc.filename)
		if err != nil {
			t.Fatalf("test %d: failed to open test file: %v", i+1, err)
		}
		defer f.Close()

		switch tc.got.(type) {
		case Header:
			var got Header
			err = read(f, tc.config, &got)
			tc.got = got
		case Footer:
			var got Footer
			err = read(f, tc.config, &got)
			tc.got = got
		}

		if err != nil {
			t.Fatalf("failed to read test file: %v", err)
		}

		if !reflect.DeepEqual(tc.expected, tc.got) {
			t.Errorf("file `%s`:\nexpected %+v, \nbut got  %+v", tc.filename, tc.expected, tc.got)
		}
	}
}

func TestPixelPosition(t *testing.T) {
	testCases := []struct {
		imageDescriptor ImageDescriptor
		expected        ImageOrigin
	}{
		{imageDescriptor: 0, expected: BottomLeft},
		{imageDescriptor: 15, expected: BottomLeft},
		{imageDescriptor: 16, expected: BottomRight},
		{imageDescriptor: 32, expected: TopLeft},
		{imageDescriptor: 37, expected: TopLeft},
		{imageDescriptor: 48, expected: TopRight},
		{imageDescriptor: 53, expected: TopRight},
		{imageDescriptor: 127, expected: TopRight},
	}

	for i, tc := range testCases {
		got := tc.imageDescriptor.ImageOrigin()

		if tc.expected != got {
			t.Errorf("test %d: expected `%s`, but got `%s`", i+1, tc.expected, got)
		}
	}
}
