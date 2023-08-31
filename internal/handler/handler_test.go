package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
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

	mockSegRepo := mocks.NewMockSegmentRepository(ctrl)
	mockUsrRepo := mocks.NewMockUserRepository(ctrl)
	mockRepRepo := mocks.NewMockReportRepository(ctrl)

	mockSegRepo.EXPECT().CreateSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockSegRepo.EXPECT().DeleteSegment(gomock.Any(), gomock.Any()).Return(appErrors.ErrorNoRows).MaxTimes(1)
	mockSegRepo.EXPECT().DeleteSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockUsrRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrorNoRows).MaxTimes(1)
	mockUsrRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockUsrRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(nil, appErrors.ErrorNoRows).MaxTimes(1)
	mockUsrRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(make([]string, 0), nil).MaxTimes(1)
	mockUsrRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(nil, appErrors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockUsrRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil).AnyTimes()

	mockRepRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, appErrors.ErrorNoRows).MaxTimes(1)
	mockRepRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, appErrors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.HistoryElem, 0), nil).MaxTimes(1)
	mockRepRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.HistoryElem{{UserID: "1", Slug: "a", Operation: "addition", DateTime: time.Now()}}, nil).AnyTimes()

	mockSegRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockSegRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockRepRepo.EXPECT().CreateCSV(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", appErrors.ErrorNoRows).MaxTimes(1)
	mockRepRepo.EXPECT().CreateCSV(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", appErrors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepRepo.EXPECT().CreateCSV(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("report1.csv", nil).AnyTimes()

	mockRepRepo.EXPECT().SendCSVReportFile(gomock.Any(), gomock.Any()).Return(appErrors.ErrorBadFilename).MaxTimes(1)
	mockRepRepo.EXPECT().SendCSVReportFile(gomock.Any(), gomock.Any()).Return(appErrors.ErrorFileNotFound).MaxTimes(1)
	mockRepRepo.EXPECT().SendCSVReportFile(gomock.Any(), gomock.Any()).Return(appErrors.ErrorEmptyFile).MaxTimes(1)
	mockRepRepo.EXPECT().SendCSVReportFile(gomock.Any(), gomock.Any()).Return(appErrors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepRepo.EXPECT().SendCSVReportFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	segSrv := service.NewSegment(mockSegRepo)
	usrSrv := service.NewUser(mockUsrRepo)
	repSrv := service.NewReport(mockRepRepo)

	segHan := NewSegment(segSrv)
	usrHan := NewUser(usrSrv)
	repHan := NewReport(repSrv)

	e.POST("/api/segment", segHan.CreateSegment, middleware.UseGzipReader())
	e.DELETE("/api/segment", segHan.DeleteSegment, middleware.UseGzipReader())
	e.POST("/api/user", usrHan.UpdateUserSegments, middleware.UseGzipReader())
	e.GET("/api/user/:user", usrHan.ReadUserSegments)
	e.GET("/api/user-history/:user", repHan.CreateUserSegmentsHistoryReport, middleware.AddServerAddressToContext(""))
	e.GET("/api/reports/:report", repHan.ReadUserSegmentsHistoryReport)

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
			http.StatusAccepted,
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

func TestCreateUserSegmentsHistoryReport(t *testing.T) {
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
			http.StatusOK,
			"",
		},
		{
			"/api/user-history/1?start=123",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/api/user-history/1?start=1971-1-11&end=123",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/api/user-history/",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/api/user-history/1?start=1971-1",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/user-history/1?end=1971-1",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/user-history/1?start=1971-1&end=1970-12",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/api/user-history/1?exact=1971-1",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/user-history/1?exact=1971-1&end=1971-2",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/api/user-history/1?end=1899-1",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
	}
	logger.InitLogger()
	log, _ := logger.GetLogger()

	for i, testCase := range testTable {
		log.Infoln(i)
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestReadUserSegmentsHistoryReport(t *testing.T) {
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
			"/api/reports/repo.csv",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/api/reports/report1.csv",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/api/reports/report1.csv",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/api/reports/report1.csv",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/api/reports/report1.csv",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
		{
			"/api/reports/",
			http.MethodGet,
			"",
			http.StatusNotFound,
			"",
		},
	}
	logger.InitLogger()
	log, _ := logger.GetLogger()

	for i, testCase := range testTable {
		log.Infoln(i)
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}
