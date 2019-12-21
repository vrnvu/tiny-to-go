package tiny

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect not found")
	ErrRedirectInvalid  = errors.New("Redirect invalid")
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (service *service) Find(code string) (*Redirect, error) {
	return service.repository.Find(code)
}

func (service *service) Save(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.Timestamp = time.Now().UTC().Unix()
	return service.repository.Save(redirect)
}
