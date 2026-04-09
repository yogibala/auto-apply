package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yogibala/auto-apply/internal/database"
	"github.com/yogibala/auto-apply/internal/evaluator"
	"github.com/yogibala/auto-apply/internal/generator"
	"github.com/yogibala/auto-apply/internal/scraper"
	"net/http"
	"os"
)

func main() {
	fmt.Println("🚀 Initializing Auto-Apply Backend...")
	database.InitDB()

	r := gin.Default()

	// Fixed: Register /ping route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/api/apply", func(c *gin.Context) {
		var input struct {
			URL string `json:"url"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON body with 'url' is required"})
			return
		}

		jdText, err := scraper.ExtractJD(input.URL)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to scrape: " + err.Error()})
			return
		}

		cv, _ := os.ReadFile("data/cv.md")
		aiResult, err := evaluator.EvaluateAndTailor(jdText, string(cv))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Error: " + err.Error()})
			return
		}

		finalData := generator.ResumeData{
			FullName:   "Balayogi Meenakshisundaram",
			Email:      "mby020801@gmail.com",
			Phone:      "+91-9600608132",
			LinkedIn:   "linkedin.com/in/balayogim",
			GitHub:     "github.com/yogibala",
			Summary:    aiResult.Summary,
			Skills:     aiResult.SkillsLatex,
			Experience: aiResult.ExperienceLatex,
			Projects:   aiResult.ProjectsLatex,
			Awards:     aiResult.AwardsLatex,
		}

		outputPath := "data/tailored_resume.tex"
		generator.GenerateResume(finalData, outputPath)
		pdfPath, err := generator.CompileToPDF(outputPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "LaTeX Error: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"grade":   aiResult.Grade,
			"pdf":     pdfPath,
		})
	})

	fmt.Println("🌐 Listening on http://localhost:8080")
	r.Run(":8080")
}