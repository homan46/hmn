package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/constant"
	"codeberg.org/rchan/hmn/dto"
	"codeberg.org/rchan/hmn/log"
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

	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(userID)
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

// GetAllNoteEndpoint will return notes in ether list or tree form.
// use can filter the output to only notes under a certain note.
//
// this api will accept the following query parameters
// - rootId=<id>
// - tree=1
//

func (n *NoteController) GetAllNoteEndpoint(c echo.Context) error {
	log.ZLog.Debug("NoteController:GetAllNoteEndpoint")
	//handle rootId query parameter
	rootIDStr := c.QueryParam("rootId")
	rootIDStr = strings.Trim(rootIDStr, " ")

	rootID := 0
	useUnder := false
	if rootIDStr != "" {
		var err error
		rootID, err = strconv.Atoi(rootIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		useUnder = true
	}

	//handle tree query parameter
	treeStr := c.QueryParam("tree")
	treeStr = strings.Trim(treeStr, " ")
	useTree := false

	if treeStr != "" {
		if treeStr != "1" {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid value for query parameter 'tree' ")
		}
		useTree = true
	}

	//get context
	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var notes []*model.Note
	if useUnder {
		notes, err = n.b.Note().GetNoteUnder(mycontext, rootID)
	} else {
		//TODO: change this to use GetAllNote. GetNoteUnder is needed for now because
		//      the tree building process currently requre input be sorted by depth
		//      but output of GetAllNote is not

		notes, err = n.b.Note().GetNoteUnder(mycontext, constant.RootNoteID)
		//notes, err = n.b.Note().GetAllNote(mycontext)
	}

	log.ZLog.Sugar().Debugf("NoteController:GetAllNoteEndpoint note retrived count: %v", len(notes))

	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if useTree {
		///TODO: This section is quite buggy and cannot return the tree correctly

		//map note to NoteTreeDto, which has a child field
		dtos := make([]*dto.NoteEntityTreeDto, 0)
		for _, n := range notes {
			dto1 := dto.NoteEntityTreeDto{}
			n.FillDto(&dto1)
			dto1.Children = make([]*dto.NoteEntityTreeDto, 0)

			dtos = append(dtos, &dto1)
		}

		// loop over and over and attach their child to them,
		// starting from the root note
		treeRootNote := dtos[0]

		// dtos is sorted by depth so just one go should be find
		// TODO: is it ACTUALLY sorted by depth tho?
		for _, fromAll := range dtos {
			for _, fromTree := range treeRootNote.Flatten() {
				if fromTree.GetID() == fromAll.GetParentID() {
					fromTree.Children = append(fromTree.Children, fromAll)
				}
			}
		}

		for _, note := range treeRootNote.Flatten() {
			sort.Slice(note.Children, func(i, j int) bool {
				return note.Children[i].GetIndex() < note.Children[j].GetIndex()
			})
		}

		tx.Commit()
		return c.JSON(http.StatusOK, treeRootNote)

	} else {
		dtos := make([]dto.NoteEntityDto, 0)
		for _, n := range notes {
			dto := dto.NoteEntityDto{}
			n.FillDto(&dto)
			dtos = append(dtos, dto)
		}

		tx.Commit()
		return c.JSON(http.StatusOK, dtos)
	}
}

func (n *NoteController) GetNoteEndpoint(c echo.Context) error {
	log.ZLog.Debug("NoteController:GetNoteEndpoint")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

func (n *NoteController) UpdateNoteEndpoint(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//bind json body
	input := new(dto.NoteDto)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	note := model.Note{}
	note.SetID(id)
	note.SetTitle(input.Title)
	note.SetContent(input.Content)
	note.SetParentID(input.ParentID)
	note.SetIndex(input.Index)

	err = n.b.Note().UpdateNote(mycontext, &note)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return c.NoContent(http.StatusOK)

}

func (n *NoteController) PatchNoteEndpoint(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	inputMap := make(map[string]interface{})
	err = json.Unmarshal(body, &inputMap)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	mycontext, tx, err := n.b.GetContextFor(userID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = n.b.Note().PatchNote(mycontext, id, inputMap)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return nil
}

func (n *NoteController) DeleteNoteEndpoint(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := readUserIDFromSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mycontext, tx, err := n.b.GetContextFor(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = n.b.Note().DeleteNote(mycontext, id)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return c.NoContent(http.StatusNoContent)

}
