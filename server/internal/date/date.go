package date

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// Date 日付を表現する
// 数値は負の値をとらないが、`time.Time`に合わせ`int`を採用
type Date struct {
	year  int
	month time.Month
	day   int
}

var (
	_       sql.Scanner   = (*Date)(nil)
	_       driver.Valuer = (*Date)(nil)
	minDate               = Min()
)

func NewDate(year int, month time.Month, day int) (Date, error) {
	d := Date{
		year:  year,
		month: month,
		day:   day,
	}
	if err := d.validate(); err != nil {
		return Date{}, err
	}

	return d, nil
}

func Min() Date {
	return Date{
		year:  1,
		month: time.January,
		day:   1,
	}
}

func Max() Date {
	return Date{
		year:  9999,
		month: time.December,
		day:   31,
	}
}

func Parse(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, newError(fmt.Errorf("invalid date string: %q: %w", s, err))
	}

	d, err := NewDate(t.Year(), t.Month(), t.Day())
	if err != nil {
		return Date{}, newError(err)
	}

	return d, nil
}

func (d Date) validate() error {
	t := time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.Local)
	if t.Year() != d.year ||
		t.Month() != d.month ||
		t.Day() != d.day {
		return newError(fmt.Errorf("invalid date: got = %d-%d-%d", d.year, d.month, d.day))
	}

	return nil
}

func (d Date) Year() int {
	return d.year
}

func (d Date) Month() time.Month {
	return d.month
}

func (d Date) Day() int {
	return d.day
}

func (d Date) Time() time.Time {
	return time.Date(
		d.Year(),
		d.Month(),
		d.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

// Scan `database/sql`の`sql.Scanner`を実装し、`DATE`カラムの値をマッピングできるようにする
// ドライバの設定により`time.Time`（`parseTime=true`）または文字列で渡されるため、両方を受け付ける
func (d *Date) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*d = Date{}
		return nil
	case time.Time:
		scanned, err := NewDate(v.Year(), v.Month(), v.Day())
		if err != nil {
			return newError(err)
		}
		*d = scanned
		return nil
	case []byte:
		return d.scanString(string(v))
	case string:
		return d.scanString(v)
	default:
		return newError(fmt.Errorf("unsupported scan type for Date: %T", value))
	}
}

func (d *Date) scanString(s string) error {
	scanned, err := Parse(s)
	if err != nil {
		return err
	}
	*d = scanned

	return nil
}

// Value `database/sql/driver`の`driver.Valuer`を実装し、`DATE`カラムへ書き込めるよう文字列表現を返す
func (d Date) Value() (driver.Value, error) {
	if err := d.validate(); err != nil {
		return nil, newError(err)
	}

	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day), nil
}

// IsPast 基準時刻`now`のカレンダー日より前の日付であれば`true`を返す
// 時刻は持たないため、`now`のロケーションにおける日単位で比較する
func (d Date) IsPast(now time.Time) bool {
	loc := now.Location()
	self := time.Date(d.year, d.month, d.day, 0, 0, 0, 0, loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	return self.Before(today)
}

func (d Date) IsMin() bool {
	return d.year == minDate.year && d.month == minDate.month && d.day == minDate.day
}

func (d Date) IsEqual(target Date) bool {
	return d.year == target.year && d.month == target.month && d.day == target.day
}
