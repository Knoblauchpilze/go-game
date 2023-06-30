package rest

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

func GetHeaderFromHttpRequest(req *http.Request, headerKey string) ([]string, error) {
	header, ok := req.Header[headerKey]
	if !ok {
		return nil, errors.NewCode(errors.ErrNoSuchHeader)
	}

	return header, nil
}

func GetSingleHeaderFromHttpRequest(req *http.Request, headerKey string) (string, error) {
	header, err := GetHeaderFromHttpRequest(req, headerKey)
	if err != nil {
		return "", err
	}

	if len(header) > 1 {
		return "", errors.NewCode(errors.ErrNonUniqueHeader)
	}

	return header[0], nil
}
