package factory

import (
	"fmt"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/repository/localcache"
)

const (
	// Memory storage
	Memory = "memory"
)

// GetRepository get repository for type
func GetRepository(repositoryType string) (repository.EventRepository, error) {
	switch repositoryType {
	case Memory:
		return localcache.NewEventLocalStorage(), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", repositoryType)
	}
}
