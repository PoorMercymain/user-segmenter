package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

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

// @Title UserSegmenter API
// @Description Сервис динамического сегментирования пользователей.
// @Version 1.0

// @BasePath /api
// @Host localhost:8080

// @Tag.name Slugs
// @Tag.description "Группа запросов для управления сегментами"

// CreateSegment godoc
// @Tags Slugs
// @Summary Запрос для создания нового сегмента
// @Description Запрос для создания сегмента по уникальному названию
// @Accept json
// @Produce json
// @Param slug string
// @Success 200
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 422 {string} string "Передан не slug"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /segment [post]
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

	var slug domain.Slug

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

func (h *segment) UpdateUserSegments(c echo.Context) error {
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

func (h *segment) ReadUserSegments(c echo.Context) error {
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

func (h *segment) ReadUserSegmentsHistory(c echo.Context) error {
	defer c.Request().Body.Close()

	userID := c.Param("user")

	startDateStr := c.QueryParam("start")
	endDateStr := c.QueryParam("end")
	if startDateStr == "" {
		startDateStr = "1900-01" // some date when user definetely could not be added
	}

	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01")
	}

	startDate, err := time.Parse("2006-01", startDateStr)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	endDateBuf, err := time.Parse("2006-01", endDateStr)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	endDateBuf = endDateBuf.AddDate(0, 1, 0)
	endDateBuf = endDateBuf.Add(-time.Millisecond)

	endDateStr = endDateBuf.Format(time.RFC3339)

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	history, err := h.srv.ReadUserSegmentsHistory(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if len(history) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	w := csv.NewWriter(c.Response().Writer)
	c.Response().Header().Set("Content-Type", "text/csv")
	w.Comma = ';'

	for _, historyElement := range history {
		historyElementStrSlice := []string{historyElement.UserID, historyElement.Slug, historyElement.Operation, historyElement.DateTime.Format(time.RFC3339)}
		if err := w.Write(historyElementStrSlice); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		w.Flush()
	}

	return err
}
