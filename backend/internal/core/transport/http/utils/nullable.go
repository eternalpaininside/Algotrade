package corehttputils

import (
	"encoding/json/v2"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n Nullable[T]) UnmarshalJSON(
	data []byte,
) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(
		data,
		&value,
	); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
