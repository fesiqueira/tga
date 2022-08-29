package tga

import (
	"image"
	"os"
	"reflect"
	"testing"
)

var testFiles = []string{
	"test.tga",
	"test2.tga",
	"flag_t16.tga",
	"xing_t24.tga",
}

func decodeTGA(filename string) (image.Image, error) {
	f, err := os.Open("./testdata/" + filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Decode(f)
}

func TestDecode(t *testing.T) {
	testCases := map[string]struct {
		filename string
		want     image.Image
	}{
		"DecodeTGA32BottomLeft": {
			filename: "test.tga",
			want:     nil,
		},
		// "DecodeTGA24BottomLeft": {
		//	filename: "test2.tga",
		//	want:     nil,
		// },
		// "DecodeTGA16TopLeft": {
		//	filename: "flag_t16.tga",
		//	want:     nil,
		// },
		// "DecodeTGA24TopLeft": {
		//	filename: "xing_t24.tga",
		//	want:     nil,
		// },
	}

	for name, tc := range testCases {
		got, err := decodeTGA(tc.filename)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("%s: different image from what is expected", name)
		}
	}
}
