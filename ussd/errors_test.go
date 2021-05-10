package ussd_test

import (
	"github.com/techcraftt/tigosdk/ussd"
	"testing"
)

func TestTxnStatus(t *testing.T) {
	st := ussd.TxnStatus
	for key, val := range st {
		t.Logf("key: %s, value: %s",key,val)
	}
}
