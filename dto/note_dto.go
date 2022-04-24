package dto

type NoteEntityDto struct {
	EntityDto
	noteTitleField
	noteContentField
	noteParentIDField
	noteIndexField
}

type NoteDto struct {
	noteTitleField
	noteContentField
	noteParentIDField
	noteIndexField
}

type NoteEntityTreeDto struct {
	EntityDto
	noteTitleField
	noteContentField
	noteParentIDField
	noteIndexField
	noteChildrenField[NoteEntityTreeDto]
}

func (n *NoteEntityTreeDto) Flatten() []*NoteEntityTreeDto {
	a := make([]*NoteEntityTreeDto, 0)
	//add this node
	a = append(a, n)
	//add all child node
	for _, child := range n.Children {
		a = append(a, child.Flatten()...)
	}
	return a
}
