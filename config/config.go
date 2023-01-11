// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 Steadybit GmbH

package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var (
	Config Specification
)

type Specification struct {
	InstanceName    string `default:"Animal Shelter" split_words:"true"`
	TargetType      string `default:"dog" split_words:"true"`
	TargetTypeLabel string `default:"Dog" split_words:"true"`
}

func init() {
	err := envconfig.Process("steadybit_extension", &Config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse configuration from environment.")
	}
}
