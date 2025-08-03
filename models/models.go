package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// JSONB é um tipo customizado para lidar com campos JSONB do PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan into JSONB")
	}
}

// Team representa um time/equipe
type Team struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Color       string     `json:"color" db:"color"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

// Developer representa um desenvolvedor
type Developer struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	Name                   string     `json:"name" db:"name"`
	Role                   string     `json:"role" db:"role"`
	LatestPerformanceScore float64    `json:"latestPerformanceScore" db:"latest_performance_score"`
	TeamID                 *uuid.UUID `json:"teamId" db:"team_id"`
	ArchivedAt             *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedAt              time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt              time.Time  `json:"updatedAt" db:"updated_at"`
}

// PerformanceReport representa um relatório de performance
type PerformanceReport struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	DeveloperID           uuid.UUID `json:"developerId" db:"developer_id"`
	Month                 string    `json:"month" db:"month"`
	QuestionScores        JSONB     `json:"questionScores" db:"question_scores"`
	CategoryScores        JSONB     `json:"categoryScores" db:"category_scores"`
	WeightedAverageScore  float64   `json:"weightedAverageScore" db:"weighted_average_score"`
	Highlights            string    `json:"highlights" db:"highlights"`
	PointsToDevelop       string    `json:"pointsToDevelop" db:"points_to_develop"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateTeamRequest representa a requisição para criar um time
type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// UpdateTeamRequest representa a requisição para atualizar um time
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// CreateDeveloperRequest representa a requisição para criar um desenvolvedor
type CreateDeveloperRequest struct {
	Name   string     `json:"name" validate:"required,min=2"`
	Role   string     `json:"role" validate:"required,min=2"`
	TeamID *uuid.UUID `json:"teamId,omitempty"`
}

// UpdateDeveloperRequest representa a requisição para atualizar um desenvolvedor
type UpdateDeveloperRequest struct {
	Name                   *string     `json:"name,omitempty"`
	Role                   *string     `json:"role,omitempty"`
	LatestPerformanceScore *float64    `json:"latestPerformanceScore,omitempty"`
	TeamID                 *uuid.UUID  `json:"teamId,omitempty"`
}

// CreatePerformanceReportRequest representa a requisição para criar um relatório de performance
type CreatePerformanceReportRequest struct {
	DeveloperID          uuid.UUID `json:"developerId" validate:"required"`
	Month                string    `json:"month" validate:"required"`
	QuestionScores       JSONB     `json:"questionScores" validate:"required"`
	CategoryScores       JSONB     `json:"categoryScores" validate:"required"`
	WeightedAverageScore float64   `json:"weightedAverageScore" validate:"required,min=0,max=10"`
	Highlights           string    `json:"highlights"`
	PointsToDevelop      string    `json:"pointsToDevelop"`
}

// ArchiveDeveloperRequest representa a requisição para arquivar um desenvolvedor
type ArchiveDeveloperRequest struct {
	Archive bool `json:"archive"`
}
