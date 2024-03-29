package gin

import (
	"bytes"
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

	"playground/internal/app"
	mock_app "playground/internal/mock/app"
	"playground/internal/pkg/apperr"
	"playground/internal/pkg/password"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher)
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
			buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
				arg := &app.User{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				mrm.User().(*mock_app.MockUserRepository).EXPECT().
					Create(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
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
			buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
				mrm.User().(*mock_app.MockUserRepository).EXPECT().
					Create(gomock.Any(), gomock.Any()).Times(1).Return(&app.User{}, sql.ErrConnDone)
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
		// 	buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
		// 		mrm.User().(*mock_app.MockUserRepository).EXPECT().
		// 			Create(gomock.Any(), gomock.Any()).Times(1).Return(&app.User{}, db.ErrUniqueViolation)
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
			buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
				mrm.User().(*mock_app.MockUserRepository).EXPECT().
					Create(gomock.Any(), gomock.Any()).Times(0)
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
			buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
				mrm.User().(*mock_app.MockUserRepository).EXPECT().
					Create(gomock.Any(), gomock.Any()).Times(0)
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
			buildStubs: func(mrm *mock_app.MockRepository, md *mock_app.MockDispatcher) {
				mrm.User().(*mock_app.MockUserRepository).EXPECT().
					Create(gomock.Any(), gomock.Any()).Times(0)
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
			mrm := NewMockRepository(t, ctrl)
			md := mock_app.NewMockDispatcher(ctrl)
			tc.buildStubs(mrm, md)

			do.OverrideValue[app.Repository](i, mrm)
			do.OverrideValue[app.Dispatcher](i, md)
			router := do.MustInvoke[*gin.Engine](i)
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
		buildStubs    func(ur *mock_app.MockUserRepository)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(ur *mock_app.MockUserRepository) {
				ur.EXPECT().Get(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				ur.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)
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
			buildStubs: func(ur *mock_app.MockUserRepository) {
				ur.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(nil, failure.Translate(sql.ErrNoRows, apperr.NotFound))
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
			buildStubs: func(ur *mock_app.MockUserRepository) {
				ur.EXPECT().Get(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
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
			buildStubs: func(ur *mock_app.MockUserRepository) {
				ur.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(&app.User{}, sql.ErrConnDone)
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
			buildStubs: func(ur *mock_app.MockUserRepository) {
				ur.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
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
			mrm := NewMockRepository(t, ctrl)
			ur := mrm.User().(*mock_app.MockUserRepository)
			tc.buildStubs(ur)

			do.OverrideValue[app.Repository](i, mrm)
			router := do.MustInvoke[*gin.Engine](i)
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

func randomUser(t *testing.T) (*app.User, string) {
	pwd := app.RandomString(6)
	hashedPassword, err := password.Hash(pwd)
	require.NoError(t, err)

	u := app.NewUser(app.RandomOwner(), hashedPassword, app.RandomOwner(), app.RandomEmail())
	return u, pwd
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user *app.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser *app.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}

func EqCreateUserParams(arg *app.User, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

type eqCreateUserParamsMatcher struct {
	arg      *app.User
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(*app.User)
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
