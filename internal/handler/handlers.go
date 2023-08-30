package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	jsonduplicatechecker "github.com/PoorMercymain/user-segmenter/pkg/json-duplicate-checker"
	jsonmimechecker "github.com/PoorMercymain/user-segmenter/pkg/json-mime-checker"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

type segment struct {
	srv domain.SegmentService
}

func NewSegment(srv domain.SegmentService) *segment {
	return &segment{srv: srv}
}

// @title UserSegmenter API
// @version 1.0
// @description Сервис динамического сегментирования пользователей

// @host localhost:8080
// @BasePath /

// @Tag.name Segments
// @Tag.description Группа запросов для управления списком существующих сегментов

// @Tag.name Users
// @Tag.description Группа запросов для управления сегментами пользователя

// @Schemes http

// @Tags Segments
// @Summary Запрос для создания нового сегмента
// @Description Запрос для создания сегмента по уникальному названию
// @Accept json
// @Param input body domain.Slug true "segment info"
// @Success 200
// @Failure 404
// @Failure 500
// @Failure 400
// @Failure 409
// @Router /api/segment [post]
func (h *segment) CreateSegment(c echo.Context) error {
	defer c.Request().Body.Close()

	if !jsonmimechecker.IsJSONContentTypeCorrect(c.Request()) {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	bytesToCheck, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	reader := bytes.NewReader(bytes.Clone(bytesToCheck))

	err = jsonduplicatechecker.CheckDuplicatesInJSON(json.NewDecoder(reader), nil)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(bytesToCheck))

	d := json.NewDecoder(c.Request().Body)
	d.DisallowUnknownFields()

	var slug domain.Slug

	if err := d.Decode(&slug); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	if slug.Slug == "" || slug.PercentOfUsers < 0 || slug.PercentOfUsers > 100 {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	err = h.srv.CreateSegment(c.Request().Context(), slug.Slug)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNotASlug) {
			c.Response().WriteHeader(http.StatusUnprocessableEntity)
			return err
		}

		if errors.Is(err, appErrors.ErrorUniqueViolation) {
			c.Response().WriteHeader(http.StatusConflict)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if slug.PercentOfUsers != 0 {
		err = h.srv.AddSegmentToPercentOfUsers(c.Request().Context(), slug.Slug, slug.PercentOfUsers)
		if err != nil && !errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// @Tags Segments
// @Summary Запрос для удаления сегмента
// @Description Запрос для удаления сегмента из списка существующих сегментов по уникальному названию
// @Accept json
// @Param input body domain.SlugNoPercent true "segment info"
// @Success 202
// @Failure 404
// @Failure 500
// @Failure 400
// @Router /api/segment [delete]
func (h *segment) DeleteSegment(c echo.Context) error {
	defer c.Request().Body.Close()

	if !jsonmimechecker.IsJSONContentTypeCorrect(c.Request()) {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	bytesToCheck, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	reader := bytes.NewReader(bytes.Clone(bytesToCheck))

	err = jsonduplicatechecker.CheckDuplicatesInJSON(json.NewDecoder(reader), nil)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(bytesToCheck))

	d := json.NewDecoder(c.Request().Body)
	d.DisallowUnknownFields()

	var slug domain.SlugNoPercent

	if err := d.Decode(&slug); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	if slug.Slug == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	err = h.srv.DeleteSegment(c.Request().Context(), slug.Slug)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		log, _ := logger.GetLogger()
		log.Infoln(err)
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}

type user struct {
	srv domain.UserService
}

func NewUser(srv domain.UserService) *user {
	return &user{srv: srv}
}

// @Tags Users
// @Summary Запрос обновления сегментов пользователя
// @Description Запрос для обновления списка сегментов пользователя
// @Accept json
// @Param input body domain.UserUpdate true "user segment info"
// @Success 200
// @Failure 404
// @Failure 500
// @Failure 400
// @Router /api/user [post]
func (h *user) UpdateUserSegments(c echo.Context) error {
	defer c.Request().Body.Close()
	log, err := logger.GetLogger()
	if err != nil {
		return err
	}

	if !jsonmimechecker.IsJSONContentTypeCorrect(c.Request()) {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	bytesToCheck, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	reader := bytes.NewReader(bytes.Clone(bytesToCheck))

	err = jsonduplicatechecker.CheckDuplicatesInJSON(json.NewDecoder(reader), nil)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(bytesToCheck))

	d := json.NewDecoder(c.Request().Body)
	d.DisallowUnknownFields()

	var userUpdate domain.UserUpdate

	if err := d.Decode(&userUpdate); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	if len(userUpdate.SlugsToAdd) != len(userUpdate.TTL) && len(userUpdate.TTL) != 0 {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	var TTLs []time.Time
	for _, TTL := range userUpdate.TTL {
		oneOfTTLs, err := time.Parse(time.RFC3339, TTL)
		if err != nil {
			log.Infoln(err)
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}
		TTLs = append(TTLs, oneOfTTLs)
	}

	err = h.srv.UpdateUserSegments(c.Request().Context(), userUpdate.UserID, userUpdate.SlugsToAdd, userUpdate.SlugsToDelete)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	for i, TTL := range TTLs {
		err = h.srv.CreateDeletionTime(c.Request().Context(), userUpdate.UserID, userUpdate.SlugsToAdd[i], TTL)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// @Tags Users
// @Summary Запрос чтения сегментов пользователя
// @Description Запрос для получения списка сегментов пользователя
// @Produce json
// @Param id path string true "user id"
// @Success 200
// @Success 204
// @Failure 404
// @Failure 500
// @Router /api/user/{id} [get]
func (h *user) ReadUserSegments(c echo.Context) error {
	defer c.Request().Body.Close()

	userID := c.Param("user")

	slugs, err := h.srv.ReadUserSegments(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if len(slugs) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	var slugsBytes []byte
	buf := bytes.NewBuffer(slugsBytes)
	err = json.NewEncoder(buf).Encode(slugs)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	c.Response().Header().Set("Content-Type", "application/json")

	if len(buf.Bytes()) > 1024 {
		acceptsEncoding := c.Request().Header.Values("Accept-Encoding")
		for _, encoding := range acceptsEncoding {
			if strings.Contains(encoding, "gzip") {
				c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
				gz := gzip.NewWriter(c.Response().Writer)
				defer gz.Close()

				c.Response().Writer = domain.RespWriter{
					Writer:         gz,
					ResponseWriter: c.Response().Writer,
				}
				break
			}
		}
	}

	_, err = c.Response().Write(buf.Bytes())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusOK)
	return nil
}

type report struct {
	srv domain.ReportService
}

func NewReport(srv domain.ReportService) *report {
	return &report{srv: srv}
}

// @Tags Users
// @Summary Запрос формирования отчета по истории сегментов пользователя
// @Description Запрос для создания отчета по истории сегментов пользователя в формате csv
// @Produce plain
// @Param id path string true "user id"
// @Param start query string false "start date"
// @Param end query string false "end date"
// @Success 200
// @Failure 404
// @Failure 400
// @Failure 500
// @Router /api/user-history/{id} [get]
func (h *report) ReadUserSegmentsHistory(c echo.Context) error {
	defer c.Request().Body.Close()

	userID := c.Param("user")

	startDateStr := c.QueryParam("start")
	endDateStr := c.QueryParam("end")
	if startDateStr == "" {
		startDateStr = "1900-01" // some date when user definitely could not be added
	}

	startDate, err := time.Parse("2006-1", startDateStr)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	endDateBuf, err := time.Parse("2006-1", endDateStr)
	if err != nil && endDateStr == "" {
		endDateStr = time.Now().Format("2006-1")
		endDateBuf, err = time.Parse("2006-1", endDateStr)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
	} else if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	endDateBuf = endDateBuf.AddDate(0, 1, 0)
	endDateBuf = endDateBuf.Add(-time.Millisecond)

	endDateStr = endDateBuf.Format(time.RFC3339)

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	filename, err := h.srv.CreateCSV(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().Write([]byte(c.Request().Context().Value(domain.Key("server")).(string) + "/api/" + strings.Replace(filename, "\\", "/", -1)))
	return nil
}

// @Tags Users
// @Summary Запрос чтения отчета по истории сегментов пользователя
// @Description Запрос для получения отчета по истории сегментов пользователя в формате csv
// @Produce text/csv
// @Param filename path string true "report filename"
// @Success 200
// @Success 204
// @Failure 404
// @Failure 500
// @Router /api/reports/{filename} [get]
func (h *report) ReadUserSegmentsHistoryReport(c echo.Context) error {
	defer c.Request().Body.Close()

	reportName := c.Param("report")

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+reportName)

	err := h.srv.SendCSVReportFile(reportName, c.Response())
	if err != nil {
		if errors.Is(err, appErrors.ErrorFileNotFound) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		if errors.Is(err, appErrors.ErrorEmptyFile) {
			c.Response().WriteHeader(http.StatusNoContent)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	return nil
}
