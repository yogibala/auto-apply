package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// ResumeData now mirrors your resume.cls macros
// We use strings for these blocks because Gemini will
// generate the formatted LaTeX items for us.
type ResumeData struct {
	FullName   string
	Email      string
	Phone      string
	LinkedIn   string
	GitHub     string
	Summary    string
	Skills     string
	Experience string
	Projects   string
	Awards     string
}

func GenerateResume(data ResumeData, outputPath string) error {
	// Engineering Bias: We must ensure the .cls is in the same folder
	// as the generated .tex during compilation.
	tmpl, err := template.ParseFiles("templates/latex/resume.tex")
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
func CompileToPDF(texPath string) (string, error) {
	// 1. Determine paths
	absPath, err := filepath.Abs(texPath)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(absPath)
	baseName := filepath.Base(absPath)
	pdfName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".pdf"

	// 2. Prepare the command
	// -interaction=nonstopmode: Don't pause for user input on errors
	// -output-directory: Where to put the resulting PDF
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+dir, absPath)

	// Engineering Bias: Set the working directory to where the .tex and .cls are.
	cmd.Dir = dir

	// 3. Execute and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// We return the output in the error so we can debug LaTeX syntax issues
		return "", fmt.Errorf("pdflatex failed: %v\nOutput: %s", err, string(output))
	}

	return filepath.Join(dir, pdfName), nil
}
