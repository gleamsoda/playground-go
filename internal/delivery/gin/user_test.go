package gin

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"playground/internal/delivery/gin/handler"
	"playground/internal/delivery/gin/middleware"
	"playground/internal/pkg/apperr"
	"playground/internal/pkg/password"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
	mock_wallet "playground/test/mock/wallet"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher) {
				arg := &wallet.User{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				mr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn wallet.TransactionFunc) error {
					return fn(ctx, mr)
				})
				mr.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
				md.EXPECT().SendVerifyEmail(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher) {
				mr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn wallet.TransactionFunc) error {
					return fn(ctx, mr)
				})
				mr.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(&wallet.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// {
		// 	name: "DuplicateUsername",
		// 	body: gin.H{
		// 		"username":  user.Username,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mock_wallet.MockRepository) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(&wallet.User{}, db.ErrUniqueViolation)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusForbidden, recorder.Code)
		// 	},
		// },
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher) {
				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "invalid-email",
			},
			buildStubs: func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher) {
				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mr *mock_wallet.MockRepository, md *mock_wallet.MockDispatcher) {
				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
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
			mr := mock_wallet.NewMockRepository(ctrl)
			md := mock_wallet.NewMockDispatcher(ctrl)
			tc.buildStubs(mr, md)

			do.OverrideValue[wallet.Repository](i, mr)
			do.OverrideValue[wallet.Dispatcher](i, md)
			h := do.MustInvoke[*handler.Handler](i)
			tm := do.MustInvoke[token.Manager](i)
			router := NewRouter(h, middleware.Auth(tm))
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_wallet.MockRepository)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_wallet.MockRepository) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mock_wallet.MockRepository) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, failure.Translate(sql.ErrNoRows, apperr.NotFound))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(store *mock_wallet.MockRepository) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_wallet.MockRepository) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&wallet.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(store *mock_wallet.MockRepository) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
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
			mrepo := mock_wallet.NewMockRepository(ctrl)
			tc.buildStubs(mrepo)

			do.OverrideValue[wallet.Repository](i, mrepo)
			h := do.MustInvoke[*handler.Handler](i)
			tm := do.MustInvoke[token.Manager](i)
			router := NewRouter(h, middleware.Auth(tm))
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (*wallet.User, string) {
	pwd := wallet.RandomString(6)
	hashedPassword, err := password.Hash(pwd)
	require.NoError(t, err)

	u := wallet.NewUser(wallet.RandomOwner(), hashedPassword, wallet.RandomOwner(), wallet.RandomEmail())
	return u, pwd
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user *wallet.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser *wallet.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}

func EqCreateUserParams(arg *wallet.User, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

type eqCreateUserParamsMatcher struct {
	arg      *wallet.User
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(*wallet.User)
	if !ok {
		return false
	}

	err := password.Verify(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}
