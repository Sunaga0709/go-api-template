//go:build !integration

package date

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	t.Parallel()

	type args struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name    string
		args    args
		want    Date
		wantErr bool
	}{
		// 正常系
		{
			name: "valid date",
			args: args{year: 2026, month: time.June, day: 27},
			want: Date{year: 2026, month: time.June, day: 27},
		},
		{
			name: "leap day",
			args: args{year: 2024, month: time.February, day: 29},
			want: Date{year: 2024, month: time.February, day: 29},
		},
		// 境界値
		{
			name: "first day of month",
			args: args{year: 2026, month: time.January, day: 1},
			want: Date{year: 2026, month: time.January, day: 1},
		},
		{
			name: "last day of month",
			args: args{year: 2026, month: time.December, day: 31},
			want: Date{year: 2026, month: time.December, day: 31},
		},
		// 異常系
		{
			name:    "day overflow (April 31)",
			args:    args{year: 2026, month: time.April, day: 31},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "non-leap year February 29",
			args:    args{year: 2026, month: time.February, day: 29},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "day zero",
			args:    args{year: 2026, month: time.June, day: 0},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "month out of range",
			args:    args{year: 2026, month: time.Month(13), day: 1},
			want:    Date{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewDate(tt.args.year, tt.args.month, tt.args.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Date
		wantErr bool
	}{
		// 正常系
		{
			name: "valid date",
			args: args{s: "2026-06-27"},
			want: Date{year: 2026, month: time.June, day: 27},
		},
		{
			name: "leap day",
			args: args{s: "2024-02-29"},
			want: Date{year: 2024, month: time.February, day: 29},
		},
		// 境界値
		{
			name: "first day of month",
			args: args{s: "2026-01-01"},
			want: Date{year: 2026, month: time.January, day: 1},
		},
		{
			name: "last day of month",
			args: args{s: "2026-12-31"},
			want: Date{year: 2026, month: time.December, day: 31},
		},
		// 異常系
		{
			name:    "invalid string format",
			args:    args{s: "2026/06/27"},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "month is not zero-padded",
			args:    args{s: "2026-6-27"},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "non-leap year February 29",
			args:    args{s: "2026-02-29"},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "day overflow",
			args:    args{s: "2026-04-31"},
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "empty string",
			args:    args{s: ""},
			want:    Date{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		date    Date
		wantErr bool
	}{
		// 正常系
		{
			name: "valid date",
			date: Date{year: 2026, month: time.June, day: 27},
		},
		// 異常系
		{
			name:    "day overflow",
			date:    Date{year: 2026, month: time.April, day: 31},
			wantErr: true,
		},
		{
			name:    "month out of range",
			date:    Date{year: 2026, month: time.Month(0), day: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := tt.date
			if err := d.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Date.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		want    Date
		wantErr bool
	}{
		// 正常系（time.Time）
		{
			name:  "time.Time",
			value: time.Date(2026, time.June, 27, 12, 0, 0, 0, time.UTC),
			want:  Date{year: 2026, month: time.June, day: 27},
		},
		// 正常系（[]byte）
		{
			name:  "byte slice",
			value: []byte("2026-06-27"),
			want:  Date{year: 2026, month: time.June, day: 27},
		},
		// 正常系（string）
		{
			name:  "string",
			value: "2024-02-29",
			want:  Date{year: 2024, month: time.February, day: 29},
		},
		// 境界値（nil はゼロ値）
		{
			name:  "nil",
			value: nil,
			want:  Date{},
		},
		// 異常系
		{
			name:    "invalid string format",
			value:   "2026/06/27",
			wantErr: true,
		},
		{
			name:    "invalid date value",
			value:   "2026-02-29",
			wantErr: true,
		},
		{
			name:    "unsupported type",
			value:   12345,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var d Date
			err := d.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if d != tt.want {
				t.Errorf("Date.Scan() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       Date
		want    driver.Value
		wantErr bool
	}{
		// 正常系
		{
			name: "valid date",
			d:    Date{year: 2026, month: time.June, day: 27},
			want: "2026-06-27",
		},
		{
			name: "single digit month and day are zero-padded",
			d:    Date{year: 2026, month: time.January, day: 1},
			want: "2026-01-01",
		},
		// 異常系
		{
			name:    "zero value is invalid",
			d:       Date{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.d.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Date.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsPast(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.June, 27, 12, 0, 0, 0, time.Local)

	type args struct {
		now time.Time
	}
	tests := []struct {
		name string
		d    Date
		args args
		want bool
	}{
		// 正常系（過去日）
		{
			name: "yesterday is past",
			d:    Date{year: 2026, month: time.June, day: 26},
			args: args{now: now},
			want: true,
		},
		{
			name: "previous year is past",
			d:    Date{year: 2025, month: time.December, day: 31},
			args: args{now: now},
			want: true,
		},
		// 境界値（当日は過去日ではない）
		{
			name: "today is not past",
			d:    Date{year: 2026, month: time.June, day: 27},
			args: args{now: now},
			want: false,
		},
		{
			name: "today at start of day is not past",
			d:    Date{year: 2026, month: time.June, day: 27},
			args: args{now: time.Date(2026, time.June, 27, 0, 0, 0, 0, time.Local)},
			want: false,
		},
		{
			name: "today at end of day is not past",
			d:    Date{year: 2026, month: time.June, day: 27},
			args: args{now: time.Date(2026, time.June, 27, 23, 59, 59, 0, time.Local)},
			want: false,
		},
		// 異常系相当（未来日）
		{
			name: "tomorrow is not past",
			d:    Date{year: 2026, month: time.June, day: 28},
			args: args{now: now},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.d.IsPast(tt.args.now); got != tt.want {
				t.Errorf("Date.IsPast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Year(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    Date
		want int
	}{
		// 正常系
		{
			name: "returns year",
			d:    Date{year: 2026, month: time.June, day: 27},
			want: 2026,
		},
		// 境界値（ゼロ値）
		{
			name: "zero value",
			d:    Date{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.d.Year(); got != tt.want {
				t.Errorf("Date.Year() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Month(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    Date
		want time.Month
	}{
		// 正常系
		{
			name: "returns month",
			d:    Date{year: 2026, month: time.June, day: 27},
			want: time.June,
		},
		// 境界値（ゼロ値）
		{
			name: "zero value",
			d:    Date{},
			want: time.Month(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.d.Month(); got != tt.want {
				t.Errorf("Date.Month() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Day(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    Date
		want int
	}{
		// 正常系
		{
			name: "returns day",
			d:    Date{year: 2026, month: time.June, day: 27},
			want: 27,
		},
		// 境界値（ゼロ値）
		{
			name: "zero value",
			d:    Date{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.d.Day(); got != tt.want {
				t.Errorf("Date.Day() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    Date
		want string
	}{
		// 正常系
		{
			name: "valid date",
			d:    Date{year: 2026, month: time.June, day: 27},
			want: "2026-06-27",
		},
		{
			name: "leap day",
			d:    Date{year: 2024, month: time.February, day: 29},
			want: "2024-02-29",
		},
		// 境界値
		{
			name: "single digit month and day are zero-padded",
			d:    Date{year: 2026, month: time.January, day: 1},
			want: "2026-01-01",
		},
		{
			name: "year is zero-padded to four digits",
			d:    Date{year: 1, month: time.January, day: 1},
			want: "0001-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := tt.d
			if got := d.String(); got != tt.want {
				t.Errorf("Date.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
