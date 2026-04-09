package models

// JobApplication represents the core data structure for the system
type JobApplication struct {
	ID             string  `json:"id"`
	Company        string  `json:"company"`
	Role           string  `json:"role"`
	JobDescription string  `json:"jd"`
	MatchScore     float64 `json:"match_score"` // 1.0 to 5.0
	Grade          string  `json:"grade"`       // A-F
	TailoredResume string  `json:"tailored_resume_latex"`
	Status         string  `json:"status"`      // Applied, Interview, etc.
}