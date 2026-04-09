package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yogibala/auto-apply/internal/database"
	"github.com/yogibala/auto-apply/internal/evaluator"
	"github.com/yogibala/auto-apply/internal/scraper"
	"github.com/yogibala/auto-apply/pkg/models"
	"net/http"
	"os"
)

func main() {
	database.InitDB() // Initialize the SQLite memory
	r := gin.Default()

	r.POST("/api/evaluate", func(c *gin.Context) {
		var input struct {
			URL string `json:"url"`
			JD  string `json:"jd"` // User can still paste text if scraper fails
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		jdText := input.JD
		// If a URL is provided, the scraper takes priority
		if input.URL != "" {
			scraped, err := scraper.ExtractJD(input.URL)
			if err == nil {
				jdText = scraped
			}
		}

		cv, _ := os.ReadFile("data/resume.md")
		result, err := evaluator.EvaluateJob(jdText, string(cv))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save to Database
		app := models.JobApplication{
			JobDescription: jdText,
			Status:         "Evaluated",
		}
		database.DB.Create(&app)

		c.JSON(http.StatusOK, gin.H{
			"evaluation": result,
			"db_id":      app.ID,
		})
	})

	r.Run(":8080")
}