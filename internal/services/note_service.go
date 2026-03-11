package services

import (
	"backend-AI-Knowledge-Assistant/internal/models"
	"backend-AI-Knowledge-Assistant/internal/repositories"
)

type NoteService struct {
	NoteRepo  *repositories.NoteRepository
	ChatRepo  *repositories.ChatRepository
	AIService *AIService
}

func NewNoteService(noteRepo *repositories.NoteRepository, chatRepo *repositories.ChatRepository, aiService *AIService) *NoteService {
	return &NoteService{
		NoteRepo:  noteRepo,
		ChatRepo:  chatRepo,
		AIService: aiService,
	}
}

func (s *NoteService) CreateNote(userID, title, content string) (*models.Note, error) {
	note, err := s.NoteRepo.Create(userID, title, content)
	if err != nil {
		return nil, err
	}

	summary := s.AIService.Summarize(content)
	keywords := s.AIService.ExtractKeywords(content)

	_ = s.NoteRepo.SaveSummary(note.ID, summary, "local-dummy-ai")
	_ = s.NoteRepo.ReplaceKeywords(note.ID, keywords)

	return note, nil
}

func (s *NoteService) UpdateNote(noteID, userID, title, content string) (*models.Note, error) {
	note, err := s.NoteRepo.Update(noteID, userID, title, content)
	if err != nil {
		return nil, err
	}

	summary := s.AIService.Summarize(content)
	keywords := s.AIService.ExtractKeywords(content)

	_ = s.NoteRepo.SaveSummary(note.ID, summary, "local-dummy-ai")
	_ = s.NoteRepo.ReplaceKeywords(note.ID, keywords)

	return note, nil
}

func (s *NoteService) AskAI(userID, question string) (string, error) {
	notes, err := s.NoteRepo.FindAll(userID)
	if err != nil {
		return "", err
	}

	answer := s.AIService.AnswerFromNotes(question, notes)
	_, _ = s.ChatRepo.Save(userID, question, answer)

	return answer, nil
}