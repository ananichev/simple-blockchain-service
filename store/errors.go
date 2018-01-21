package store

import "fmt"

// ErrGetLast is returned when there is an error retrieving last blocks
type ErrGetLast struct {
	err error
}

func (e ErrGetLast) Error() string {
	return fmt.Sprintf("GetLast error: %s", e.err)
}

// ErrAddBlock is returned when there is an error storing block
type ErrAddBlock struct {
	err error
}

func (e ErrAddBlock) Error() string {
	return fmt.Sprintf("AddBlock error: %s", e.err)
}
