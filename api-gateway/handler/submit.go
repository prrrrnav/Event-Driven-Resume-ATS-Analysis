package handler

import (
"encoding/json"	
"net/http"
"time"

"github.com/google/uuid"

"resume-ats-platform/api-gateway/kafka"
"resume-ats-platform/api-gateway/models"

)

type SubmitHandler struct {
	Producer *kafka.Producer
}

type SubmitRequest struct {
	UserID string `json:"user_id"`
	ResumeURL       string `json:"resume_url"`
	ResumeFileType string `json:"resume_file_type"`
	JobDescription string `json:"job_description"`
}

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if r.Method !=
}