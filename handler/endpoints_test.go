package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/config"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/generated"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/repository"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupTestPostEstate(t *testing.T) (*Server, *repository.MockRepositoryInterface, *echo.Echo) {
	t.Parallel()
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	server := &Server{
		Repository: mockRepo,
		Config: &config.Config{
			ScaleFactor: 10,
		},
	}

	e := echo.New()
	e.Validator = validator.NewRequestValidator()

	e.POST("/estate", server.PostEstate)

	return server, mockRepo, e
}

// TestPostEstate tests the PostEstate handler function.
// It sets up a mock repository, creates a new HTTP request, and calls the
// PostEstate handler. It then checks the response to ensure it
// matches the expected behavior.
func TestPostEstate(t *testing.T) {

	t.Run("Valid request - parameter follows happy path", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstate(t)
		// Create a valid request body
		reqBody := generated.PostEstateJSONRequestBody{
			Length: 10,
			Width:  10,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		// Set up the mock repository response
		mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(&repository.CreateEstateOutput{
			Id: uuid.New().String(),
		}, nil)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err = server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp generated.EstateResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp.Id)
	})

	t.Run("Invalid request body - input as string", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{"length": "abc", "width": 20}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - input as zero", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{"length": "0", "width": 20}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - input as negative numbers", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{"length": -5, "width": -7}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - input out of bound", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{"length": 57000, "width": 80000}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - incomplete parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{"width": 80000}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - nil parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{}`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - random parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{ adjksfboiasdfhu8798 }`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Unexcepted internal server error - create estate", func(t *testing.T) {
		server, _, e := setupTestPostEstate(t)
		// Create an invalid request body
		requestBody := []byte(`{ adjksfboiasdfhu8798 }`)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err := server.PostEstate(c)
		require.NoError(t, err)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Error creating estate", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstate(t)
		// Create a valid request body
		reqBody := generated.PostEstateJSONRequestBody{
			Length: 10,
			Width:  10,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to create estate"))

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the PostEstate handler
		err = server.PostEstate(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Something happens in our end. Let us check.", errResp["error"])
	})
}

// setupTestPostEstateEstateIdTree sets up a test environment for the PostEstateEstateIdTree handler function.
// It creates a new mock repository, a new Server instance, and a new Echo instance with a route for the handler.
// The function returns the Server, mock repository, and Echo instance for use in tests.
func setupTestPostEstateEstateIdTree(t *testing.T) (*Server, *repository.MockRepositoryInterface, *echo.Echo) {
	t.Parallel()
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	server := &Server{
		Repository: mockRepo,
		Config:     &config.Config{},
	}

	e := echo.New()
	e.Validator = validator.NewRequestValidator()

	e.POST("/estate/:estateId/tree", func(c echo.Context) error {
		estateID, err := uuid.Parse(c.Param("estateId"))
		if err != nil {
			return err
		}
		return server.PostEstateEstateIdTree(c, estateID)
	})

	return server, mockRepo, e
}

// TestPostEstateEstateIdTree tests the PostEstateEstateIdTree handler function.
// It sets up a mock repository, creates a new HTTP request, and calls the
// PostEstateEstateIdTree handler. It then checks the response to ensure it
// matches the expected behavior.
func TestPostEstateEstateIdTree(t *testing.T) {

	t.Run("Valid request - parameter follows happy path", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		treeId := uuid.New()
		reqBody := generated.PostEstateEstateIdTreeJSONRequestBody{
			X:      5,
			Y:      30,
			Height: 15,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{
			Id: estateId.String(),
		}).Return(&repository.GetEstateByEstateIdOutput{
			Estate: repository.Estate{
				Length: 20,
				Width:  30,
			},
		}, nil)

		mockRepo.EXPECT().IsTreeExist(gomock.Any(), gomock.Any()).Return(&repository.IsTreeExistOutput{
			IsExist: false,
		}, nil)

		mockRepo.EXPECT().CreateTree(gomock.Any(), gomock.Any()).Return(&repository.CreateTreeOutput{
			Id: treeId.String(),
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp generated.TreeResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp.Id)
	})

	t.Run("Invalid request - input as string", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": "invalid", "y": 10, "height": 15}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request - input as 0", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 0, "y": 0, "height": 0}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request - input out of bound", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 50001, "y": 50001, "height": 8}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request - incomplete parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"height": 15}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request - out of bound height", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 5, "y": 5, "height": 55}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request - estate not found", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		validRequestBody := []byte(`{"x": 5, "y": 5, "height": 15}`)

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{
			Id: estateId.String(),
		}).Return(nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(validRequestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Estate not found", errResp["error"])
	})

	t.Run("Invalid request body - coordinates out of bound", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 200, "y": 100, "height": 15}`)

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{
			Id: estateId.String(),
		}).Return(&repository.GetEstateByEstateIdOutput{
			Estate: repository.Estate{
				Length: 20,
				Width:  30,
			},
		}, nil)

		mockRepo.EXPECT().IsTreeExist(gomock.Any(), gomock.Any()).Return(&repository.IsTreeExistOutput{
			IsExist: false,
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - tree already exists", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 3, "y": 8, "height": 15}`)

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{
			Id: estateId.String(),
		}).Return(&repository.GetEstateByEstateIdOutput{
			Estate: repository.Estate{
				Length: 20,
				Width:  30,
			},
		}, nil)

		mockRepo.EXPECT().IsTreeExist(gomock.Any(), gomock.Any()).Return(&repository.IsTreeExistOutput{
			IsExist: true,
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - negative height", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": 3, "y": 8, "height": -10}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - negative coordinates", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{"x": -3, "y": -8, "height": 8}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - nil parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})

	t.Run("Invalid request body - random parameter", func(t *testing.T) {
		server, _, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()
		requestBody := []byte(`{ asdjkasbdjlkjlk tgtrbh}`)

		req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId.String()+"/tree", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.PostEstateEstateIdTree(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", errResp["error"])
	})
}

// setupTestGetEstateEstateIdDronePlan sets up a test environment for the GetEstateEstateIdDronePlan handler.
// It creates a new Server with a mock repository, an Echo instance, and a route for the handler.
// The function returns the Server, the mock repository, and the Echo instance.
func setupTestGetEstateEstateIdDronePlan(t *testing.T) (*Server, *repository.MockRepositoryInterface, *echo.Echo) {
	t.Parallel()
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	server := &Server{
		Repository: mockRepo,
		Config: &config.Config{
			ScaleFactor: 10,
		},
	}

	e := echo.New()
	e.Validator = validator.NewRequestValidator()

	e.POST("/estate/:estateId/drone-plan", func(c echo.Context) error {
		estateID, err := uuid.Parse(c.Param("estateId"))
		if err != nil {
			return err
		}
		return server.GetEstateEstateIdDronePlan(c, estateID, generated.GetEstateEstateIdDronePlanParams{})
	})

	return server, mockRepo, e
}

// TestGetEstateEstateIdDronePlan tests the GetEstateEstateIdDronePlan handler function.
// It sets up a mock repository and tests the handler function with valid and invalid requests.
// The test cases cover the following scenarios:
// - Valid request with a happy path parameter
// - Valid request with a different happy path parameter
// - Estate not found
func TestGetEstateEstateIdDronePlan(t *testing.T) {

	t.Run("Valid request #1 - parameter follows happy path", func(t *testing.T) {
		server, mockRepo, e := setupTestGetEstateEstateIdDronePlan(t)
		estateId := uuid.New()

		mockRepo.EXPECT().GetEstateTreesByEstateId(gomock.Any(), &repository.GetEstateTreesByEstateIdInput{
			EstateId: estateId.String(),
		}).Return(&repository.GetEstateTreesByEstateIdOutput{
			Estate: repository.Estate{
				Length: 5,
				Width:  1,
			},
			Trees: []repository.Tree{
				{X: 1, Y: 1, Height: 5},
				{X: 2, Y: 1, Height: 2},
				{X: 3, Y: 1, Height: 1},
				{X: 4, Y: 1, Height: 5},
				{X: 5, Y: 1, Height: 3},
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId.String()+"/drone-plan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.GetEstateEstateIdDronePlan(c, estateId, generated.GetEstateEstateIdDronePlanParams{})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp generated.DronePlanResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, resp.Distance, 60)
	})

	t.Run("Valid request #2 - parameter follows happy path", func(t *testing.T) {
		server, mockRepo, e := setupTestGetEstateEstateIdDronePlan(t)
		estateId := uuid.New()

		mockRepo.EXPECT().GetEstateTreesByEstateId(gomock.Any(), &repository.GetEstateTreesByEstateIdInput{
			EstateId: estateId.String(),
		}).Return(&repository.GetEstateTreesByEstateIdOutput{
			Estate: repository.Estate{
				Length: 5,
				Width:  1,
			},
			Trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId.String()+"/drone-plan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.GetEstateEstateIdDronePlan(c, estateId, generated.GetEstateEstateIdDronePlanParams{})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp generated.DronePlanResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, resp.Distance, 54)
	})

	t.Run("Estate not found", func(t *testing.T) {
		server, mockRepo, e := setupTestPostEstateEstateIdTree(t)
		estateId := uuid.New()

		mockRepo.EXPECT().GetEstateTreesByEstateId(gomock.Any(), &repository.GetEstateTreesByEstateIdInput{
			EstateId: estateId.String(),
		}).Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId.String()+"/drone-plan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.GetEstateEstateIdDronePlan(c, estateId, generated.GetEstateEstateIdDronePlanParams{})
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Estate not found", errResp["error"])
	})

}

// TestGetEstateEstateIdStats tests the GetEstateEstateIdStats handler function.
// It checks the happy path scenario where the estate and its stats are successfully retrieved,
// as well as the error scenario where the estate is not found.
func TestGetEstateEstateIdStats(t *testing.T) {
	t.Parallel()

	t.Run("Valid request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		server := &Server{Repository: mockRepo, Config: &config.Config{}}
		estateId := uuid.New()

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{Id: estateId.String()}).Return(&repository.GetEstateByEstateIdOutput{
			Estate: repository.Estate{Length: 10, Width: 20},
		}, nil)
		mockRepo.EXPECT().GetEstateStatsByEstateId(gomock.Any(), &repository.GetEstateStatsByEstateIdInput{EstateId: estateId.String()}).Return(&repository.GetEstateStatsByEstateIdOutput{
			Count:  5,
			Max:    15,
			Min:    5,
			Median: 10,
		}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId.String()+"/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.GetEstateEstateIdStats(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp generated.EstateStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, resp.Count, 5)
		assert.Equal(t, resp.Max, 15)
		assert.Equal(t, resp.Min, 5)
		assert.Equal(t, resp.Median, float32(10))
	})

	t.Run("Estate not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		server := &Server{Repository: mockRepo, Config: &config.Config{}}
		estateId := uuid.New()

		mockRepo.EXPECT().GetEstateByEstateId(gomock.Any(), &repository.GetEstateByEstateIdInput{Id: estateId.String()}).Return(nil, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId.String()+"/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := server.GetEstateEstateIdStats(c, estateId)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var errResp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "Estate not found", errResp["error"])
	})

}

func TestCalculateDroneDistance(t *testing.T) {
	t.Parallel()

	t.Run("Invalid input - nil", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(nil, nil)
		assert.Nil(t, calculateDroneDistanceOutput)
		assert.EqualError(t, err, "err CalculateDroneDistance: invalid input -- nothing to calculate drone distance")
	})

	t.Run("Valid input - single tree - parameter follows happy path", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 5},
			Trees:  []repository.Tree{{X: 3, Y: 3, Height: 5}},
		}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, nil)
		assert.NoError(t, err)
		assert.Equal(t, 252, calculateDroneDistanceOutput.TotalDistance)
	})

	t.Run("Valid input - single estate single tree - parameter follows happy path", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 1, Width: 1},
			Trees:  []repository.Tree{{X: 1, Y: 1, Height: 5}},
		}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, nil)
		assert.NoError(t, err)
		assert.Equal(t, 12, calculateDroneDistanceOutput.TotalDistance)
	})

	t.Run("Valid input - multiple trees - parameter follows happy path", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 5},
			Trees:  []repository.Tree{{X: 2, Y: 2, Height: 5}, {X: 3, Y: 3, Height: 3}, {X: 4, Y: 4, Height: 4}},
		}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, nil)
		assert.NoError(t, err)
		assert.Equal(t, 266, calculateDroneDistanceOutput.TotalDistance)
	})

	t.Run("Valid input - no trees - parameter follows happy path", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 5},
			Trees:  []repository.Tree{},
		}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, nil)
		assert.NoError(t, err)
		assert.Equal(t, 242, calculateDroneDistanceOutput.TotalDistance)
	})

	t.Run("Valid input - large estate - parameter follows happy path", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 100, Width: 100},
			Trees:  []repository.Tree{{X: 50, Y: 50, Height: 10}, {X: 75, Y: 25, Height: 20}},
		}
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, nil)
		assert.NoError(t, err)
		assert.Equal(t, 100052, calculateDroneDistanceOutput.TotalDistance)
	})

	t.Run("Test coordinates and travelled distance #1", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 1},
			Trees:  []repository.Tree{{X: 5, Y: 1, Height: 5}},
		}

		expectedValue := 46
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, &expectedValue)
		assert.NoError(t, err)
		assert.Equal(t, 4, calculateDroneDistanceOutput.LastAchievableXCoordinate)
		assert.Equal(t, 1, calculateDroneDistanceOutput.LastAchievableYCoordinate)
	})

	t.Run("Test coordinates and travelled distance #2", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 1},
			Trees:  []repository.Tree{{X: 5, Y: 1, Height: 5}},
		}

		expectedValue := 32
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, &expectedValue)
		assert.NoError(t, err)
		assert.Equal(t, 4, calculateDroneDistanceOutput.LastAchievableXCoordinate)
		assert.Equal(t, 1, calculateDroneDistanceOutput.LastAchievableYCoordinate)
	})

	t.Run("Test coordinates and travelled distance #3", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 1},
			Trees:  []repository.Tree{{X: 5, Y: 1, Height: 5}},
		}
		expectedValue := 33
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, &expectedValue)
		assert.NoError(t, err)
		assert.Equal(t, 4, calculateDroneDistanceOutput.LastAchievableXCoordinate)
		assert.Equal(t, 1, calculateDroneDistanceOutput.LastAchievableYCoordinate)
	})

	t.Run("Test coordinates and travelled distance #4", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 2},
			Trees:  []repository.Tree{{X: 5, Y: 1, Height: 5}, {X: 5, Y: 2, Height: 10}},
		}
		expectedValue := 111
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, &expectedValue)
		assert.NoError(t, err)
		assert.Equal(t, 2, calculateDroneDistanceOutput.LastAchievableXCoordinate)
		assert.Equal(t, 2, calculateDroneDistanceOutput.LastAchievableYCoordinate)
	})

	t.Run("Test coordinates and travelled distance #5", func(t *testing.T) {
		server := &Server{Config: &config.Config{ScaleFactor: 10}}
		input := &repository.CalculateDroneDistanceInput{
			Estate: repository.Estate{Length: 5, Width: 2},
			Trees:  []repository.Tree{{X: 5, Y: 1, Height: 5}, {X: 5, Y: 2, Height: 10}},
		}
		expectedValue := 112
		calculateDroneDistanceOutput, err := server.CalculateDroneDistance(input, &expectedValue)
		assert.NoError(t, err)
		assert.Equal(t, 1, calculateDroneDistanceOutput.LastAchievableXCoordinate)
		assert.Equal(t, 2, calculateDroneDistanceOutput.LastAchievableYCoordinate)
	})

}
