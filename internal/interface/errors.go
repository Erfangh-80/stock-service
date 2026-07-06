package iface

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrNotFound     Error = "not found"
	ErrInvalidInput Error = "invalid input"
	ErrConflict     Error = "conflict"
	ErrInternal     Error = "internal error"
)
