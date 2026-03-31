package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
)

type Price struct {
	Min float64
	Max float64
}

type Supplier struct {
	Name   string
	Rating float64
	Region string
}

type SourceItem struct {
	ID          int64
	Platform    string
	SourceURL   string
	ExternalID  string
	Title       string
	Description string
	Images      []string
	Price       Price
	Supplier    Supplier
	Category    string
	Tags        []string
	SalesVolume int
	MinOrder    int
	Status      string
	FetchedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSourceItem(platform, sourceURL, externalID, title, description string,
	images []string, price Price, supplier Supplier, category string,
	salesVolume, minOrder int) *SourceItem {
	now := time.Now()
	return &SourceItem{
		Platform:    platform,
		SourceURL:   sourceURL,
		ExternalID:  externalID,
		Title:       title,
		Description: description,
		Images:      images,
		Price:       price,
		Supplier:    supplier,
		Category:    category,
		Tags:        []string{},
		SalesVolume: salesVolume,
		MinOrder:    minOrder,
		Status:      consts.SourceItemStatusNew,
		FetchedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *SourceItem) Select() error {
	if s.Status == consts.SourceItemStatusSelected {
		return errors.New("source item already selected")
	}
	s.Status = consts.SourceItemStatusSelected
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SourceItem) Ignore() error {
	if s.Status == consts.SourceItemStatusIgnored {
		return errors.New("source item already ignored")
	}
	s.Status = consts.SourceItemStatusIgnored
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SourceItem) AddTag(tag string) {
	for _, t := range s.Tags {
		if t == tag {
			return
		}
	}
	s.Tags = append(s.Tags, tag)
	s.UpdatedAt = time.Now()
}
