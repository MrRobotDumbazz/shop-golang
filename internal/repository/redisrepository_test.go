package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v9"
)

var client *redis.Client

var (
	key = "key"
	val = "val"
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	code := m.Run()
	os.Exit(code)
}

func TestSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	claimsID := 1
	key := fmt.Sprintf("token-%d", claimsID)
	mock.ExpectSetNX(key, claimsID, 10*time.Minute).SetErr(errors.New("FAIL"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	err := db.SetNX(ctx, key, claimsID, 10*time.Minute)
	err2 := err.Err()
	if err2 == nil || err2.Error() != "FAIL" {
		t.Error(err2)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	claimsID := 1
	key := fmt.Sprintf("token-%d", claimsID)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	mock.ExpectGet(key).RedisNil()
	err := db.Get(ctx, key)
	if err == nil {
		t.Error(err)
	}
}
