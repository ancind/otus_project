package image

import (
	"context"
	"reflect"
	"testing"
)

func Test_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ctx         context.Context
		originalImg []byte
		resizedImg  []byte
		width       int
		height      int
	}{
		{
			name:        "success resize 256x126",
			ctx:         ctx,
			width:       235,
			height:      545,
			originalImg: loadImage("gopher.jpg"),
			resizedImg:  loadImage("gopher_333x666.jpg"),
		},
	}

	id := NewResizer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := id.Resize(tt.originalImg, tt.width, tt.height)
			if err != nil {
				t.Errorf("Resize() error = %v", err)
				return
			}

			if !reflect.DeepEqual(gotImg, tt.resizedImg) {
				t.Errorf("Resize()  gotImg = %v, want %v", gotImg, tt.resizedImg)
			}
		})
	}
}
