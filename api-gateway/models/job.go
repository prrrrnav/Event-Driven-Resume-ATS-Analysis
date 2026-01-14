package models

import "time"

type ResumeUploadEvent struct {
	EventID string `json:"event_id"`
	EventType string `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`

	JobID string `json:"job_id"`
	UserID string `json:"user_id"`

	Resume ResumeInfo `json:"resume"`
	JD JobDescription `json:"jobDescription"`

}

type ResumeInfo struct {
	FileURL string `json:"file_url"`
	FileType string `json:"file_type"`
}

type JobDescription struct {
	Text string `json:"text"`
}