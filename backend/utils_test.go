package backend

import "time"

// EqualTime is 日付を比べる時に、Millisecondまでの制度でのみ比べるEqual
// 引数にZeroTimeを渡した場合は、必ずfalseが返ってくるので、ZeroTime通しをEqualで比べることはできない
func EqualTime(t1 time.Time, t2 time.Time) bool {
	if t1.IsZero() || t2.IsZero() {
		return false
	}
	return t1.Truncate(time.Millisecond).Equal(t2.Round(time.Millisecond))
}
