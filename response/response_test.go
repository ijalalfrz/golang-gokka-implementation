package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ijalalfrz/coinbit-test/exception"
	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	resp := response.NewErrorResponse(
		exception.ErrNotFound, http.StatusNotFound, nil, response.StatNotFound, "Resource not found",
	)

	assert.NotNil(t, resp)
	assert.Equal(t, exception.ErrNotFound, resp.Error())
	assert.Equal(t, http.StatusNotFound, resp.HTTPStatusCode())
	assert.Nil(t, resp.Data())
	assert.Equal(t, response.StatNotFound, resp.Status())
	assert.Equal(t, "Resource not found", resp.Message())
}

func TestSuccessResponse(t *testing.T) {
	t.Run("when status is common ok", func(t *testing.T) {
		resp := response.NewSuccessResponse(
			nil, response.StatOK, "OK",
		)

		assert.NotNil(t, resp)
		assert.Nil(t, resp.Error())
		assert.Equal(t, http.StatusOK, resp.HTTPStatusCode())
		assert.Nil(t, resp.Data())
		assert.Equal(t, response.StatOK, resp.Status())
		assert.Equal(t, "OK", resp.Message())
	})

	t.Run("when status is created", func(t *testing.T) {
		resp := response.NewSuccessResponse(
			nil, response.StatCreated, "Created",
		)

		assert.NotNil(t, resp)
		assert.Nil(t, resp.Error())
		assert.Equal(t, http.StatusCreated, resp.HTTPStatusCode())
		assert.Nil(t, resp.Data())
		assert.Equal(t, response.StatCreated, resp.Status())
		assert.Equal(t, "Created", resp.Message())
	})
}

func TestRESTResponse(t *testing.T) {
	t.Run("responding json as success", func(t *testing.T) {
		recoreder := httptest.NewRecorder()
		resp := response.NewSuccessResponse(
			nil, response.StatOK, "OK",
		)
		response.JSON(recoreder, resp)

		assert.Equal(t, http.StatusOK, recoreder.Code)
	})

	t.Run("responding json as error", func(t *testing.T) {
		recoreder := httptest.NewRecorder()
		resp := response.NewErrorResponse(
			exception.ErrNotFound, http.StatusNotFound, nil, response.StatNotFound, "Resource not found",
		)
		response.JSON(recoreder, resp)

		assert.Equal(t, http.StatusNotFound, recoreder.Code)
	})
}
