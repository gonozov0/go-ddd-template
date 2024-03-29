package echomiddleware_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go-echo-ddd-template/pkg/echomiddleware"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	m.Called(ctx, level, msg, attrs)
}

func TestSlogLoggerMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(echomiddleware.RequestIDHeader, "some-id")
	req.Header.Add(echomiddleware.TraceParentHeader, "aa-bb-cc-dd")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger := new(MockLogger)
	mockLogger.On("LogAttrs", mock.Anything, slog.LevelInfo, "REQUEST", mock.Anything)

	m := echomiddleware.SlogLoggerMiddleware(mockLogger)
	h := m(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
	err := h(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	mockLogger.AssertExpectations(t)
}
