package corehttprequest

import (
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(
	r *http.Request,
	dest any,
) error {
	r.Body = http.MaxBytesReader(
		nil,
		r.Body,
		1<<20,
	)
	if err := jsonv2.UnmarshalRead(
		r.Body,
		dest,
		jsonv2.RejectUnknownMembers(true),
	); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](
			err,
		); ok {
			return fmt.Errorf(
				"%w: %w",
				coreerrors.ErrTooLarge,
				err,
			)
		}

		return fmt.Errorf(
			"%w: json unmarshal read: %w",
			coreerrors.ErrInvalidArgument,
			err,
		)
	}

	if err := requestValidator.Struct(
		dest,
	); err != nil {
		if _, ok := errors.AsType[validator.ValidationErrors](
			err,
		); ok {
			return fmt.Errorf(
				"%w: %w",
				coreerrors.ErrInvalidArgument,
				err,
			)
		}
		return fmt.Errorf(
			"request validation: %w",
			err,
		)
	}

	return nil
}
