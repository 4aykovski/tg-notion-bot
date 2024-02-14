package notion

const (
	TitleType      = "title"
	MultiSelectTyp = "multi_select"
	CheckboxType   = "checkbox"
	TextType       = "text"
	ParagraphType  = "paragraph"
	Heading2Type   = "heading_2"
)

type Text struct {
	Content string `json:"content,omitempty"`
}

func newText(text string) *Text {
	return &Text{
		Content: text,
	}
}

type title struct {
	Type string `json:"type"`
	Text Text   `json:"text,omitempty"`
}

func newTitle(text string) *title {
	return &title{
		Type: "text",
		Text: *newText(text),
	}
}

type property struct {
	Type  string  `json:"type"`
	Title []title `json:"title,omitempty"`
}

func newTitleProperty(text string) *property {
	return newProperty("title", []title{*newTitle(text)})
}

func newProperty(propertyType string, title []title) *property {
	return &property{
		Type:  propertyType,
		Title: title,
	}
}

type richText struct {
	Type string `json:"type,omitempty"`
	Text Text   `json:"text"`
}

func newRichText(text string) *richText {
	return &richText{
		Type: "text",
		Text: *newText(text),
	}
}

type paragraph struct {
	RichText []richText `json:"rich_text,omitempty"`
}

func newParagraph(text string) *paragraph {
	return &paragraph{
		RichText: []richText{*newRichText(text)},
	}
}

type heading2 struct {
	RichText []richText `json:"rich_text,omitempty"`
}

func newHeading2(text string) *heading2 {
	return &heading2{
		RichText: []richText{*newRichText(text)},
	}
}

type block struct {
	Object    string     `json:"object"`
	Type      string     `json:"type"`
	Paragraph *paragraph `json:"paragraph,omitempty"`
	Heading2  *heading2  `json:"heading_2,omitempty"`
}

func newParagraphBlock(text string) block {
	return block{
		Object:    "block",
		Type:      "paragraph",
		Paragraph: newParagraph(text),
	}
}

func newHeading2Block(text string) block {
	return block{
		Object:   "block",
		Type:     "heading_2",
		Heading2: newHeading2(text),
	}
}

type parent struct {
	Type       string `json:"type,omitempty"`
	DatabaseId string `json:"database_id,omitempty"`
}

func newDatabaseParent(databaseId string) *parent {
	return newParent("database_id", databaseId)
}

func newParent(parentType string, databaseId string) *parent {
	return &parent{
		Type:       parentType,
		DatabaseId: databaseId,
	}
}

type page struct {
	Parent     parent              `json:"parent"`
	Properties map[string]property `json:"properties,omitempty"`
	Children   []block             `json:"children,omitempty"`
}

func newPage(parent parent, properties map[string]property, children []block) *page {
	return &page{
		Parent:     parent,
		Properties: properties,
		Children:   children,
	}
}
