/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { BasePipelineConfig } from '@/plugins/register/base';
import { AEConfig } from './register/ae';
import { CustomizeConfig } from './register/customize';
import { DBTConfig } from './register/dbt';
import { DORAConfig } from './register/dora';
import { FeiShuConfig } from './register/feishu';
import { GiteeConfig } from './register/gitee';
import { GitExtractorConfig } from './register/gitextractor';
import { GitHubConfig } from './register/github';
import { GitHubGraphqlConfig } from './register/github_graphql';
import { JenkinsConfig } from './register/jenkins';
import { JiraConfig } from './register/jira';
import { KubeDeploymentConfig } from './register/myplug';
import { adfDeploymentConfig } from './register/azure_data_factory';
import { OrgConfig } from './register/org';
import { RefDiffConfig } from './register/refdiff';
import { StarRocksConfig } from './register/starrocks';
import { WebhookConfig } from './register/webook';
// import { AzureConfig } from './register/azure';
import type { PluginConfigType } from './types';

export const PluginConfig: PluginConfigType[] = [
  AEConfig,
  // AzureConfig,
  // BitBucketConfig,
  CustomizeConfig,
  DBTConfig,
  DORAConfig,
  FeiShuConfig,
  GiteeConfig,
  GitExtractorConfig,
  GitHubConfig,
  GitHubGraphqlConfig,
  // GitLabConfig,
  JenkinsConfig,
  JiraConfig,
  OrgConfig,
  // PagerDutyConfig,
  RefDiffConfig,
  // SonarQubeConfig,
  StarRocksConfig,
  // TAPDConfig,
  // TeambitionConfig,
  // ZenTaoConfig,
  WebhookConfig,
  KubeDeploymentConfig,
  adfDeploymentConfig
].sort((a, b) => a.sort - b.sort);

// get plugin config by plugin name
export function getPluginConfig(pluginName: string): PluginConfigType {
  let pluginConfig = PluginConfig.find((plugin) => plugin.plugin === pluginName) as PluginConfigType;
  if (!pluginConfig) {
    pluginConfig = {
      ...BasePipelineConfig,
      plugin: pluginName,
      name: pluginName,
    } as PluginConfigType;
  }
  return {
    ...pluginConfig,
    transformationType: pluginConfig.transformationType || (pluginConfig.transformation ? 'for-connection' : 'none'),
  };
}
