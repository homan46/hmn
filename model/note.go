package model

type Note struct {
	Entity
	title    string
	content  string
	parentID int
	index    int
}

func (n *Note) GetTitle() string {
	return n.title
}

func (n *Note) SetTitle(title string) {
	n.title = title
}

func (n *Note) GetContent() string {
	return n.content
}

func (n *Note) SetContent(content string) {
	n.content = content
}

func (n *Note) GetParentID() int {
	return n.parentID
}

func (n *Note) SetParentID(parentID int) {
	n.parentID = parentID
}

func (n *Note) GetIndex() int {
	return n.index
}

func (n *Note) SetIndex(index int) {
	n.index = index
}

type NoteLikeRO interface {
	EntityLikeRO

	GetTitle() string
	GetContent() string
	GetParentID() int
	GetIndex() int
}

type NoteLikeWO interface {
	EntityLikeWO
	SetTitle(title string)
	SetContent(content string)
	SetParentID(parentID int)
	SetIndex(index int)
}

func NewNoteFrom(nl NoteLikeRO) *Note {
	newNote := Note{}

	// mBy, mTime := nl.GetModifiedBy()
	// cBy, cTime := nl.GetCreatedBy()

	newNote.Entity = *NewEntity(nl.GetID(),
		nl.GetModifiedTime(), nl.GetModifiedBy(),
		nl.GetCreatedTime(), nl.GetCreatedBy())
	newNote.title = nl.GetTitle()
	newNote.content = nl.GetContent()
	newNote.parentID = nl.GetParentID()
	newNote.index = nl.GetIndex()

	return &newNote
}
func (n *Note) FillDto(emptyDto NoteLikeWO) error {
	emptyDto.SetID(n.id)
	emptyDto.SetCreatedBy(n.createdBy)
	emptyDto.SetCreatedTime(n.createdTime)
	emptyDto.SetModifiedBy(n.modifiedBy)
	emptyDto.SetModifiedTime(n.modifiedTime)
	emptyDto.SetTitle(n.title)
	emptyDto.SetContent(n.content)
	emptyDto.SetParentID(n.parentID)
	emptyDto.SetIndex(n.index)
	return nil
}
