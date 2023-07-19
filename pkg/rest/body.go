package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

func GetBodyFromHttpRequestAs(req *http.Request, out interface{}) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return errors.WrapCode(err, errors.ErrFailedToGetBody)
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return errors.WrapCode(err, errors.ErrBodyParsingFailed)
	}

	return nil
}

func GetBodyFromHttpResponseAs(resp *http.Response, out interface{}) error {
	if resp == nil {
		return errors.NewCode(errors.ErrNoResponse)
	}
	if resp.Body == nil {
		if resp.StatusCode != http.StatusOK {
			return errors.WrapCode(errors.New(http.StatusText(resp.StatusCode)), errors.ErrResponseIsError)
		}
		return errors.NewCode(errors.ErrFailedToGetBody)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.WrapCode(err, errors.ErrFailedToGetBody)
	}

	var in ResponseTemplate
	err = json.Unmarshal(data, &in)
	if err != nil {
		return errors.WrapCode(err, errors.ErrBodyParsingFailed)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.WrapCode(errors.New(string(in.Details)), errors.ErrResponseIsError)
	}

	if err = json.Unmarshal(in.Details, out); err != nil {
		return errors.WrapCode(err, errors.ErrBodyParsingFailed)
	}

	return nil
}
