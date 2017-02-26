package jsonapi

import (
	"os"
	"testing"

	"github.com/DexterLB/mvm/testutils"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.RecordHTTP(m, "fixtures/jsonapi"))
}
