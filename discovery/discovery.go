// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package discovery

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-crud/config"
	"github.com/steadybit/extension-crud/db"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/exthttp"
	"net/http"
)

func RegisterDiscoveryHandlers() {
	exthttp.RegisterHttpHandler("/discovery-description", exthttp.GetterAsHandler(getDiscoveryDescription))
	exthttp.RegisterHttpHandler("/target-description", exthttp.GetterAsHandler(getTargetDescription))
	exthttp.RegisterHttpHandler("/attribute-descriptions", exthttp.GetterAsHandler(getAttributeDescriptions))
	exthttp.RegisterHttpHandler("/discover", discover)
}

func getDiscoveryDescription() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         config.Config.TargetType,
		RestrictTo: discovery_kit_api.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discover",
			CallInterval: discovery_kit_api.Ptr("5s"),
		},
	}
}

func getTargetDescription() discovery_kit_api.TargetDescription {
	return discovery_kit_api.TargetDescription{
		Id:       config.Config.TargetType,
		Label:    discovery_kit_api.PluralLabel{One: config.Config.TargetTypeLabel, Other: config.Config.TargetTypeLabel},
		Category: discovery_kit_api.Ptr("CRUD"),
		Version:  extbuild.GetSemverVersionStringOrUnknown(),
		Table: discovery_kit_api.Table{
			Columns: []discovery_kit_api.Column{
				{Attribute: "steadybit.label"},
				{Attribute: "crud.instance.name"},
			},
			OrderBy: []discovery_kit_api.OrderBy{
				{
					Attribute: "steadybit.label",
					Direction: "ASC",
				},
			},
		},
	}
}

func getAttributeDescriptions() discovery_kit_api.AttributeDescriptions {
	return discovery_kit_api.AttributeDescriptions{
		Attributes: []discovery_kit_api.AttributeDescription{
			{
				Attribute: "crud.instance.name",
				Label: discovery_kit_api.PluralLabel{
					One:   "CRUD instance name",
					Other: "CRUD instance names",
				},
			},
		},
	}
}

func discover(w http.ResponseWriter, _ *http.Request, _ []byte) {
	exthttp.WriteBody(w, discovery_kit_api.DiscoveredTargets{Targets: db.Read()})
}
