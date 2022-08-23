package tga

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		sections [][]byte
		expected File
	}{
		{
			sections: [][]byte{
				{0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 16, 0},
				{12, 12},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 'R', 'U', 'E', 'V', 'I', 'S', 'I', 'O', 'N', '-', 'X', 'F', 'I', 'L', 'E', '.', 0x00},
			},
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
					BitsPerPixel:    16,
					ImageDescriptor: 0,
				},
				Image: Image{
					ID:       []byte{},
					ColorMap: []byte{},
					Data:     []byte{12, 12},
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
		input := bytes.NewReader(bytes.Join(tc.sections, nil))

		got, err := Read(input)
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

func TestPixelAt(t *testing.T) {
	testCases := []struct {
		img      File
		x        int
		y        int
		expected []byte
	}{
		{
			x:        0,
			y:        0,
			expected: []byte{0, 1},
			img: File{
				Header: Header{
					Width:        10,
					Height:       10,
					BitsPerPixel: 16,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
						16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
						32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
						48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
						64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
						80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
						96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
						112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
						128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
						144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
						160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
						176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191,
						192, 193, 194, 195, 196, 197, 198, 199,
					},
				},
			},
		},
		{
			x:        5,
			y:        0,
			expected: []byte{10, 11},
			img: File{
				Header: Header{
					Width:        10,
					Height:       10,
					BitsPerPixel: 16,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
						16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
						32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
						48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
						64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
						80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
						96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
						112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
						128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
						144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
						160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
						176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191,
						192, 193, 194, 195, 196, 197, 198, 199,
					},
				},
			},
		},
		{
			x:        5,
			y:        5,
			expected: []byte{110, 111},
			img: File{
				Header: Header{
					Width:        10,
					Height:       10,
					BitsPerPixel: 16,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
						16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
						32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
						48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
						64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
						80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
						96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
						112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
						128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
						144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
						160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
						176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191,
						192, 193, 194, 195, 196, 197, 198, 199,
					},
				},
			},
		},
		{
			x:        1,
			y:        1,
			expected: []byte{12, 13, 14},
			img: File{
				Header: Header{
					Width:        3,
					Height:       3,
					BitsPerPixel: 24,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 3, 4, 5, 6, 7, 8,
						9, 10, 11, 12, 13, 14, 15, 16, 17,
						18, 19, 20, 21, 22, 23, 24, 25, 26,
					},
				},
			},
		},
		{
			x:        3,
			y:        3,
			expected: []byte{33},
			img: File{
				Header: Header{
					Width:        10,
					Height:       10,
					BitsPerPixel: 8,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
						10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
						20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
						30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
						40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
						50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
						60, 61, 62, 63, 64, 65, 66, 67, 68, 69,
						70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
						80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
						90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
					},
				},
			},
		},
		{
			x:        3,
			y:        1,
			expected: []byte{29, 30, 31, 32},
			img: File{
				Header: Header{
					Width:        4,
					Height:       4,
					BitsPerPixel: 32,
				},
				Image: Image{
					Data: []byte{
						0, 1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
						17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
						33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
						49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
						65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		got := tc.img.PixelAt(tc.x, tc.y)

		if !reflect.DeepEqual(tc.expected, got) {
			t.Errorf("test %d: expected `%v`, but got `%v`", i+1, tc.expected, got)
		}
	}
}
