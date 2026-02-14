package domain

import "time"

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusRunning  BuildStatus = "running"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failed"
	BuildStatusCanceled BuildStatus = "canceled"
)

type Build struct {
	ID                string      `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	RepoUrl           string      `json:"repo_url" validate:"required,url"`
	Ref               string      `json:"ref" validate:"required"`
	Command           string      `json:"command" validate:"required"`
	Status            BuildStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	FinishedAt        *time.Time  `json:"finished_at"`
	CancelRequestedAt *time.Time  `json:"cancel_requested_at"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ExitCode          int         `json:"exit_code"`
	Error             string      `json:"error"`
}

type BuildLog struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BuildID   string    `json:"build_id" gorm:"type:uuid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
