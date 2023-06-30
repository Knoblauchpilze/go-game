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
