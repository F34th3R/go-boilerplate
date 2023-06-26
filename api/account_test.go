package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/F34th3R/go_simplebank/db/mock"
	db "github.com/F34th3R/go_simplebank/db/sqlc"
	"github.com/F34th3R/go_simplebank/db/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// Start test server and send request
	server := NewServer(store)
	record := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%s", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(record, request)
	// TODO: fix this test
	// the test fails because the response is not what we expect 500 != 200
	// the error is on url := fmt.Sprintf("/accounts/%s", account.ID)
	// the account.ID is not being converted to a string
	// the error is: panic: interface conversion: interface {} is uuid.UUID, not string
	require.Equal(t, http.StatusOK, record.Code)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomUUID(),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
