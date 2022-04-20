package controller

import (
	"net/http"
	"strconv"

	"codeberg.org/rchan/hmn/business"
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

type AddNoteInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ParentID int    `json:"parentId"`
	Index    int    `json:"index"`
}

func (n *NoteController) AddNoteEndpoint(c echo.Context) error {
	input := new(AddNoteInput)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	model := model.Note{}
	model.SetTitle(input.Title)
	model.SetContent(input.Content)
	model.SetParentID(input.Index)
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
	mycontext, tx, err := n.b.GetContextFor(1)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	notes, err := n.b.Note().GetAllNote(mycontext)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return c.JSON(http.StatusOK, notes)
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

	tx.Commit()
	return c.JSON(http.StatusOK, note)
}