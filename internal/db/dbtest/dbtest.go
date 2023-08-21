// Package dbtest provides helper functions for tests that require database access
package dbtest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/joho/godotenv"
)

func Prepare() *db.DB {
	var (
		testDB *db.DB
		err    error
	)

	err = godotenv.Load("../../../.env.test")
	if err != nil {
		log.Printf("Did not load from .env.test file: %v", err)
	}

	testDB, err = db.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = testDB.AutoMigrate()
	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}

	return testDB
}

type Factory[T any] interface {
	Build(int) T
}

func Create[T any](t *testing.T, tx *db.DB, factory Factory[T], index int) T {
	t.Helper()

	m := factory.Build(index)
	if err := tx.Create(&m).Error; err != nil {
		t.Fatalf("Failed to create %T: %v", m, err)
	}

	return m
}

func CreateList[T any](t *testing.T, tx *db.DB, factory Factory[T], count int) []T {
	t.Helper()

	ms := make([]T, count)
	for i := range ms {
		ms[i] = Create(t, tx, factory, i)
	}

	return ms
}

// CreateApp is deprecated in favor of Create with AppFactory
func CreateApp(t *testing.T, tx *db.DB, index int, user *db.User) *db.App {
	t.Helper()

	if user == nil {
		user = CreateUser(t, tx, index)
	}

	app := &db.App{
		UserID:     user.ID,
		PlatformID: db.AndroidPlatformID,
		HumanName:  "Test App",
		PackageName: sql.NullString{
			String: fmt.Sprintf("app.package%d", index),
			Valid:  true,
		},
		AppKey: sql.NullString{
			String: fmt.Sprintf("appkey%d", index),
			Valid:  true,
		},
		Settings: make(map[string]any),
	}

	if err := tx.Create(app).Error; err != nil {
		t.Fatalf("Failed to create app: %v", app)
	}
	return app
}

// CreateAppsList is deprecated in favor of CreateList with AppFactory
func CreateAppsList(t *testing.T, tx *db.DB, count int) []*db.App {
	t.Helper()

	apps := make([]*db.App, count)
	for i := range apps {
		apps[i] = CreateApp(t, tx, i, nil)
	}

	return apps
}

type DemandSourceOption func(db *db.DemandSource)

func WithDemandSourceOptions(opt *db.DemandSource) DemandSourceOption {
	return func(demandSource *db.DemandSource) {
		if opt.HumanName != "" {
			demandSource.HumanName = opt.HumanName
		}
		if opt.APIKey != "" {
			demandSource.APIKey = opt.APIKey
		}
	}
}

// CreateDemandSource is deprecated in favor of Create with DemandSourceFactory
func CreateDemandSource(t *testing.T, tx *db.DB, opts ...DemandSourceOption) *db.DemandSource {
	t.Helper()

	demandSource := &db.DemandSource{
		APIKey:    fmt.Sprintf("apikey%d", time.Now().UnixNano()),
		HumanName: "demandsource",
	}

	for _, opt := range opts {
		opt(demandSource)
	}

	if err := tx.Create(demandSource).Error; err != nil {
		t.Fatalf("Failed to create demand source: %v", err)
	}
	return demandSource
}

// CreateDemandSourcesList is deprecated in favor of CreateList with DemandSourceFactory
func CreateDemandSourcesList(t *testing.T, tx *db.DB, count int) []*db.DemandSource {
	t.Helper()

	demandSources := make([]*db.DemandSource, count)
	for i := range demandSources {
		demandSources[i] = CreateDemandSource(t, tx)
	}

	return demandSources
}

type DemandSourceAccountOption func(*db.DemandSourceAccount)

func WithDemandSourceAccountOptions(optAccount *db.DemandSourceAccount) DemandSourceAccountOption {
	return func(account *db.DemandSourceAccount) {
		if optAccount.DemandSourceID != 0 {
			account.DemandSourceID = optAccount.DemandSourceID
		}
		if optAccount.UserID != 0 {
			account.UserID = optAccount.UserID
		}
		if optAccount.Type != "" {
			account.Type = optAccount.Type
		}
		if optAccount.Extra != nil {
			account.Extra = optAccount.Extra
		}
		if optAccount.IsBidding != nil {
			account.IsBidding = optAccount.IsBidding
		}
		if optAccount.IsDefault.Valid {
			account.IsDefault = optAccount.IsDefault
		}
		if optAccount.DemandSource != (db.DemandSource{}) {
			account.DemandSource = optAccount.DemandSource
		}
	}
}

// CreateDemandSourceAccount is deprecated in favor of Create with DemandSourceAccountFactory
func CreateDemandSourceAccount(t *testing.T, tx *db.DB, opts ...DemandSourceAccountOption) *db.DemandSourceAccount {
	t.Helper()

	account := &db.DemandSourceAccount{
		Type:      "DemandSourceAccount::Admob",
		Extra:     []byte(`{}`),
		IsBidding: new(bool),
		IsDefault: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}

	for _, opt := range opts {
		opt(account)
	}

	index := int(time.Now().UnixNano())

	if account.UserID == 0 {
		account.UserID = CreateUser(t, tx, index).ID
	}

	if account.DemandSourceID == 0 {
		account.DemandSourceID = CreateDemandSource(t, tx).ID
	}

	if err := tx.Create(account).Error; err != nil {
		t.Fatalf("Failed to create demand source account: %v", err)
	}
	return account
}

// CreateDemandSourceAccountsList is deprecated in favor of CreateList with DemandSourceAccountFactory
func CreateDemandSourceAccountsList(t *testing.T, tx *db.DB, count int) []*db.DemandSourceAccount {
	t.Helper()

	accounts := make([]*db.DemandSourceAccount, count)
	for i := range accounts {
		accounts[i] = CreateDemandSourceAccount(t, tx)
	}

	return accounts
}

func CreateSegment(t *testing.T, tx *db.DB, index int, app *db.App) *db.Segment {
	t.Helper()

	if app == nil {
		app = CreateApp(t, tx, index, nil)
	}

	segment := &db.Segment{
		Name:        "Test Segment",
		Description: "description",
		Filters:     []segment.Filter{},
		Enabled:     new(bool),
		AppID:       app.ID,
		Priority:    0,
	}

	if err := tx.Create(segment).Error; err != nil {
		t.Fatalf("Failed to create segment: %v", err)
	}
	return segment
}

func CreateSegmentsList(t *testing.T, tx *db.DB, count int) []*db.Segment {
	t.Helper()

	segments := make([]*db.Segment, count)
	for i := range segments {
		segments[i] = CreateSegment(t, tx, i, nil)
	}

	return segments
}

// CreateUser is deprecated in favor of Create with UserFactory
func CreateUser(t *testing.T, tx *db.DB, index int) *db.User {
	t.Helper()

	user := &db.User{
		Email: fmt.Sprintf("test%d@email.com", index),
	}

	if err := tx.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	return user
}

// CreateUsersList is deprecated in favor of CreateList with UserFactory
func CreateUsersList(t *testing.T, tx *db.DB, usersCount int) []*db.User {
	t.Helper()

	users := make([]*db.User, usersCount)
	for i := range users {
		users[i] = CreateUser(t, tx, i)
	}

	return users
}
