package rest_test

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "zog-news/domain"
    "zog-news/internal/rest"
    "zog-news/internal/rest/mocks"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestTopicHappyPath(t *testing.T) {
    t.Parallel()

    mockTopicService := new(mocks.TopicService)

    newTopic := domain.Topic{
        ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
        Name:  "Technology",
    }

    // data for update
    updatedTopic := newTopic
    updatedTopic.Name = "Social"

    handler := rest.TopicHandler{
        Service: mockTopicService,
    }

    // --- Create Topic
    t.Run("CreateTopic", func(t *testing.T) {
        createReq := domain.CreateTopicRequest{
            Name: newTopic.Name,
        }
        mockTopicService.
            On("CreateTopic", mock.Anything, &createReq).
            Return(&newTopic, nil).
            Once()

        body, err := json.Marshal(createReq)
        require.NoError(t, err)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err = handler.CreateTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusCreated, rec.Code)

        var resp domain.ResponseSingleData[domain.Topic]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)
        assert.Equal(t, "success", resp.Status)
        assert.Equal(t, &newTopic, &resp.Data)

        mockTopicService.AssertExpectations(t)
    })

    // --- Update Topic
    t.Run("UpdateTopic", func(t *testing.T) {
        id, err := uuid.Parse(newTopic.ID)
        require.NoError(t, err)

        mockTopicService.
            On("UpdateTopic", mock.Anything, id, mock.MatchedBy(func(u *domain.Topic) bool {
                return u.Name == updatedTopic.Name
            })).
            Return(&updatedTopic, nil).
            Once()

        body, err := json.Marshal(updatedTopic)
        require.NoError(t, err)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPut, "/api/v1/topics/"+newTopic.ID, bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newTopic.ID)

        err = handler.UpdateTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusOK, rec.Code)

        var resp domain.ResponseSingleData[domain.Topic]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)
        assert.Equal(t, "success", resp.Status)
        assert.Equal(t, &updatedTopic, &resp.Data)

        mockTopicService.AssertExpectations(t)
    })

    // 	// --- Get Topic
    t.Run("GetTopic", func(t *testing.T) {
        id, err := uuid.Parse(newTopic.ID)
        require.NoError(t, err)

        mockTopicService.
            On("GetTopic", mock.Anything, id).
            Return(&updatedTopic, nil).
            Once()

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/"+newTopic.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newTopic.ID)

        err = handler.GetTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusOK, rec.Code)

        var resp domain.ResponseSingleData[domain.Topic]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)
        assert.Equal(t, "success", resp.Status)
        assert.Equal(t, updatedTopic, resp.Data)

        mockTopicService.AssertExpectations(t)
    })

    // --- Delete Topic
    t.Run("DeleteTopic", func(t *testing.T) {
        id, err := uuid.Parse(newTopic.ID)
        require.NoError(t, err)
        mockTopicService.
            On("DeleteTopic", mock.Anything, id).
            Return(nil).
            Once()

        e := echo.New()
        req := httptest.NewRequest(http.MethodDelete, "/api/v1/topics/"+newTopic.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newTopic.ID)

        err = handler.DeleteTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNoContent, rec.Code)
        mockTopicService.AssertExpectations(t)
    })

    // --- Verify Deletion
    t.Run("GetTopicAfterDeletion", func(t *testing.T) {
        id, err := uuid.Parse(newTopic.ID)
        require.NoError(t, err)
        mockTopicService.
            On("GetTopic", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/"+newTopic.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newTopic.ID)

        err = handler.GetTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "Topic not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockTopicService.AssertExpectations(t)
    })
}



func TestTopicUnhappyPath(t *testing.T) {
    mockTopicService := new(mocks.TopicService)

    newTopic := domain.Topic{
        ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
        Name:  "Technology",
    }

    handler := rest.TopicHandler{
        Service: mockTopicService,
    }


    // --- Get Non-Existent Topic
    t.Run("GetNonExistingTopic", func(t *testing.T) {
        id, err := uuid.Parse(newTopic.ID)
        require.NoError(t, err)
        mockTopicService.
            On("GetTopic", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/"+newTopic.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newTopic.ID)

        err = handler.GetTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "Topic not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockTopicService.AssertExpectations(t)
    })

    // // --- Create Invalid JSON
    t.Run("CreateTopic_InvalidNameType", func(t *testing.T) {
        body := []byte(`{
            "Name": "",
        }`)

        e := echo.New()
        req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := handler.CreateTopic(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusBadRequest, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)
        assert.Equal(t, "error", resp.Status)
        assert.NotEmpty(t, resp.Message)
    })
}
