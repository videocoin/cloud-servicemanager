package datastore

import (
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

// database implements the DataStore interface.
type database struct {
	*gorm.DB
}

// Close closes the database connection.
func (db *database) Close() error {
	return db.DB.Close()
}

// CreateService creates a managed service.
func (db *database) CreateService(svc *Service) (*Service, error) {
	if err := db.Create(svc).Error; err != nil {
		return nil, err
	}
	return svc, nil
}

// GetService gets a managed service.
func (db *database) GetService(name string) (*Service, error) {
	var svc Service
	if err := db.Where("name = ?", name).First(&svc).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

// ListServices lists managed services.
func (db *database) ListServices() ([]*Service, error) {
	var svcs []*Service
	if err := db.Find(&svcs).Error; err != nil {
		return nil, err
	}
	return svcs, nil
}

// DeleteService deletes a managed service.
func (db *database) DeleteService(name string) error {
	return db.Delete(Service{}, "name = ?", name).Error
}

// CreateServiceConsumer creates the association between a service and a consumer.
func (db *database) CreateServiceConsumer(svcName string, consumerID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// create consumer if not exists
		consumer := &Consumer{ID: consumerID}
		if err := tx.FirstOrCreate(consumer).Error; err != nil {
			return err
		}

		var svc Service
		if err := tx.Where("name = ?", svcName).First(&svc).Error; err != nil {
			return err
		}

		return tx.Model(&svc).Association("Consumers").Append(consumer).Error
	})
}

// ListConsumerServices lists consumer services.
func (db *database) ListConsumerServices(ID string) ([]*Service, error) {
	var svcs []*Service
	if err := db.Model(&Consumer{ID: ID}).Association("Services").Find(&svcs).Error; err != nil {
		return nil, err
	}
	return svcs, nil
}

// DeleteServiceConsumer deletes the association between a service and a consumer.
func (db *database) DeleteServiceConsumer(svcName string, consumerID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var svc Service
		if err := tx.Where("name = ?", svcName).First(&svc).Error; err != nil {
			return err
		}
		association := tx.Model(svc).Association("Consumers")

		// delete association
		if err := association.Delete(&svc).Error; err != nil {
			return err
		}

		// if no associations left, delete the consumer
		if association.Count() == 0 {
			if err := tx.Delete(Consumer{ID: consumerID}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
