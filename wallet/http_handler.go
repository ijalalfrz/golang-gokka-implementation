package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/ijalalfrz/coinbit-test/webmodel"
	"github.com/sirupsen/logrus"
)

const (
	basePath = "/wallet"
)

// HTTPHandler is a concrete struct of wallet http handler.
type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

func NewWalletHTTPHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, usecase Usecase) {
	handler := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}
	router.HandleFunc(basePath+"/v1/deposit", handler.DepositWallet).Methods(http.MethodPost)
	router.HandleFunc(basePath+"/v1/details/{walletId}", handler.GetDetailWallet).Methods(http.MethodGet)

}

// GetDetailWallet is a function  handle get detail wallet
func (handler HTTPHandler) GetDetailWallet(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	ctx := r.Context()

	pathVariables := mux.Vars(r)
	walletId := pathVariables["walletId"]

	resp = handler.Usecase.GetDetail(ctx, walletId)
	response.JSON(w, resp)
	return
}

// DepositWallet is a function to handle deposit request
func (handler HTTPHandler) DepositWallet(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var payload webmodel.DepositWalletPayload

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	if err := handler.validateRequestBody(payload); err != nil {
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = handler.Usecase.Deposit(ctx, payload)
	response.JSON(w, resp)
	return
}

// validateRequestBody will validate payload to be processed
func (handler HTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = handler.Validate.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("Invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
