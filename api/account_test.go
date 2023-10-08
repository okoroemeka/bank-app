package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdbb "github.com/okoroemeka/simple_bank/db/mock"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := generateRandomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdbb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:      "NotFound",
		accountID: account.ID,
		buildStubs: func(store *mockdbb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
		// Add more cases
	}, {
		name:      "OK",
		accountID: account.ID,
		buildStubs: func(store *mockdbb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)

			requireBodyMatchAccount(t, recorder.Body, account)
		},
	}, {
		name:      "InternalError",
		accountID: account.ID,
		buildStubs: func(store *mockdbb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	}, {
		name:      "BadRequest",
		accountID: 0,
		buildStubs: func(store *mockdbb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdbb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			testCase.checkResponse(t, recorder)
		})
	}

}

//func TestCreateAccount(t *testing.T) {
//	account := generateRandomAccount()
//
//	testCases := []struct {
//		name          string
//		buildStubs    func(store *mockdbb.MockStore)
//		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
//	}{{
//		name: "OK",
//		buildStubs: func(store *mockdbb.MockStore) {
//			store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
//		},
//		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//			require.Equal(t, http.StatusCreated, recorder.Code)
//		},
//	}}
//}

func generateRandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
