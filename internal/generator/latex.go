package generator

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func CompileToPDF(texPath string) (string, error) {
	absPath, err := filepath.Abs(texPath)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(absPath)

	// --- NEW: Copy resume.cls to the data directory ---
	clsSource := "templates/latex/resume.cls"
	clsDest := filepath.Join(dir, "resume.cls")
	
	sourceFile, err := os.Open(clsSource)
	if err != nil {
		return "", fmt.Errorf("failed to open source .cls: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(clsDest)
	if err != nil {
		return "", fmt.Errorf("failed to create dest .cls: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy .cls file: %v", err)
	}
	// --------------------------------------------------

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+dir, absPath)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pdflatex failed: %v\nOutput: %s", err, string(output))
	}

	pdfName := filepath.Base(absPath[:len(absPath)-len(filepath.Ext(absPath))]) + ".pdf"
	return filepath.Join(dir, pdfName), nil
}