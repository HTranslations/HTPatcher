package service

import (
	"fmt"
	"htpatcher/internal/domain"
)

// PatchDownloadRepository interface for patch download operations
type PatchDownloadRepository interface {
	Download(patchDownloadId string) (string, error)
}

// DownloadService handles patch downloading
type DownloadService struct {
	patchRepo PatchDownloadRepository
	logger    Logger
}

// NewDownloadService creates a new download service
func NewDownloadService(patchRepo PatchDownloadRepository, logger Logger) *DownloadService {
	return &DownloadService{
		patchRepo: patchRepo,
		logger:    logger,
	}
}

// DownloadPatch downloads a patch and returns its information
func (s *DownloadService) DownloadPatch(patchDownloadId string, loadPatchFunc func(filePath string) (*domain.PatchInfo, error)) (*domain.PatchInfo, error) {
	s.logger.Info(fmt.Sprintf("Downloading patch with download ID %s", patchDownloadId))
	filePath, err := s.patchRepo.Download(patchDownloadId)
	if err != nil {
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("Downloaded patch to %s", filePath))
	s.logger.Info("Loading patch into memory...")

	// Use the provided function to load patch info
	return loadPatchFunc(filePath)
}

