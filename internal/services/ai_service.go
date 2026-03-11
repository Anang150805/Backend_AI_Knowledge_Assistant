package services

import (
	"sort"
	"strings"

	"backend-AI-Knowledge-Assistant/internal/models"
)

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

func (s *AIService) Summarize(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "Belum ada isi catatan."
	}
	if len(text) <= 180 {
		return text
	}
	return text[:180] + "..."
}

func (s *AIService) ExtractKeywords(text string) []models.NoteKeyword {
	words := strings.Fields(strings.ToLower(text))
	stopwords := map[string]bool{
		"dan": true, "yang": true, "di": true, "ke": true, "dari": true,
		"untuk": true, "dengan": true, "hari": true, "ini": true, "itu": true,
		"saya": true, "atau": true, "the": true, "a": true, "an": true,
		"pada": true, "dalam": true, "adalah": true, "karena": true,
	}

	count := map[string]int{}
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:\"'()[]{}")
		if len(w) < 3 || stopwords[w] {
			continue
		}
		count[w]++
	}

	type pair struct {
		word  string
		count int
	}

	var pairs []pair
	for k, v := range count {
		pairs = append(pairs, pair{word: k, count: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	limit := 5
	if len(pairs) < limit {
		limit = len(pairs)
	}

	var result []models.NoteKeyword
	for i := 0; i < limit; i++ {
		result = append(result, models.NoteKeyword{
			Keyword: pairs[i].word,
			Score:   float64(pairs[i].count),
		})
	}

	return result
}

func (s *AIService) AnswerFromNotes(question string, notes []models.Note) string {
	if len(notes) == 0 {
		return "Saya tidak menemukan catatan yang relevan."
	}

	var builder strings.Builder
	builder.WriteString("Berdasarkan catatan kamu:\n")
	for i, note := range notes {
		if i >= 3 {
			break
		}
		builder.WriteString("- ")
		if strings.TrimSpace(note.Title) == "" {
			builder.WriteString("Tanpa Judul")
		} else {
			builder.WriteString(note.Title)
		}
		builder.WriteString(": ")

		content := strings.TrimSpace(note.Content)
		if len(content) > 100 {
			builder.WriteString(content[:100] + "...")
		} else {
			builder.WriteString(content)
		}
		builder.WriteString("\n")
	}

	builder.WriteString("\nPertanyaan: " + question)
	return builder.String()
}