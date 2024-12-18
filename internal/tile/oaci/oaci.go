package oaci

import (
	"context"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"

	pkgerrors "github.com/pkg/errors"

	"github.com/cenkalti/backoff"
)

type Client struct{}

func (c Client) LoadImage(ctx context.Context, zoomLevel uint64, col uint64, row uint64) (image.Image, error) {
	client := http.Client{}

	picUrl, err := url.Parse("https://data.geopf.fr/private/wmts")
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot parse tile source url")
	}

	values := url.Values{}
	values.Add("apikey", "geoportail")
	values.Add("layer", "GEOGRAPHICALGRIDSYSTEMS.MAPS.SCAN-OACI")
	values.Add("style", "normal")
	values.Add("tilematrixset", "PM")
	values.Add("Service", "WMTS")
	values.Add("Request", "GetTile")
	values.Add("Version", "1.0.0")
	values.Add("Format", "image/jpeg")
	values.Add("TileMatrix", strconv.FormatUint(zoomLevel, 10))
	values.Add("TileCol", strconv.FormatUint(col, 10))
	values.Add("TileRow", strconv.FormatUint(row, 10))

	picUrl.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, picUrl.String(), nil)
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot build request to the tile source server")
	}

	request.Header.Add("Referer", "https://www.geoportail.gouv.fr/")

	resp, err := client.Do(request)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot GET image")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		return nil, os.ErrNotExist
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot decode JPEG image")
	}

	return img, nil
}
