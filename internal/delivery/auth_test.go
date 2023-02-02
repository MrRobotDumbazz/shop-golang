package delivery_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"shop/internal/delivery"
	"shop/internal/repository"
	"shop/internal/service"
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestSignUp(t *testing.T) {
	db, err := repository.Init()
	if err != nil {
		t.Errorf("Error opening db: %v\n", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Can't close db err: %v\n", err)
		} else {
			log.Print("db closed")
		}
	}()
	repositories := repository.NewRepository(db)
	redis, _ := redismock.NewClientMock()
	redisrepository := repository.NewRedisRepository(redis)
	services := service.NewServices(repositories, *redisrepository)
	handlers := delivery.NewHandler(services)

	form := url.Values{}
	form.Add("email", "example@gmail.com")
	form.Add("password", "secret123")
	req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Errorf("Error httprequest because: %v\n", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.SignUp)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}
}

func TestSignIn(t *testing.T) {
	db, err := repository.Init()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		} else {
			log.Print("db closed")
		}
	}()
	repositories := repository.NewRepository(db)
	redis, mock := redismock.NewClientMock()
	redisrepository := repository.NewRedisRepository(redis)
	services := service.NewServices(repositories, *redisrepository)
	handlers := delivery.NewHandler(services)
	mux := http.NewServeMux()
	mux.HandleFunc("/signin", handlers.SignIn)
	server := httptest.NewServer(mux)
	defer server.Close()
	body := []byte("email=2232@gmail.com&password=123456")
	res, err := http.Post(server.URL+"/signin", "application/x-www-form-urlencoded", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}
}
