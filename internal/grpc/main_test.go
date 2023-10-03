package grpc

import (
	"context"
	"os"
	"testing"

	"github.com/kien-manabie/rio-mocker/internal/test"
)

func TestMain(m *testing.M) {
	test.ResetDB(context.Background(), "../..")

	code := m.Run()
	os.Exit(code)
}
