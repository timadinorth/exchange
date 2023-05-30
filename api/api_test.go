package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/timadinorth/bet-exchange/model"
)

type ApiTestSuite struct {
	suite.Suite
	server *Server
}

func (t *ApiTestSuite) makeRequest(method, url string, body interface{}) *http.Response {
	requestBody, _ := json.Marshal(body)
	rq, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	rq.Header.Add("Content-Type", "application/json")
	resp, err := t.server.Web.Test(rq, -1)
	assert.Nil(t.T(), err, nil)
	return resp
}

func parseResponse(t *testing.T, resp *http.Response) (response map[string]map[string]string) {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Error(err)
	}
	return

}

func (ts *ApiTestSuite) TestSignUp() {
	newUser := SignUpReq{
		Username: "tim",
		Password: "adi",
	}
	badReq := SignUpReq{
		Username: "bad",
	}
	ts.T().Run("should not allow malformed request", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signup", "test")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	ts.T().Run("user should be able to signup", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signup", newUser)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dbUser := model.User{}
		if err := ts.server.DB.Model(model.User{}).Where("username = ?", newUser.Username).Take(&dbUser).Error; err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		response := parseResponse(t, resp)
		data, exists := response["data"]
		assert.Equal(t, true, exists)
		_, exists = data["username"]
		assert.Equal(t, true, exists)
		assert.Equal(t, data["username"], newUser.Username)
	})

	ts.T().Run("duplicated username should not be allowed", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signup", newUser)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	ts.T().Run("password should be provided in request", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signup", badReq)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func (ts *ApiTestSuite) TestSignIn() {
	newUser := SignUpReq{
		Username: "tim",
		Password: "adi",
	}

	ts.T().Run("should not allow malformed request", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signin", "test")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	ts.T().Run("user should be able to signin after signup", func(t *testing.T) {
		resp := ts.makeRequest("POST", "/api/v1/auth/signup", newUser)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		resp = ts.makeRequest("POST", "/api/v1/auth/signin", newUser)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	ts.T().Run("signin without signup should not work", func(t *testing.T) {
		nonExistentUser := SignUpReq{
			Username: "tim2",
			Password: "adi",
		}

		resp := ts.makeRequest("POST", "/api/v1/auth/signin", nonExistentUser)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	ts.T().Run("signin with wrong password should not work", func(t *testing.T) {
		userWithWrongPassword := SignUpReq{
			Username: "tim",
			Password: "adi2",
		}

		resp := ts.makeRequest("POST", "/api/v1/auth/signin", userWithWrongPassword)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, &ApiTestSuite{})
}

func (t *ApiTestSuite) SetupSuite() {
	t.server = &Server{}
	t.server.InitLogger()
	t.server.LoadConfig("../")
	t.server.ConnectDB()
	t.server.ConnectCache()
	t.server.InitWeb()
	t.server.RegisterRoutes()
}

func (t *ApiTestSuite) SetupTest() {
	err := t.server.SetupModels()
	if err != nil {
		t.T().Errorf("test setup failed: %v", err)
	}
}

func (t *ApiTestSuite) TearDownTest() {
	err := t.server.CleanupModels()
	if err != nil {
		t.T().Errorf("test cleanup failed: %v", err)
	}
}
