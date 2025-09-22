package main

import (
	"context"

	"github.com/go-gorote/auth"
	"github.com/go-gorote/gorote"
	"github.com/go-gorote/temp/app/example"
	"github.com/go-gorote/temp/env"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
)

func Config(ctx context.Context, app *fiber.App) error {
	sql, err := gorote.NewGorm(postgres.Open(env.Sql01.DSN()))
	if err != nil {
		return err
	}

	coreRouter, err := auth.New(&auth.Config{
		DB:               sql,
		AppName:          env.App.Name,
		PrivateKey:       env.PrivateKey,
		JwtExpireAccess:  env.JwtExpireAccess,
		JwtExpireRefresh: env.JwtExpireRefresh,
		SuperEmail:       env.SuperEmail,
		SuperPass:        env.SuperPassword,
		Domain:           env.Domain,
	})
	if err != nil {
		return err
	}

	coreRouter.RegisterRouter(app.Group("/api/v1/core"))

	exampleRouter, err := example.New(&example.Config{
		DB:        sql,
		PublicKey: env.PublicKey,
	})
	if err != nil {
		return err
	}
	exampleRouter.RegisterRouter(app.Group("/api/v1/example"))
	return nil
}
