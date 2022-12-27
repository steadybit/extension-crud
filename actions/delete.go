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

func RegisterDeleteAction() {
	exthttp.RegisterHttpHandler("/actions/delete", exthttp.GetterAsHandler(getDeleteDescription))
	exthttp.RegisterHttpHandler("/actions/delete/prepare", prepareDelete)
	exthttp.RegisterHttpHandler("/actions/delete/start", startDelete)
}

func getDeleteDescription() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("%s.delete", config.Config.TargetType),
		Label:       fmt.Sprintf("delete %s", config.Config.TargetTypeLabel),
		Description: fmt.Sprintf("Removes a %s entity from the CRUD extension's in-memory data store.", config.Config.TargetTypeLabel),
		Version:     "1.0.0-SNAPSHOT",
		Category:    extutil.Ptr("CRUD"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.Instantaneous,
		Parameters:  []action_kit_api.ActionParameter{},
		TargetType:  extutil.Ptr(config.Config.TargetType),
		TargetSelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
			{
				Label: fmt.Sprintf("by %s label", config.Config.TargetTypeLabel),
				Query: "steadybit.label=\"\"",
			},
		}),
		Prepare: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/delete/prepare",
		},
		Start: action_kit_api.MutatingEndpointReference{
			Method: "POST",
			Path:   "/actions/delete/start",
		},
	}
}

type DeleteState struct {
	Label string
}

func prepareDelete(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.PrepareActionRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	label := request.Target.Name
	var convertedState action_kit_api.ActionState
	err = extconversion.Convert(DeleteState{Label: label}, &convertedState)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to initialize delete action state", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.PrepareResult{
		State: convertedState,
	})
}

func startDelete(w http.ResponseWriter, _ *http.Request, body []byte) {
	var request action_kit_api.ActionStatusRequestBody
	err := json.Unmarshal(body, &request)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to parse request body", err))
		return
	}

	var state DeleteState
	err = extconversion.Convert(request.State, &state)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to decode action state", err))
		return
	}

	err = db.Delete(state.Label)
	if err != nil {
		exthttp.WriteError(w, extension_kit.ToError("Failed to delete from database", err))
		return
	}

	exthttp.WriteBody(w, action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Deleted '%s' labeled '%s' from instance '%s'", config.Config.TargetTypeLabel, state.Label, config.Config.InstanceName),
			},
		}),
	})
}
