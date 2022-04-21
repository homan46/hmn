package controller

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/dto"
	"codeberg.org/rchan/hmn/model"
	"github.com/labstack/echo/v4"
)

type NoteController struct {
	b business.BusinessLayer
}

func NewNoteController(b business.BusinessLayer) *NoteController {
	return &NoteController{
		b: b,
	}
}

func (n *NoteController) AddNoteEndpoint(c echo.Context) error {
	input := new(dto.NoteDto)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	model := model.Note{}
	model.SetTitle(input.Title)
	model.SetContent(input.Content)
	model.SetParentID(input.ParentID)
	model.SetIndex(input.Index)

	mycontext, tx, err := n.b.GetContextFor(1)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = n.b.Note().AddNote(mycontext, &model)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return c.JSON(http.StatusCreated, model)
}

func (n *NoteController) GetAllNoteEndpoint(c echo.Context) error {
	rootIDStr := c.QueryParam("rootId")
	rootIDStr = strings.Trim(rootIDStr, " ")

	var rootID int = 0
	var useUnder = false
	if rootIDStr != "" {
		var err error
		rootID, err = strconv.Atoi(rootIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		useUnder = true
	}

	mycontext, tx, err := n.b.GetContextFor(1)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var notes []*model.Note
	if useUnder {
		notes, err = n.b.Note().GetNoteUnder(mycontext, rootID)
	} else {
		notes, err = n.b.Note().GetAllNote(mycontext)
	}

	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	dtos := make([]dto.NoteEntityDto, 0)
	for _, n := range notes {
		dto := dto.NoteEntityDto{}
		n.FillDto(&dto)
		dtos = append(dtos, dto)
	}

	tx.Commit()
	return c.JSON(http.StatusOK, dtos)
}

func (n *NoteController) GetNoteEndpoint(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(1)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	note, err := n.b.Note().GetNote(mycontext, id)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	dto := dto.NoteEntityDto{}
	note.FillDto(&dto)

	tx.Commit()
	return c.JSON(http.StatusOK, dto)
}
