// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 Steadybit GmbH

package main

import (
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-crud/actions"
	"github.com/steadybit/extension-crud/discovery"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kong/utils"
)

func main() {
	extlogging.InitZeroLog()
	extbuild.PrintBuildInformation()

	utils.RegisterHttpHandler("/", utils.GetterAsHandler(getExtensionList))
	discovery.RegisterDiscoveryHandlers()
	actions.RegisterCreateAction()
	actions.RegisterUpdateAction()
	actions.RegisterDeleteAction()

	exthttp.Listen(exthttp.ListenOpts{
		Port: 8091,
	})
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
