/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	adfDeploymentHelper "github.com/apache/incubator-devlake/plugins/adf_deployment/helper"
)

type AdfDeploymentApiParams struct {
	ConnectionId   uint64 `json:"connectionId"`
	Namespace      string `json:"namespace" mapstructure:"namespace,omitempty"`
	DeploymentName string `json:"deploymentName" mapstructure:"deploymentName,omitempty"`
}

type AdfDeploymentOptions struct {
	// TODO add some custom options here if necessary
	// options means some custom params required by plugin running.
	// Such As How many rows do your want
	// You can use it in subtasks, and you need to pass it to main.go and pipelines.
	ConnectionId     uint64   `json:"connectionId"`
	Tasks            []string `json:"tasks,omitempty"`
	CreatedDateAfter string   `json:"createdDateAfter" mapstructure:"createdDateAfter,omitempty"`
	Namespace        string   `json:"namespace" mapstructure:"namespace,omitempty"`
	DeploymentName   string   `json:"deploymentName" mapstructure:"deploymentName,omitempty"`
}

type AdfDeploymentTaskData struct {
	Options          *AdfDeploymentOptions
	ApiClient        *helper.ApiAsyncClient
	AdfAPIClient    *adfDeploymentHelper.AdfApiClientSet
	CreatedDateAfter *time.Time
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*AdfDeploymentOptions, errors.Error) {
	var op AdfDeploymentOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}
	return &op, nil
}

// NOT BEING USED
func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*helper.RawDataSubTaskArgs, *AdfDeploymentTaskData) {
	data := taskCtx.GetData().(*AdfDeploymentTaskData)
	filteredData := *data
	filteredData.Options = &AdfDeploymentOptions{}
	*filteredData.Options = *data.Options
	var params = AdfDeploymentApiParams{
		ConnectionId: data.Options.ConnectionId,
	}
	rawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}
