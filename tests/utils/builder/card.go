package builder

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

type CardBuilder struct {
	list models.Card
}

func NewCardBuilder() *CardBuilder {
	return &CardBuilder{}
}

func (b *CardBuilder) WithID(id int) *CardBuilder {
	b.list.ID = id
	return b
}

func (b *CardBuilder) WithListID(listID int) *CardBuilder {
	b.list.ListID = listID
	return b
}

func (b *CardBuilder) WithTitle(title string) *CardBuilder {
	b.list.Title = title
	return b
}

func (b *CardBuilder) WithContent(content string) *CardBuilder {
	b.list.Content = content
	return b
}

func (b *CardBuilder) WithPosition(position int) *CardBuilder {
	b.list.Position = position
	return b
}

func (b *CardBuilder) WithCreatedAt(createdAt time.Time) *CardBuilder {
	b.list.CreatedAt = createdAt
	return b
}

func (b *CardBuilder) WithUpdatedAt(updatedAt time.Time) *CardBuilder {
	b.list.UpdatedAt = updatedAt
	return b
}

func (b *CardBuilder) Build() models.Card {
	return b.list
}
