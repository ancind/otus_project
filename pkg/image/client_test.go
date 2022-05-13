package image

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

var defaultImgURL = "http://raw.githubusercontent.com/ancind/otus_project/master/tests/static/"

func Test_Download_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
	}{
		{
			name:    "success_get_img",
			ctx:     ctx,
			imgName: "gopher.jpg",
		},
	}

	logger := log.With().Logger()
	id := NewImageGetter(logger, 15*time.Second, 15*time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := id.Get(tt.ctx, defaultImgURL+tt.imgName, nil)
			if err != nil {
				t.Errorf("Download() error = %v", err)
				return
			}

			wantImg := loadImage(tt.imgName)
			if !reflect.DeepEqual(gotImg, wantImg) {
				t.Errorf("Download() gotImg = %v, want %v", gotImg, wantImg)
			}
		})
	}
}

func Test_Download_Negative(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeOut, c := context.WithTimeout(ctx, 5*time.Microsecond)
	defer c()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
		url     string
		err     error
	}{
		{
			name:    "timeout_case",
			ctx:     ctxWithTimeOut,
			imgName: "gopher.jpg",
			url:     defaultImgURL,
			err:     nil,
		},
		{
			name:    "not_allowed_type_img_case",
			ctx:     ctx,
			imgName: "text.txt",
			url:     defaultImgURL,
			err:     nil,
		},
		{
			name:    "bad_response_case",
			ctx:     ctx,
			imgName: "",
			url:     defaultImgURL,
			err:     nil,
		},
	}

	logger := log.With().Logger()
	id := NewImageGetter(logger, 5*time.Second, 5*time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := id.Get(tt.ctx, tt.url+tt.imgName, nil)
			require.Error(t, err, tt.err)
		})
	}
}
