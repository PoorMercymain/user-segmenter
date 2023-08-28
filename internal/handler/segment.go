package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
	jsonduplicatechecker "github.com/PoorMercymain/user-segmenter/pkg/json-duplicate-checker"
	jsonmimechecker "github.com/PoorMercymain/user-segmenter/pkg/json-mime-checker"
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

	if slug.Slug == "" {
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

	err = h.srv.UpdateUserSegments(c.Request().Context(), userUpdate.UserID, userUpdate.SlugsToAdd, userUpdate.SlugsToDelete)

	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (h *segment) ReadUserSegments(c echo.Context) error {
	defer c.Request().Body.Close()

	userID := c.Param("user")
	if userID == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	slugs, err := h.srv.ReadUserSegments(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, appErrors.ErrorNoRows) {
			c.Response().WriteHeader(http.StatusNotFound)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
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
	return nil
}
