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
