package generator

import (
	"os"
	"text/template"
)

// ResumeData now mirrors your resume.cls macros
// We use strings for these blocks because Gemini will 
// generate the formatted LaTeX items for us.
type ResumeData struct {
	FullName    string
	Email       string
	Phone       string
	LinkedIn    string
	GitHub      string
	Summary     string
	Skills      string 
	Experience  string
	Projects    string
	Awards      string
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