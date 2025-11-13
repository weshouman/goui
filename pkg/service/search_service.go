package service

import (
	"github.com/sirupsen/logrus"

	"github.com/ourorg/goui/pkg/domain"
)

// SearchService handles search operations
type SearchService struct{}

// NewSearchService creates a new search service
func NewSearchService() *SearchService {
	return &SearchService{}
}

// ExecuteSearch performs search operations based on the current state
// This is a generic implementation that apps can extend or override
func (s *SearchService) ExecuteSearch(term string, currentState *domain.State) {
	logrus.Debugf("Executing search with term: %s", term)

	// Update the current state's Args with search results
	if currentState != nil {
		// Generic search implementation
		currentState.Args["searchTerm"] = term
		currentState.Args["searchActive"] = true

		logrus.Debugf("Search results stored in state args for term: %s", term)
	}
}