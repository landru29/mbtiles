package oaci_test

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"testing"

	"github.com/landru29/mbtiles/internal/matcher"
	mockoaci "github.com/landru29/mbtiles/internal/mocks"
	"github.com/landru29/mbtiles/internal/model"
	"github.com/landru29/mbtiles/internal/tile/oaci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoad(t *testing.T) {
	t.Run("coordinate shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockTrip := mockoaci.NewMockRoundTripper(ctrl)

		// Prepare jpeg image in a buffer.
		buffer := bytes.Buffer{}
		require.NoError(
			t,
			jpeg.Encode(
				&buffer,
				image.NewRGBA(image.Rectangle{
					Min: image.Point{X: 0, Y: 0},
					Max: image.Point{X: 256, Y: 256},
				}),
				&jpeg.Options{Quality: 100},
			),
		)

		// Mock the response.
		mockTrip.EXPECT().RoundTrip(matcher.NewRequest(
			matcher.RequestWithValue("TileMatrix", "10"),
			matcher.RequestWithValue("TileCol", "500"),
			matcher.RequestWithValue("TileRow", "353"),
		)).Return(&http.Response{
			Body:       io.NopCloser(&buffer),
			StatusCode: http.StatusOK,
		}, nil)

		imgOut, err := oaci.New(oaci.WithTransport(mockTrip)).LoadImage(
			ctx,
			model.TileRequest{
				ZoomLevel: 10,
				Col:       500,
				Row:       670,
			},
		)

		require.NoError(t, err)

		assert.EqualValues(t, 256, imgOut.Bounds().Max.X)
		assert.EqualValues(t, 256, imgOut.Bounds().Max.Y)
	})
}
