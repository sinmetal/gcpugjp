package backend

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/favclip/testerator"
	_ "github.com/favclip/testerator/datastore"
	_ "github.com/favclip/testerator/memcache"
	_ "github.com/favclip/testerator/search"
)

func TestMain(m *testing.M) {
	_, _, err := testerator.SpinUp() // 最初の1プロセスを起動！

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run() // UnitTest実行！

	err = testerator.SpinDown() // 最初に立ち上げたプロセスを落とす
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// EqualTime is 日付を比べる時に、Millisecondまでの制度でのみ比べるEqual
// 引数にZeroTimeを渡した場合は、必ずfalseが返ってくるので、ZeroTime通しをEqualで比べることはできない
func EqualTime(t1 time.Time, t2 time.Time) bool {
	if t1.IsZero() || t2.IsZero() {
		return false
	}
	return t1.Truncate(time.Millisecond).Equal(t2.Round(time.Millisecond))
}
