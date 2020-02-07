package datastore

import (
	"io"

	"github.com/jinzhu/gorm"
)

// DataStore is a repository for persistently storing collections of data
// related to services.
type DataStore interface {
	CreateService(svc *Service) (*Service, error)
	GetService(name string) (*Service, error)
	ListServices() ([]*Service, error)
	ListConsumerServices(ID string) ([]*Service, error)
	DeleteService(name string) error
	CreateServiceConsumer(svcName string, consumerID string) error
	DeleteServiceConsumer(svcName string, consumerID string) error
	io.Closer
}

// Open gets a handle for a database.
func Open(uri string) (DataStore, error) {
	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	return &database{DB: db}, nil
}
