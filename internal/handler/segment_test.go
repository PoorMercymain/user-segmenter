package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	"github.com/PoorMercymain/user-segmenter/internal/domain/mocks"
	"github.com/PoorMercymain/user-segmenter/internal/middleware"
	"github.com/PoorMercymain/user-segmenter/internal/service"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

func testRouter(t *testing.T) *echo.Echo {
	e := echo.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	mockRepo.EXPECT().CreateSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockRepo.EXPECT().DeleteSegment(gomock.Any(), gomock.Any()).Return(errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().DeleteSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(nil, errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(make([]string, 0), nil).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(nil, errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil).AnyTimes()

	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.HistoryElem, 0), nil).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.HistoryElem{{UserID: "1", Slug: "a", Operation: "addition", DateTime: time.Now()}}, nil).AnyTimes()

	mockRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	segSrv := service.NewSegment(mockRepo)
	segHan := NewSegment(segSrv)

	e.POST("/api/segment", segHan.CreateSegment, middleware.UseGzipReader())
	e.DELETE("/api/segment", segHan.DeleteSegment, middleware.UseGzipReader())
	e.POST("/api/user", segHan.UpdateUserSegments, middleware.UseGzipReader())
	e.GET("/api/user/:user", segHan.ReadUserSegments)
	e.GET("/api/user-history/:user", segHan.ReadUserSegmentsHistory)

	return e
}

func request(t *testing.T, ts *httptest.Server, code int, method, content, body, endpoint string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+endpoint, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", content)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	logger.InitLogger()
	log, _ := logger.GetLogger()

	log.Infoln(string(b))

	require.Equal(t, code, resp.StatusCode)

	return resp
}

func TestCreateSegment(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError, // logger should not be initialized for the first test case to get 500
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusOK,
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"text/plain",
			http.StatusBadRequest,
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\", \"slug\":\"test1\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\", \"test\":\"testing\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusUnprocessableEntity,
			"{\"slug\":\"~test\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"\"}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"slug\":\"test\",\"percent\":10}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusOK,
			"{\"slug\":\"test\",\"percent\":10}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\",\"percent\":150}",
		},
		{
			"/api/segment",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\",\"percent\":\"10\"}",
		},
	}

	for i, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()

		if i == 0 {
			logger.InitLogger()
		}
	}
}

func TestDeleteSegment(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/api/segment",
			http.MethodDelete,
			"application/json",
			http.StatusNotFound,
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodDelete,
			"application/json",
			http.StatusAccepted,
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodDelete,
			"text/plain",
			http.StatusBadRequest,
			"{\"slug\":\"test\"}",
		},
		{
			"/api/segment",
			http.MethodDelete,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\", \"slug\":\"test1\"}",
		},
		{
			"/api/segment",
			http.MethodDelete,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"test\", \"test\":\"test1\"}",
		},
		{
			"/api/segment",
			http.MethodDelete,
			"application/json",
			http.StatusBadRequest,
			"{\"slug\":\"\"}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestUpdateUserSegments(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/api/user",
			http.MethodPost,
			"application/json",
			http.StatusNotFound,
			"{\"slugs_to_add\":[\"test\"], \"slugs_to_delete\":[\"test\"], \"user_id\": \"123\"}",
		},
		{
			"/api/user",
			http.MethodPost,
			"application/json",
			http.StatusOK,
			"{\"slugs_to_add\":[\"test\"], \"slugs_to_delete\":[\"test\"], \"user_id\": \"123\"}",
		},
		{
			"/api/user",
			http.MethodPost,
			"text/plain",
			http.StatusBadRequest,
			"{\"slugs_to_add\":[\"test\"], \"slugs_to_delete\":[\"test\"], \"user_id\": \"123\"}",
		},
		{
			"/api/user",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slugs_to_add\":[\"test\"], \"slugs_to_delete\":[\"test\"], \"user_id\": \"123\", \"slugs_to_add\":[\"test1\"]}",
		},
		{
			"/api/user",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"slugs_to_add\":[\"test\"], \"slugs_to_delete\":[\"test\"], \"user_id\": \"123\", \"ip\":\"0.0.0.0\"}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestReadUserSegments(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/api/user/1",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/api/user/1",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/api/user/1",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/api/user/1",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/user/",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestReadUserSegmentsHistory(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/api/user-history/1",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/api/user-history/1",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/api/user-history/1",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/api/user-history/1",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/user-history/",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}
