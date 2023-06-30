package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
)

type ResponseBuilder interface {
	Pass() ResponseBuilder
	Fail() ResponseBuilder

	WithDetails(details interface{}) ResponseBuilder
	WithCode(httpCode int) ResponseBuilder

	Write(w http.ResponseWriter)
}

type responseImpl struct {
	RequestId uuid.UUID
	Status    string
	Details   json.RawMessage `json:",omitempty"`
	code      int
}

var StatusOK = "SUCCESS"
var StatusNOK = "ERROR"

func NewSuccessResponse(id uuid.UUID) ResponseBuilder {
	return &responseImpl{
		RequestId: id,
		Status:    StatusOK,
		code:      http.StatusOK,
	}
}

func NewErrorResponse(id uuid.UUID) ResponseBuilder {
	return &responseImpl{
		RequestId: id,
		Status:    StatusNOK,
		code:      http.StatusOK,
	}
}

func (ri *responseImpl) Pass() ResponseBuilder {
	ri.Status = StatusOK
	return ri
}

func (ri *responseImpl) Fail() ResponseBuilder {
	ri.Status = StatusNOK
	return ri
}

func (ri *responseImpl) WithDetails(details interface{}) ResponseBuilder {
	var out []byte
	var err error

	if out, err = json.Marshal(details); err == nil {
		ri.Details = out
	}

	return ri
}

func (ri *responseImpl) WithCode(httpCode int) ResponseBuilder {
	ri.code = httpCode
	if ri.code != http.StatusOK {
		return ri.Fail()
	}

	return ri.Pass()
}

func (ri *responseImpl) Write(w http.ResponseWriter) {
	out, _ := json.Marshal(ri)

	w.WriteHeader(ri.code)
	w.Write(out)
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

	var in responseImpl
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
