// Package dbtest provides helper functions for tests that require database access
package dbtest

import (
	"hash/fnv"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bwmarrin/snowflake"
	"github.com/joho/godotenv"
)

// Prepare sets up the database for testing and initializes test factories.
// It should be called once per test package, usually in TestMain.
func Prepare() *db.DB {
	var (
		testDB *db.DB
		err    error
	)

	// Calculate a number from the file path of the calling function
	// to use as a snowflake node ID and base for factory counters.
	// Because running tests from multiple packages in parallel can cause transaction deadlocks.
	var hashNum uint32
	if _, file, _, ok := runtime.Caller(1); ok {
		hash := fnv.New32()
		_, _ = hash.Write([]byte(file))
		hashNum = hash.Sum32()
	}

	counter.base = hashNum >> 16 // cut in half for readability

	err = godotenv.Load("../../../.env.test")
	if err != nil {
		log.Printf("Did not load from .env.test file: %v", err)
	}

	nodeID := int64(hashNum >> 22) // max value is 1023 or 10 bits
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Error creating snowflake node: %v", err)
	}

	testDB, err = db.Open(
		os.Getenv("DATABASE_URL"),
		db.WithIgnoreRecordNotFoundError(),
		db.WithSnowflakeNode(node),
	)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	return testDB
}

var counter = newCounters()

type counters struct {
	mu       sync.Mutex
	counters map[string]uint32
	base     uint32
}

func newCounters() *counters {
	return &counters{
		counters: make(map[string]uint32),
	}
}

func (c *counters) get(name string) uint32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.counters[name]; ok {
		c.counters[name]++
	} else {
		c.counters[name] = c.base
	}

	return c.counters[name]
}
