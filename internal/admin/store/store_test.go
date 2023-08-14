package adminstore_test

import (
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
)

var testDB *db.DB

func TestMain(m *testing.M) {
	testDB = dbtest.Prepare()

	os.Exit(m.Run())
}

func ptr[T any](t T) *T {
	return &t
}
