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

package impl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/api"
	adfDeploymentHelper "github.com/apache/incubator-devlake/plugins/adf_deployment/helper"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/models"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*AdfDeployment)(nil)
var _ plugin.PluginInit = (*AdfDeployment)(nil)
var _ plugin.PluginTask = (*AdfDeployment)(nil)
var _ plugin.PluginApi = (*AdfDeployment)(nil)
var _ plugin.PluginMigration = (*AdfDeployment)(nil)
var _ plugin.CloseablePluginTask = (*AdfDeployment)(nil)

type AdfDeployment struct{}

func (p AdfDeployment) Connection() interface{} {
	return &models.AdfConnection{}
}

func (p AdfDeployment) Scope() interface{} {
	return &models.AdfDeployment{}
}

func (p AdfDeployment) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.AdfConnection{},
		&models.AdfDeploymentRevision{},
		&models.AdfDeployment{},
	}
}

func (p AdfDeployment) Description() string {
	return "collect some AdfDeployment data"
}

func (p AdfDeployment) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p AdfDeployment) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectAdfDeploymentRevisionsMeta,
		tasks.ExtractAdfDeploymentRevisionsMeta,
	}
}

func (p AdfDeployment) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	logger := taskCtx.GetLogger()
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.AdfConnection{}
	fmt.Println("Connection ID -->", op.ConnectionId)
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	println("Credentials -->", connection.Credentials)
	strType := reflect.TypeOf(connection.Credentials)
	fmt.Println(strType) // Output: string
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get AdfDeployment connection by the given connection ID")
	}

	apiClient, err := tasks.NewAdfDeploymentApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get AdfDeployment API client instance")
	}

	var adfAPIClientParams map[string]interface{}
	// Define a struct to hold the JSON data

	// Convert the string to a JSON object
	errUnmarshall := json.Unmarshal([]byte(connection.Credentials), &adfAPIClientParams)

	if errUnmarshall != nil {
		return nil, errors.Default.Wrap(err, "Error unmarshalling message")
	}
	adfClient, err := adfDeploymentHelper.NewKubeApiClient(adfAPIClientParams)

	if err != nil {
		return nil, err
	}

	taskData := &tasks.AdfDeploymentTaskData{
		Options:       op,
		ApiClient:     apiClient,
		AdfAPIClient:  adfClient,
	}
	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}
	return taskData, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p AdfDeployment) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/adf_deployment"
}

func (p AdfDeployment) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p AdfDeployment) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"fallback-endpoint": { // This is a fallback endpoint for the UI to use when the plugin is not configured
			"POST": api.FallbackEndpoint,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		"connections/:connectionId/namespaces": {
			"GET": api.GetNameSpaces,
		},
		"connections/:connectionId/:namespace/deployments": {
			"GET": api.GetDeployments,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
	}
}

func (p AdfDeployment) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
}

func (p AdfDeployment) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p AdfDeployment) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.AdfDeploymentTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
