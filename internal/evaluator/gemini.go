package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIResponse struct {
	Summary         string `json:"summary"`
	SkillsLatex     string `json:"skills_latex"`
	ExperienceLatex string `json:"experience_latex"`
	ProjectsLatex   string `json:"projects_latex"`
	AwardsLatex     string `json:"awards_latex"`
	Grade           string `json:"grade"`
	Score           int    `json:"score"`
}

func EvaluateAndTailor(jd string, cv string) (*AIResponse, error) {
	ctx := context.Background()
	// Engineering Bias: v1beta is often required for the latest models like Gemini 3
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Updating to the 2026 standard model
	model := client.GenerativeModel("gemini-3-flash-preview")
	model.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`
		Act as a Senior Technical Recruiter. Your task is to generate a tailored resume based ONLY on the provided CV and JD.

		CRITICAL CONSTRAINTS:
		1. NO HALLUCINATIONS: Do not invent experience. The candidate has approximately 4 years of experience. Do NOT claim 8+ years.
		2. DATA INTEGRITY: Use only dates, companies, and roles listed in the CV.
		3. COMPLETE SECTIONS: You MUST populate the following JSON fields with valid LaTeX using the resume.cls macros:
		   - "summary": A 3-4 sentence professional summary.
		   - "skills_latex": Multiple \skillItem[category={...}, skills={...}] macros.
		   - "experience_latex": \experienceItem[...] followed by \begin{itemize} and \item entries for each relevant role.
		   - "projects_latex": \projectItem[title={...}, duration={...}, keyHighlight={...}] followed by an itemized list.
		   - "awards_latex": \skillItem macros for certifications (GCP/AWS).

		CV DATA: %s
		JOB DESCRIPTION: %s
	`, cv, jd)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	var aiResult AIResponse
	rawJSON := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Robust cleaning for AI output
	rawJSON = strings.TrimSpace(rawJSON)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimSuffix(rawJSON, "```")

	err = json.Unmarshal([]byte(rawJSON), &aiResult)
	if err != nil {
		return nil, fmt.Errorf("AI parsing error: %v", err)
	}

	return &aiResult, nil
}
