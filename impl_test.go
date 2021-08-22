package tigopesa_test

import (
	"github.com/techcraftlabs/tigopesa"
	"testing"
)

func TestInterfaceImplementations(t *testing.T) {
	t.Run("test if tigopesa.Client implements tigopesa.Service",
		func(t *testing.T) {
			var i interface{} = new(tigopesa.Client)
			if _, ok := i.(tigopesa.Service); !ok {
				t.Fatalf("expected %t to implement Service", i)
			}
		})

}
