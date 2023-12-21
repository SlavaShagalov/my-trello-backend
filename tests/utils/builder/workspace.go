package builder

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

type WorkspaceBuilder struct {
	workspace models.Workspace
}

func NewWorkspaceBuilder() *WorkspaceBuilder {
	return &WorkspaceBuilder{}
}

func (b *WorkspaceBuilder) WithID(id int) *WorkspaceBuilder {
	b.workspace.ID = id
	return b
}

func (b *WorkspaceBuilder) WithUserID(userID int) *WorkspaceBuilder {
	b.workspace.UserID = userID
	return b
}

func (b *WorkspaceBuilder) WithTitle(title string) *WorkspaceBuilder {
	b.workspace.Title = title
	return b
}

func (b *WorkspaceBuilder) WithDescription(description string) *WorkspaceBuilder {
	b.workspace.Description = description
	return b
}

func (b *WorkspaceBuilder) WithCreatedAt(createdAt time.Time) *WorkspaceBuilder {
	b.workspace.CreatedAt = createdAt
	return b
}

func (b *WorkspaceBuilder) WithUpdatedAt(updatedAt time.Time) *WorkspaceBuilder {
	b.workspace.UpdatedAt = updatedAt
	return b
}

func (b *WorkspaceBuilder) Build() models.Workspace {
	return b.workspace
}
