package builder

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

type BoardBuilder struct {
	workspace models.Board
}

func NewBoardBuilder() *BoardBuilder {
	return &BoardBuilder{}
}

func (b *BoardBuilder) WithID(id int) *BoardBuilder {
	b.workspace.ID = id
	return b
}

func (b *BoardBuilder) WithWorkspaceID(workspacesID int) *BoardBuilder {
	b.workspace.WorkspaceID = workspacesID
	return b
}

func (b *BoardBuilder) WithTitle(title string) *BoardBuilder {
	b.workspace.Title = title
	return b
}

func (b *BoardBuilder) WithDescription(description string) *BoardBuilder {
	b.workspace.Description = description
	return b
}

func (b *BoardBuilder) WithBackground(background string) *BoardBuilder {
	b.workspace.Background = &background
	return b
}

func (b *BoardBuilder) WithCreatedAt(createdAt time.Time) *BoardBuilder {
	b.workspace.CreatedAt = createdAt
	return b
}

func (b *BoardBuilder) WithUpdatedAt(updatedAt time.Time) *BoardBuilder {
	b.workspace.UpdatedAt = updatedAt
	return b
}

func (b *BoardBuilder) Build() models.Board {
	return b.workspace
}
