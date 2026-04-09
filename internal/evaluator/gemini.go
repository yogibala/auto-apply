package evaluator

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func EvaluateJob(jd string, cvContent string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")

	prompt := fmt.Sprintf(`
    Act as an Expert LaTeX Resume Writer. 
    Use the provided JD and CV to generate a tailored resume.
    
    CRITICAL: Use these LaTeX macros ONLY:
    - \skillItem[category={...}, skills={...}]
    - \experienceItem[company={...}, location={...}, position={...}, duration={...}]
    - \projectItem[title={...}, duration={...}, keyHighlight={...}]
    
    Return a JSON object:
    {
        "summary": "...",
        "skills_latex": "\\skillItem[category={...}, skills={...}] \\\\ ...",
        "experience_latex": "\\experienceItem[...] \\begin{itemize} ... \\end{itemize}",
        "projects_latex": "..."
    }
    
    JD: %s
    CV: %s
`, jd, cv)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// Extract text from the first candidate
	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
