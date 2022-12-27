// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-crud/actions"
	"github.com/steadybit/extension-crud/config"
	"github.com/steadybit/extension-crud/discovery"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kong/utils"
	"net/http"
)

func main() {
	extlogging.InitZeroLog()

	utils.RegisterHttpHandler("/", utils.GetterAsHandler(getExtensionList))
	discovery.RegisterDiscoveryHandlers()
	actions.RegisterCreateAction()
	actions.RegisterUpdateAction()
	actions.RegisterDeleteAction()

	port := config.Config.Port
	log.Info().Msgf("Starting extension-crud server on port %d. Get started via /", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start extension-crud server on port %d", port)
	}
}

type ExtensionListResponse struct {
	Actions          []action_kit_api.DescribingEndpointReference    `json:"actions"`
	Discoveries      []discovery_kit_api.DescribingEndpointReference `json:"discoveries"`
	TargetTypes      []discovery_kit_api.DescribingEndpointReference `json:"targetTypes"`
	TargetAttributes []discovery_kit_api.DescribingEndpointReference `json:"targetAttributes"`
}

func getExtensionList() ExtensionListResponse {
	return ExtensionListResponse{
		Actions: []action_kit_api.DescribingEndpointReference{
			{
				"GET",
				"/actions/create",
			},
			{
				"GET",
				"/actions/update",
			},
			{
				"GET",
				"/actions/delete",
			},
		},
		Discoveries: []discovery_kit_api.DescribingEndpointReference{
			{
				"GET",
				"/discovery-description",
			},
		},
		TargetTypes: []discovery_kit_api.DescribingEndpointReference{
			{
				"GET",
				"/target-description",
			},
		},
		TargetAttributes: []discovery_kit_api.DescribingEndpointReference{
			{
				"GET",
				"/attribute-descriptions",
			},
		},
	}
}
