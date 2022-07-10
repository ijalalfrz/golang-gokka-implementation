package wallet_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/ijalalfrz/coinbit-test/wallet/mocks"
	"github.com/ijalalfrz/coinbit-test/webmodel"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	vld *validator.Validate
)

func TestMain(m *testing.M) {
	vld = validator.New()
	m.Run()
}

func TestNewWalletHTTPHandlerConstruct(t *testing.T) {
	t.Run("should construct the wallet http handler", func(t *testing.T) {

		logger := logrus.New()
		usecase := &mocks.Usecase{}
		router := &mux.Router{}
		wallet.NewWalletHTTPHandler(logger, vld, router, usecase)
	})
}

func TestDepositWallet_Error_UnprocessableEntity(t *testing.T) {
	usecase := &mocks.Usecase{}
	hh := wallet.HTTPHandler{
		Logger:   logrus.New(),
		Validate: vld,
		Usecase:  usecase,
	}
	r := httptest.NewRequest(http.MethodPost, "/just/for/testing", strings.NewReader(`should error`))
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(hh.DepositWallet)

	handler.ServeHTTP(recorder, r)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestDepositWallet_Error_BadRequest(t *testing.T) {
	usecase := &mocks.Usecase{}
	hh := wallet.HTTPHandler{
		Logger:   logrus.New(),
		Validate: vld,
		Usecase:  usecase,
	}
	payload, _ := json.Marshal(webmodel.DepositWalletPayload{})
	r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(payload))
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(hh.DepositWallet)

	handler.ServeHTTP(recorder, r)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
