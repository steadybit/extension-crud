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
	"github.com/steadybit/extension-kit/extconversion"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extutil"
	"net/http"
)

func RegisterUpdateAction() {
	exthttp.RegisterHttpHandler("/actions/update", exthttp.GetterAsHandler(getUpdateDescription))
	exthttp.RegisterHttpHandler("/actions/update/prepare", prepareUpdate)
	exthttp.RegisterHttpHandler("/actions/update/start", startUpdate)
}

func getUpdateDescription() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("%s.update", config.Config.TargetType),
		Label:       fmt.Sprintf("update %s", config.Config.TargetTypeLabel),
		Description: fmt.Sprintf("Renames a %s entity within the CRUD extension's in-memory data store.", config.Config.TargetTypeLabel),
		Version:     "1.0.0-SNAPSHOT",
		Category:    extutil.Ptr("CRUD"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.Instantaneous,
		TargetType:  extutil.Ptr(config.Config.TargetType),
		TargetSelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
			{
				Label: fmt.Sprintf("by %s label", config.Config.TargetTypeLabel),
				Query: "steadybit.label=\"\"",
			},
		}),
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:        "label",
				Label:       "New Label",
				Description: extutil.Ptr(fmt.Sprintf("What should the new label of this %s be?", config.Config.TargetTypeLabel)),
				Type:        action_kit_api.String,
				Order:       extutil.Ptr(1),
				Required:    extutil.Ptr(true),
			},
		},
		Prepare: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/update/prepare",
		},
		Start: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/update/start",
		},
	}
}

type UpdateState struct {
	OldLabel string
	NewLabel string
}

func prepareUpdate(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.PrepareActionRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	oldLabel := request.Target.Name
	newLabel := request.Config["label"].(string)
	var convertedState action_kit_api.ActionState
	err = extconversion.Convert(UpdateState{OldLabel: oldLabel, NewLabel: newLabel}, &convertedState)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to initialize update action state", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.PrepareResult{
		State: convertedState,
	})
}

func startUpdate(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.ActionStatusRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	var state UpdateState
	err = extconversion.Convert(request.State, &state)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to decode action state", err))
		return
	}

	err = db.Update(state.OldLabel, state.NewLabel)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to update in database", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Updated '%s' labeled '%s' (previously '%s') in instance '%s'", config.Config.TargetTypeLabel, state.NewLabel, state.OldLabel, config.Config.InstanceName),
			},
		}),
	})
}
