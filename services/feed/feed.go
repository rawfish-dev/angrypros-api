package feed

import (
	"github.com/rawfish-dev/angrypros-api/models"
	"github.com/rawfish-dev/angrypros-api/services/storage"
)

var _ FeedService = new(Service)

type FeedService interface {
	GetFeedItems(beforeTimestamp int64, size int) ([]models.Entry, error)
}

type Service struct {
	entryStorage storage.EntryStorage
}

func NewService(e storage.EntryStorage) *Service {
	return &Service{
		entryStorage: e,
	}
}

func (s Service) GetFeedItems(beforeTimestamp int64, size int) ([]models.Entry, error) {
	entries, err := s.entryStorage.GetEntries(beforeTimestamp, size)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
