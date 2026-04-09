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
		Act as a Senior Tech Recruiter. Using the JD and CV provided, generate a tailored resume.
		Return a JSON object that fits the resume.cls macros exactly.
		
		JD: %s
		CV: %s
	`, jd, cv)

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
