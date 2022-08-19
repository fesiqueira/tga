package tga

import (
	"os"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		filename string
		expected Header
	}{
		{
			filename: "./assets/test.tga",
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
			filename: "./assets/test2.tga",
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
			filename: "./assets/flag_t16.tga",
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
			filename: "./assets/xing_t24.tga",
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
	}

	for i, tc := range testCases {
		f, err := os.Open(tc.filename)
		if err != nil {
			t.Fatalf("test %d: failed to open test file: %v", i+1, err)
		}
		defer f.Close()

		header, err := Read(f)
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
		}

		if !reflect.DeepEqual(tc.expected, header) {
			t.Errorf("\nexpected %+v, \nbut got  %+v", tc.expected, header)
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
