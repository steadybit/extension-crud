// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-crud/config"
	"golang.org/x/exp/slices"
	"sync"
)

var (
	targets = make([]discovery_kit_api.Target, 0, 10)
	mu      sync.Mutex
)

func Create(label string) error {
	mu.Lock()
	defer mu.Unlock()

	for _, target := range targets {
		if target.Label == label {
			return fmt.Errorf("a %s with label '%s' already exists", config.Config.TargetTypeLabel, label)
		}
	}

	targets = append(targets, createTarget(uuid.New().String(), label))
	return nil
}

func Read() []discovery_kit_api.Target {
	return targets
}

func Update(oldLabel string, newLabel string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, target := range targets {
		if target.Label == oldLabel {
			targets[i] = createTarget(target.Id, newLabel)
			return nil
		}
	}

	return fmt.Errorf("a %s with label '%s' could not be found", config.Config.TargetTypeLabel, oldLabel)
}

func Delete(label string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, target := range targets {
		if target.Label == label {
			targets = slices.Delete(targets, i, i+1)
			return nil
		}
	}

	return fmt.Errorf("a %s with label '%s' could not be found", config.Config.TargetTypeLabel, label)
}

func createTarget(id string, label string) discovery_kit_api.Target {
	return discovery_kit_api.Target{
		Id:         label,
		Label:      label,
		TargetType: config.Config.TargetType,
		Attributes: map[string][]string{
			"crud.instance.name": {config.Config.InstanceName},
		},
	}
}
