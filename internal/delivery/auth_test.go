package delivery_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"shop/internal/delivery"
	"shop/internal/repository"
	"shop/internal/service"
	mock_service "shop/internal/service/mocks"
	"shop/models"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, seller models.Seller)
	testTable := []struct {
		name                string
		inputBody           map[string]string
		inputSeller         models.Seller
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: map[string]string{
				"email":    "2232@gmail.com",
				"password": "123456",
			},
			inputSeller: models.Seller{
				Email:    "2232@gmail.com",
				Password: "123456",
			},
			mockBehavior: func(s *mock_service.MockAuth, seller models.Seller) {
				s.EXPECT().CreateSeller(seller).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(auth, testCase.inputSeller)
			services := service.Service{Auth: auth}
			handlers := delivery.NewHandler(&services)
			handler := http.HandlerFunc(handlers.SignUp)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", nil)
			req.PostFormValue(testCase.inputBody["email"])
			req.PostFormValue(testCase.inputBody["password"])
			handler.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
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
