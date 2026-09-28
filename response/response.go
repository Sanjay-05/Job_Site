package response

import "jobsite/models"

type JobListResult struct {
	Jobs  []models.Job `json:"jobs"`
	Total int64        `json:"total"`
}
