package example

import (
	"crypto/rsa"

	"gorm.io/gorm"
)

type Config struct {
	*gorm.DB
	*rsa.PublicKey
}

func (c *Config) db() *gorm.DB {
	return c.DB
}

func (c *Config) publicKeyRSA() *rsa.PublicKey {
	return c.PublicKey
}

type configLoad interface {
	db() *gorm.DB
	publicKeyRSA() *rsa.PublicKey
}

type appRouter struct {
	publicKey  *rsa.PublicKey
	controller controller
}

type appController struct {
	service servicer
}

type appService struct {
	configLoad
}

func New(config configLoad) (*appRouter, error) {
	if err := config.db().AutoMigrate(); err != nil {
		return nil, err
	}

	service := appService{config}

	controller := appController{
		service: &service,
	}

	router := appRouter{
		publicKey:  config.publicKeyRSA(),
		controller: &controller,
	}

	return &router, nil
}
