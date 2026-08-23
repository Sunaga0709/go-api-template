//go:build !integration

package context

import (
	stdcontext "context"
	"reflect"
	"testing"
	"time"
)

func TestSetLocation(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("JST", 9*60*60)

	type args struct {
		ctx stdcontext.Context
		loc *time.Location
	}
	tests := []struct {
		name string
		args args
		want *time.Location
	}{
		// 正常系
		{
			name: "set location",
			args: args{
				ctx: stdcontext.Background(),
				loc: jst,
			},
			want: jst,
		},
		{
			name: "set nil location",
			args: args{
				ctx: stdcontext.Background(),
				loc: nil,
			},
			want: nil,
		},
		{
			name: "overwrite existing location",
			args: args{
				ctx: SetLocation(stdcontext.Background(), time.UTC),
				loc: jst,
			},
			want: jst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SetLocation(tt.args.ctx, tt.args.loc)
			if loc := GetLocation(got); !reflect.DeepEqual(loc, tt.want) {
				t.Errorf("SetLocation() = %v, want %v", loc, tt.want)
			}
		})
	}
}

func TestGetLocation(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("JST", 9*60*60)

	type args struct {
		ctx stdcontext.Context
	}
	tests := []struct {
		name string
		args args
		want *time.Location
	}{
		// 正常系
		{
			name: "location is set",
			args: args{
				ctx: SetLocation(stdcontext.Background(), jst),
			},
			want: jst,
		},
		// 異常系
		{
			name: "location is not set",
			args: args{
				ctx: stdcontext.Background(),
			},
			want: nil,
		},
		{
			name: "value is not *time.Location",
			args: args{
				ctx: stdcontext.WithValue(stdcontext.Background(), locationKey{}, "not a location"),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := GetLocation(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
