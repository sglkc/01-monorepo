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

func TestArticleHappyPath(t *testing.T) {
	t.Parallel()

	mockArticleService := new(mocks.ArticleService)

    newArticle := domain.Article{
	    ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
        Title: "Test judul",
	    Author: "John Doe",
        Content: "Test content",
        Status: "published",
    }

    // data for update
	updatedArticle := newArticle
	updatedArticle.Author = "Jane Doe"
    updatedArticle.Title = "Test judul 2"
	updatedArticle.Content = "Test content update"
    updatedArticle.Status = "published"

	handler := rest.ArticleHandler{
		Service: mockArticleService,
	}

    // --- Create Article
	t.Run("CreateArticle", func(t *testing.T) {
		createReq := domain.CreateArticleRequest{
			Title: newArticle.Title,
			Author: newArticle.Author,
			Content: newArticle.Content,
            Status: newArticle.Status,
		}
		mockArticleService.
			On("CreateArticle", mock.Anything, &createReq).
			Return(&newArticle, nil).
			Once()

		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateArticle(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp domain.ResponseSingleData[domain.Article]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, &newArticle, &resp.Data)

		mockArticleService.AssertExpectations(t)
	})

    // --- Update Article
    t.Run("UpdateArticle", func(t *testing.T) {
        id, err := uuid.Parse(newArticle.ID)
		require.NoError(t, err)

		mockArticleService.
			On("UpdateArticle", mock.Anything, id, mock.MatchedBy(func(u *domain.Article) bool {
                return u.Title == updatedArticle.Title &&
                    u.Author == updatedArticle.Author &&
                    u.Content == updatedArticle.Content &&
                    u.Status == updatedArticle.Status
			})).
			Return(&updatedArticle, nil).
			Once()

		body, err := json.Marshal(updatedArticle)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/"+newArticle.ID, bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newArticle.ID)

		err = handler.UpdateArticle(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.ResponseSingleData[domain.Article]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, &updatedArticle, &resp.Data)

		mockArticleService.AssertExpectations(t)
	})

// 	// --- Get Article
    t.Run("GetArticle", func(t *testing.T) {
        id, err := uuid.Parse(newArticle.ID)
		require.NoError(t, err)

		mockArticleService.
			On("GetArticle", mock.Anything, id).
			Return(&updatedArticle, nil).
			Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+newArticle.ID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newArticle.ID)

        err = handler.GetArticle(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.ResponseSingleData[domain.Article]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, updatedArticle, resp.Data)

		mockArticleService.AssertExpectations(t)
	})

   // --- Delete Article
	t.Run("DeleteArticle", func(t *testing.T) {
        id, err := uuid.Parse(newArticle.ID)
		require.NoError(t, err)
		mockArticleService.
			On("DeleteArticle", mock.Anything, id).
			Return(nil).
			Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/"+newArticle.ID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newArticle.ID)

		err = handler.DeleteArticle(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockArticleService.AssertExpectations(t)
	})

    // --- Verify Deletion
	t.Run("GetArticleAfterDeletion", func(t *testing.T) {
        id, err := uuid.Parse(newArticle.ID)
		require.NoError(t, err)
		mockArticleService.
            On("GetArticle", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

		e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+newArticle.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newArticle.ID)

        err = handler.GetArticle(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "Article not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockArticleService.AssertExpectations(t)
	})
}



func TestArticleUnhappyPath(t *testing.T) {
    mockArticleService := new(mocks.ArticleService)

    newArticle := domain.Article{
	    ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
        Title: "Test judul",
	    Author: "John Doe",
        Content: "Test content",
        Status: "published",
    }

    handler := rest.ArticleHandler{
		Service: mockArticleService,
	}


	// --- Get Non-Existent Article
	t.Run("GetNonExistingArticle", func(t *testing.T) {
        id, err := uuid.Parse(newArticle.ID)
		require.NoError(t, err)
		mockArticleService.
            On("GetArticle", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

		e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+newArticle.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newArticle.ID)

        err = handler.GetArticle(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "Article not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockArticleService.AssertExpectations(t)
	})

	// // --- Create Invalid JSON
	t.Run("CreateArticle_InvalidNameType", func(t *testing.T) {
	    body := []byte(`{
            "Title": "",
            "Author": "",
            "Content": "",
            "Status": "",
	    }`)

	    e := echo.New()
	    req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewReader(body))
	    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	    rec := httptest.NewRecorder()
	    c := e.NewContext(req, rec)

	    err := handler.CreateArticle(c)
	    require.NoError(t, err)

	    assert.Equal(t, http.StatusBadRequest, rec.Code)

	    var resp domain.ResponseSingleData[domain.Empty]
	    err = json.Unmarshal(rec.Body.Bytes(), &resp)
	    require.NoError(t, err)
	    assert.Equal(t, "error", resp.Status)
	    assert.NotEmpty(t, resp.Message)
})
}
