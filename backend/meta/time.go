package meta

import "time"

var startup time.Time

func Setup() {
	startup = time.Now()
}

func Timestamp() int64 {
	return time.Now().Sub(startup).Milliseconds()
}
