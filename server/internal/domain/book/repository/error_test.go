package repository

import (
	"errors"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	cause := errors.New("cause")
	tests := []struct {
		kind ErrorKind
		want string
		new  func(error) Error
	}{{ErrorKindUnknown, errorKindUnknownString, NewUnknownError}, {ErrorKindGet, errorKindGetString, NewGetError}, {ErrorKindNotFound, errorKindNotFoundString, NewNotFoundError}, {ErrorKindCreate, errorKindCreateString, NewCreateError}}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if tt.kind.String() != tt.want {
				t.Errorf("String() = %q", tt.kind.String())
			}
			got := tt.new(cause)
			if got.Kind != tt.kind || !strings.Contains(got.Error(), "cause") {
				t.Errorf("constructor() = %#v", got)
			}
		})
	}
	if newError(ErrorKind(99), cause).Kind.String() != errorKindUnknownString {
		t.Error("unknown kind mismatch")
	}
}
