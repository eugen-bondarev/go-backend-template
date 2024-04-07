package localization

type Message struct {
	id       string
	content  string
	template map[string]any
}

func NewMessage(id, content string, template ...map[string]any) Message {
	if len(template) == 0 {
		return Message{
			id:      id,
			content: content,
		}
	}

	return Message{
		id:       id,
		content:  content,
		template: template[0],
	}
}

func (m *Message) GetID() string      { return m.id }
func (m *Message) GetContent() string { return m.content }
