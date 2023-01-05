package dao

import (
	"findings/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IRepositoryDao - dao
type IRepositoryDao interface {
	// Create a repository.
	Create(repository *model.Repository) error
	// FindAll returns all repositories
	FindAll() ([]model.Repository, error)
	// Update updates a repository.
	Update(repository model.Repository) error
	// Delete a repository.
	Delete(name string) error
	// AddScanDetails - Updates repository and saves findings.
	AddScanDetails(repository model.Repository, status model.StatusType) error
	// Updates the ScanDetails and SaveFindings.
	UpdateScanDetails(scanDetail *model.ScanDetail, status model.StatusType) error
	// Finds a repository by status
	FindRepositoryByStatus(status model.StatusType) ([]model.Repository, error)
	// SaveFindings saves the findings to the database.
	SaveFindings(finding *model.Finding, scanDetail *model.ScanDetail) error
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryDao() *repository {
	return &repository{NewDatabaseInstance().GetConnection()}
}

func (dao *repository) Create(repo *model.Repository) error {
	dao.db.Debug().Create(repo)
	return nil
}

func (dao *repository) FindAll() ([]model.Repository, error) {
	var repos []model.Repository
	dao.db.Find(&repos)
	return repos, nil
}

func (dao *repository) Update(repository model.Repository) error {
	dao.db.Debug().Where("name = ?", repository.Name).Updates(repository)
	return nil
}

func (dao *repository) Delete(name string) error {
	var repository model.Repository
	dao.db.Debug().Preload("ScanDetails").Where("name = ?", name).Find(&repository).Select(clause.Associations).Delete(&repository)
	return nil
}

func (dao *repository) AddScanDetails(repository model.Repository, status model.StatusType) error {
	scanDetails := &model.ScanDetail{
		Status: status,
	}
	repository.ScanDetails = append(repository.ScanDetails, *scanDetails)
	dao.db.Debug().Where("name = ?", repository.Name).Find(&repository).Updates(scanDetails)
	return nil
}

func (dao *repository) FindRepositoryByStatus(status model.StatusType) ([]model.Repository, error) {
	var repos []model.Repository
	dao.db.Debug().Preload("ScanDetails", "status = ?", status).Preload("ScanDetails.Findings").Find(&repos)
	return repos, nil
}

func (dao *repository) SaveFindings(finding *model.Finding, scanDetail *model.ScanDetail) error {
	scanDetail.Findings = append(scanDetail.Findings, *finding)
	dao.db.Debug().Save(scanDetail)
	return nil
}

func (dao *repository) UpdateScanDetails(scanDetail *model.ScanDetail, status model.StatusType) error {
	scanDetail.Status = status
	dao.db.Save(scanDetail)
	return nil
}
