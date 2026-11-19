package handler

import (
	"errors"
	"net/http"
	"strconv"

	"phone-book/model"
	"phone-book/repository"
	"phone-book/service"

	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	svc service.ContactService
}

func NewContactHandler(svc service.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

// GetAll — GET /api/contacts
func (h *ContactHandler) GetAll(c echo.Context) error {
	contacts, err := h.svc.ListContacts(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, contacts)
}

// GetByID — GET /api/contacts/:id
func (h *ContactHandler) GetByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	contact, err := h.svc.GetContact(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, contact)
}

// Create — POST /api/contacts
func (h *ContactHandler) Create(c echo.Context) error {
	var contact model.Contact
	if err := c.Bind(&contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.svc.CreateContact(c.Request().Context(), &contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, contact)
}

// Update — PUT /api/contacts/:id
func (h *ContactHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var contact model.Contact
	if err := c.Bind(&contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	contact.ID = id
	if err := h.svc.UpdateContact(c.Request().Context(), &contact); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, contact)
}

// Delete — DELETE /api/contacts/:id
func (h *ContactHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.svc.DeleteContact(c.Request().Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func parseID(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}