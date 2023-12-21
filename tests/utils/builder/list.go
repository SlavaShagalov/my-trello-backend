package builder

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

type ListBuilder struct {
	list models.List
}

func NewListBuilder() *ListBuilder {
	return &ListBuilder{}
}

func (b *ListBuilder) WithID(id int) *ListBuilder {
	b.list.ID = id
	return b
}

func (b *ListBuilder) WithBoardID(boardID int) *ListBuilder {
	b.list.BoardID = boardID
	return b
}

func (b *ListBuilder) WithTitle(title string) *ListBuilder {
	b.list.Title = title
	return b
}

func (b *ListBuilder) WithPosition(position int) *ListBuilder {
	b.list.Position = position
	return b
}

func (b *ListBuilder) WithCreatedAt(createdAt time.Time) *ListBuilder {
	b.list.CreatedAt = createdAt
	return b
}

func (b *ListBuilder) WithUpdatedAt(updatedAt time.Time) *ListBuilder {
	b.list.UpdatedAt = updatedAt
	return b
}

func (b *ListBuilder) Build() models.List {
	return b.list
}
