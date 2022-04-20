package model

type Note struct {
	EntityEmbed
	title    string
	content  string
	parentID *int
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

func (n *Note) GetParentID() *int {
	return n.parentID
}

func (n *Note) SetParentID(parentID int) {
	n.parentID = &parentID
}

func (n *Note) GetIndex() int {
	return n.index
}

func (n *Note) SetIndex(index int) {
	n.index = index
}

type NoteLikeRO interface {
	EntityEmbedLikeRO

	GetTitle() string
	GetContent() string
	GetParentID() *int
	GetIndex() int
}

/*
type NoteLikeWO interface {
	SetTitle(title string)
	SetContent(content string)
	SetParentID(parentID int)
	SetIndex(index int)
}
*/
func NewNoteFrom(nl NoteLikeRO) *Note {
	newNote := Note{}
	mBy, mTime := nl.GetModifiedBy()
	cBy, cTime := nl.GetCreatedBy()

	newNote.EntityEmbed = *NewEntityEmbed(nl.GetID(), mTime, mBy, cTime, cBy)
	newNote.title = nl.GetTitle()
	newNote.content = nl.GetContent()
	newNote.parentID = nl.GetParentID()
	newNote.index = nl.GetIndex()

	return &newNote
}
