package gin

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"playground/internal/app"
	"playground/internal/delivery/gin/helper"
	mock_app "playground/internal/mock/app"
	"playground/internal/pkg/apperr"
	"playground/internal/pkg/token"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)

	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Manager)
		buildStubs    func(ar *mock_app.MockAccountRepository)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},

			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(&app.Account{}, failure.Translate(sql.ErrNoRows, apperr.NotFound))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(&app.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			i := GetInjector().Clone()
			defer i.Shutdown()
			ctrl := gomock.NewController(t)
			mrm := mock_app.NewMockRepositoryManager(ctrl)
			ar := mock_app.NewMockAccountRepository(ctrl)
			mrm.EXPECT().Account().AnyTimes().Return(ar)
			tc.buildStubs(ar)

			do.OverrideValue[app.RepositoryManager](i, mrm)
			tm := do.MustInvoke[token.Manager](i)
			router := do.MustInvoke[*gin.Engine](i)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, tm)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Manager)
		buildStubs    func(ar *mock_app.MockAccountRepository)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				arg := &app.Account{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				ar.EXPECT().Create(gomock.Any(), gomock.Eq(arg)).Times(1).Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(&app.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			i := GetInjector().Clone()
			defer i.Shutdown()
			ctrl := gomock.NewController(t)
			mrm := mock_app.NewMockRepositoryManager(ctrl)
			ar := mock_app.NewMockAccountRepository(ctrl)
			mrm.EXPECT().Account().AnyTimes().Return(ar)
			tc.buildStubs(ar)

			do.OverrideValue[app.RepositoryManager](i, mrm)
			tm := do.MustInvoke[token.Manager](i)
			router := do.MustInvoke[*gin.Engine](i)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, tm)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts := make([]app.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = *randomAccount(user.Username)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Manager)
		buildStubs    func(ar *mock_app.MockAccountRepository)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				arg := &app.ListAccountsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}
				ar.EXPECT().List(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().List(gomock.Any(), gomock.Any()).Times(1).Return([]app.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Manager) {
				addAuthorization(t, request, tokenMaker, helper.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(ar *mock_app.MockAccountRepository) {
				ar.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			i := GetInjector().Clone()
			defer i.Shutdown()
			ctrl := gomock.NewController(t)
			mrm := mock_app.NewMockRepositoryManager(ctrl)
			ar := mock_app.NewMockAccountRepository(ctrl)
			mrm.EXPECT().Account().AnyTimes().Return(ar)
			tc.buildStubs(ar)

			do.OverrideValue[app.RepositoryManager](i, mrm)
			tm := do.MustInvoke[token.Manager](i)
			router := do.MustInvoke[*gin.Engine](i)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, tm)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAccount(owner string) *app.Account {
	return &app.Account{
		ID:       app.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  app.RandomMoney(),
		Currency: app.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account *app.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount *app.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []app.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []app.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
