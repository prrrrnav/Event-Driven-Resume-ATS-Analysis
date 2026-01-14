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
	if r.Method != http.MethodPost{
		http.Error(w,"method not allowed", http.StatusMethodNotAllowed)
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err!= nil{
		http.Error(w,"invalid request body", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()

	event := models.ResumeUploadEvent{
		EventID: uuid.New().String(),
		EventType: "resume.uploaded",
		Timestamp: time.Now().UTC(),
		JobID: jobID,
		UserID:req.UserID,
	}

	event.Resume.FileURL = req.ResumeURL
	event.Resume.FileType = req.ResumeFileType
	event.JD.Text = req.JobDescription

	if err := h.Producer.Publish(r.Context(), event); err != nil{
		http.Error(w,"Failed to queue job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status" : "queued",
	})
}