//go:build wireinject
// +build wireinject

// Package di contains all the logic related to di
package di

import (
	"LiteraTest/create/v1/internal"
	"LiteraTest/create/v1/internal/uc"

	"github.com/google/wire"
)

var stdSet = wire.NewSet(
	newAWSSessionProvider,
	uc.NewFindDoubleBookedEventsUC,
	uc.NewParseEventsToUTCUC,
	internal.NewHandler,

	wire.Bind(new(internal.FindDoubleBookedEventsUCInterface), new(*uc.FindDoubleBookedEventsUC)),
	wire.Bind(new(internal.ParseEventsToUTCUCInterface), new(*uc.ParseEventsToUTCUC)),
)
