package repository_test

import (
	"shop/internal/repository"
	"shop/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateSeller(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	r := repository.NewRepository(db)
	tests := []struct {
		name    string
		mock    func()
		input   models.Seller
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO shopdb.sellers").WithArgs("2232@gmail.com", "123456").WillReturnResult(sqlmock.NewResult(1, 0))
			},
			input: models.Seller{
				Email:    "2232@gmail.com",
				Password: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := r.CreateSeller(&tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
