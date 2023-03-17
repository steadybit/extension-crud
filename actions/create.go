// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package actions

import (
	"encoding/json"
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/extension-crud/config"
	"github.com/steadybit/extension-crud/db"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extconversion"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extutil"
	"net/http"
)

func RegisterCreateAction() {
	exthttp.RegisterHttpHandler("/actions/create", exthttp.GetterAsHandler(getCreateDescription))
	exthttp.RegisterHttpHandler("/actions/create/prepare", prepareCreate)
	exthttp.RegisterHttpHandler("/actions/create/start", startCreate)
}

func getCreateDescription() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("%s.create", config.Config.TargetType),
		Label:       fmt.Sprintf("create %s", config.Config.TargetTypeLabel),
		Description: fmt.Sprintf("Stores a new %s entity within the CRUD extension's in-memory data store.", config.Config.TargetTypeLabel),
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Category:    extutil.Ptr("CRUD"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.Instantaneous,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:        "label",
				Label:       "Label",
				Description: extutil.Ptr(fmt.Sprintf("What label should the newly created %s carry?", config.Config.TargetTypeLabel)),
				Type:        action_kit_api.String,
				Order:       extutil.Ptr(1),
				Required:    extutil.Ptr(true),
			},
		},
		Prepare: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/create/prepare",
		},
		Start: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/create/start",
		},
	}
}

type CreateState struct {
	Label string
}

func prepareCreate(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.PrepareActionRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	label := request.Config["label"].(string)
	var convertedState action_kit_api.ActionState
	err = extconversion.Convert(CreateState{Label: label}, &convertedState)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to initialize create action state", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.PrepareResult{
		State: convertedState,
	})
}

func startCreate(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.ActionStatusRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	var state CreateState
	err = extconversion.Convert(request.State, &state)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to decode action state", err))
		return
	}

	err = db.Create(state.Label)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to persist in database", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Created '%s' labeled '%s' in instance '%s'", config.Config.TargetTypeLabel, state.Label, config.Config.InstanceName),
			},
		}),
	})
}
