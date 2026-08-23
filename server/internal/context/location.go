package context

import (
	stdcontext "context"
	"time"
)

type locationKey struct{}

func SetLocation(ctx stdcontext.Context, loc *time.Location) stdcontext.Context {
	return stdcontext.WithValue(ctx, locationKey{}, loc)
}

func GetLocation(ctx stdcontext.Context) *time.Location {
	loc, ok := ctx.Value(locationKey{}).(*time.Location)
	if ok {
		return loc
	}

	return nil
}
