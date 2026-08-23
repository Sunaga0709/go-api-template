package location

import (
	"sync"
	"time"
)

var (
	jstLocation   *time.Location
	jstLocationMu sync.Mutex
)

func init() {
	jstLocation = jst()
}

func JST() *time.Location {
	jstLocationMu.Lock()
	defer jstLocationMu.Unlock()

	if jstLocation != nil {
		return jstLocation
	}
	jstLocation = jst()

	return jstLocation
}

func jst() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}

	return loc
}
