package helpers_test

import (
	"api-challenges/internal/helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		status  int
		data    any
		wantErr bool
	}{
		{
			name:    "valid input",
			status:  http.StatusOK,
			data:    map[string]string{"message": "success"},
			wantErr: false,
		},
		{
			name:    "invalid input",
			status:  http.StatusNotFound,
			data:    make(chan int), // channels cannot be JSON encoded
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			err := helpers.WriteJson(rec, tt.status, tt.data)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.status, rec.Code)
				assert.JSONEq(t, `{"message":"success"}`, rec.Body.String())
			}
		})
	}
}
