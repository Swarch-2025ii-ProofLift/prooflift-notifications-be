package events

import (
	"fmt"

	"github.com/google/uuid"
)

type EventType string

const (
	EventCommentCreated EventType = "COMMENT_CREATED"
	EventReactionAdded  EventType = "REACTION_ADDED"
)

type Event struct {
	Type EventType `json:"type"`

	UserID string `json:"user_id"`

	ActorID string `json:"actor_id"`

	PostID string `json:"post_id"`

	CommentID string `json:"comment_id,omitempty"`

	Message string `json:"message"`
}

func (e *Event) ParseIDs() (userID uuid.UUID, actorID uuid.UUID, postID uuid.UUID, commentID *uuid.UUID, err error) {
	if e == nil {
		err = fmt.Errorf("event is nil")
		return
	}

	userID, err = uuid.Parse(e.UserID)
	if err != nil {
		err = fmt.Errorf("invalid user_id: %w", err)
		return
	}

	actorID, err = uuid.Parse(e.ActorID)
	if err != nil {
		err = fmt.Errorf("invalid actor_id: %w", err)
		return
	}

	postID, err = uuid.Parse(e.PostID)
	if err != nil {
		err = fmt.Errorf("invalid post_id: %w", err)
		return
	}

	if e.CommentID != "" {
		parsedCommentID, parseErr := uuid.Parse(e.CommentID)
		if parseErr != nil {
			err = fmt.Errorf("invalid comment_id: %w", parseErr)
			return
		}
		commentID = &parsedCommentID
	}

	if e.Type == EventCommentCreated && commentID == nil {
		err = fmt.Errorf("comment_id is required for COMMENT_CREATED events")
		return
	}

	return
}

func (e *Event) IsValidType() bool {
	switch e.Type {
	case EventCommentCreated, EventReactionAdded:
		return true
	default:
		return false
	}
}
