// Package oaci is the loader implementation to grab tiles from Geoportail.
package oaci

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/cenkalti/backoff"
	"github.com/landru29/mbtiles/internal/model"
	pkgerrors "github.com/pkg/errors"
)

//go:generate mockgen -destination=../../mocks/oaci.go -package=oaci net/http RoundTripper

// Client is the Geoportail (OACI) loader.
type Client struct {
	transport http.RoundTripper
}

// Configurator is the client option.
type Configurator func(*Client)

// New creates a new client.
func New(options ...Configurator) Client {
	output := Client{}

	for _, opt := range options {
		opt(&output)
	}

	return output
}

// WithTransport is for testing purpose.
func WithTransport(transport http.RoundTripper) Configurator {
	return func(client *Client) {
		client.transport = transport
	}
}

// LoadImage implements the tile.Loader interface.
func (c Client) LoadImage(ctx context.Context, request model.TileRequest) (image.Image, error) {
	client := http.Client{
		Transport: c.transport,
	}

	picURL, err := url.Parse("https://data.geopf.fr/private/wmts")
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot parse tile source url")
	}

	offset := uint64(1)
	for range request.ZoomLevel {
		offset *= 2
	}

	row := offset - 1 - request.Row

	values := url.Values{}
	values.Add("apikey", "geoportail")
	values.Add("layer", "GEOGRAPHICALGRIDSYSTEMS.MAPS.SCAN-OACI")
	values.Add("style", "normal")
	values.Add("tilematrixset", "PM")
	values.Add("Service", "WMTS")
	values.Add("Request", "GetTile")
	values.Add("Version", "1.0.0")
	values.Add("Format", "image/jpeg")
	values.Add("TileMatrix", strconv.FormatUint(request.ZoomLevel, 10))
	values.Add("TileCol", strconv.FormatUint(request.Col, 10))
	values.Add("TileRow", strconv.FormatUint(row, 10))

	picURL.RawQuery = values.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, picURL.String(), nil)
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot build request to the tile source server")
	}

	httpRequest.Header.Add("Referer", "https://www.geoportail.gouv.fr/")

	resp, err := client.Do(httpRequest)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot GET image")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		return nil, pkgerrors.WithMessage(
			os.ErrNotExist,
			fmt.Sprintf("[🔍%d ↓%d →%d]", request.ZoomLevel, request.Col, row),
		)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, pkgerrors.WithMessage(backoff.Permanent(err), "cannot decode JPEG image")
	}

	return img, nil
}
