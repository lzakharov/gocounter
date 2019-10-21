package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounter_Process(t *testing.T) {
	t.Run("valid url", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Go is an open source programming language...")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := 1

		counter := Counter{
			Query: []byte("Go"),
		}

		actual, err := counter.Process(server.URL)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid url", func(t *testing.T) {
		counter := Counter{
			Query: []byte("Go"),
		}

		_, err := counter.Process("invalid")
		assert.Error(t, err)
	})
}
