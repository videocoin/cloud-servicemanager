package datastore

import (
	"github.com/videocoin/common/models"
)

// Service represents a managed service. Ex: symphony.videocoin.network.
type Service struct {
	model.Base
	ID        string `gorm:"primary_key"`
	Name      string
	Consumers []*Consumer `gorm:"many2many:services_consumers"`
}

// TableName set Service's table name to be `services`.
func (svc *Service) TableName() string { return "services" }

// Consumer represents a VideoCoin Studio project.
type Consumer struct {
	model.Base
	ID       string    `gorm:"primary_key"`
	Services []Service `gorm:"many2many:services_consumers"`
}

// TableName set Consumer's table name to be `consumers`.
func (c *Consumer) TableName() string { return "consumers" }
