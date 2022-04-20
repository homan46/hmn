package dto

type noteTitleField struct {
	Title string `db:"title" json:"title"`
}

func (n *noteTitleField) GetTitle() string {
	return n.Title
}
func (n *noteTitleField) SetTitle(title string) {
	n.Title = title
}

type noteContentField struct {
	Content string `db:"content" json:"content"`
}

func (n *noteContentField) GetContent() string {
	return n.Content
}
func (n *noteContentField) SetContent(content string) {
	n.Content = content
}

type noteParentIDField struct {
	ParentID int `db:"parent_id" json:"parentId"`
}

func (n *noteParentIDField) GetParentID() int {
	return n.ParentID
}
func (n *noteParentIDField) SetParentID(parentID int) {
	n.ParentID = parentID
}

type noteIndexField struct {
	Index int `db:"idx" json:"index"`
}

func (n *noteIndexField) GetIndex() int {
	return n.Index
}
func (n *noteIndexField) SetParentID(index int) {
	n.Index = index
}
