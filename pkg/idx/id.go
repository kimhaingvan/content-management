package idx

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"time"
)

type Timestamp int64
type ID int64

const e6 = 1e6
const e3 = 1e3
const magicNumber = int64(1012345678909e6>>24 + 1)

var zeroDate = time.Date(2018, 02, 15, 17, 0, 0, 0, time.UTC)
var startDate = FromMillis(zeroDate.UnixNano()/1e6 - magicNumber*10)

func NewID() ID {
	return ID(NewIDWithTime(time.Now()))
}

// NewIDWithTime ...
func NewIDWithTime(t time.Time) int64 {
	var b [3]byte
	_, err := rand.Read(b[:])
	if err != nil {
	}

	d := int64(t.Sub(startDate)) / 1e7
	id := d<<24 | int64(b[0])<<16 | int64(b[1])<<8 | int64(b[2])
	return id & ^int64(1<<24)
}

func FromMillis(t int64) time.Time {
	return Timestamp(t).ToTime()
}

// ToTime converts from timestamp to Go time.
func (t Timestamp) ToTime() time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(int64(t)/1e3, int64(t%1e3)*e6).UTC()
}

// Unix converts Timestamp to seconds
func (t Timestamp) Unix() int64 {
	return int64(t) / 1e3
}

// UnixNano extracts nanoseconds from Timestamp
func (t Timestamp) UnixNano() int64 {
	return int64(t) * e6
}

func (i *ID) String() string {
	return fmt.Sprintf("%v", i)
}

func (i *ID) ParseID(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	return ID(id), err
}
