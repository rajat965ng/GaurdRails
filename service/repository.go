package service

import (
	"findings/dao"
	"findings/model"
	"strconv"
	"strings"
)

// IRepositoryService - a service that is a repository service.
type IRepositoryService interface {
	// Create a new repository.
	Create(repo *model.Repository) error
	// FindAll returns all repositories
	FindAll() ([]model.Repository, error)
	// Update updates a repository.
	Update(repository model.Repository) error
	// Delete a repository.
	Delete(name string) error
	// AddScanDetails adds details for a specific scan.
	AddScanDetails(repository string, status model.StatusType) error
	// ExecuteScanner executes a scanner
	ExecuteScanner() error
	// Finds repositories by status
	FindByStatus(status model.StatusType) ([]model.Repository, error)
	// Lists all findings.
	ListFindings() (*model.FindingResponse, error)
}

type repository struct {
	repositoryDao dao.IRepositoryDao
}

func NewRepositoryService() *repository {
	return &repository{dao.NewRepositoryDao()}
}

func (svc *repository) Create(repo *model.Repository) error {
	return svc.repositoryDao.Create(repo)
}

func (svc *repository) FindAll() ([]model.Repository, error) {
	return svc.repositoryDao.FindAll()
}

func (svc *repository) Update(repository model.Repository) error {
	return svc.repositoryDao.Update(repository)
}

func (svc *repository) Delete(name string) error {
	return svc.repositoryDao.Delete(name)
}

func (svc *repository) AddScanDetails(name string, status model.StatusType) error {
	return svc.repositoryDao.AddScanDetails(model.Repository{Name: name}, status)
}

func (svc *repository) ExecuteScanner() error {
	repos, err := svc.FindByStatus(model.QUEUED)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		for _, scanDetail := range repo.ScanDetails {
			if scanDetail.Status == model.QUEUED {
				gitService := NewGitService(&repo, "./dir")
				_, err := gitService.GitClone()
				if err != nil {
					svc.repositoryDao.UpdateScanDetails(&scanDetail, model.FAILURE)
					return err
				}
				results, _ := gitService.GrepText([]string{"private-key", "public-key"})
				for _, result := range results {
					r := strings.Split(result, ":")
					u64, err := strconv.ParseUint(r[1], 10, 32)
					if err != nil {
						svc.repositoryDao.UpdateScanDetails(&scanDetail, model.FAILURE)
						return err
					}
					finding := model.Finding{Type: "SAST", RuleId: "G001", Path: r[0], LineNumber: uint(u64), Severity: "High", Description: r[2]}
					scanDetail.Status = model.SUCCESS
					svc.repositoryDao.SaveFindings(&finding, &scanDetail)
				}
				gitService.CleanUp()
			}
		}
	}
	return nil
}

func (svc *repository) FindByStatus(status model.StatusType) ([]model.Repository, error) {
	return svc.repositoryDao.FindRepositoryByStatus(status)
}

func (svc *repository) ListFindings() (*model.FindingResponse, error) {
	rows, err := svc.repositoryDao.FindRepositoryByStatus(model.SUCCESS)
	if err != nil {
		return nil, err
	}
	return model.GenerateFindingsByRepositoryResults(rows), nil
}
