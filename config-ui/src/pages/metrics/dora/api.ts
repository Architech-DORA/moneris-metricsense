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
import { request } from '@/utils';

/**
 * Deployment Frequency
 */
type GetDeploymentFrequencyParams = {
  projects: string;
  dateFrom: string;
  dateTo: string;
};

export type ChartData = {
  dataCount: number;
  month: string;
};

export type GetDeploymentFrequencyResponse = {
  deploymentFrequency: string;
  chartData: ChartData[];
};

export const getDeploymentFrequency = (params: GetDeploymentFrequencyParams) =>
  request('/deployment-frequency', { data: params });

/**
 * Lead time for changes
 */
type GetLeadTimeForChangesParams = {
  projects: string;
  dateFrom: string;
  dateTo: string;
};

export type LeadTimeForChangesResponse = {
  medianChangeLeadTimeInHour: number;
  month: string;
};
export type GetLeadTimeForChangesResponse = {
  leadTimeForChanges: string;
  chartData: ChartData[];
};

export const getLeadTimeForChanges = (params: GetLeadTimeForChangesParams) =>
  request('/lead-time-for-changes', { data: params });

/**
 * Time to restore service
 */
type GetTimeToRestoreServiceParams = {
  projects: string;
  dateFrom: string;
  dateTo: string;
};

export type MedianTimeToRestoreServiceChartData = {
  medianTimeToResolveInHour: number;
  month: string;
};

export type GetTimeToRestoreServiceResponse = {
  timeToRestoreService: string;
  chartData: MedianTimeToRestoreServiceChartData[];
};

export const GetTimeToRestoreService = (params: GetTimeToRestoreServiceParams) =>
  request('/median-time-to-restore', { data: params });

/**
 * Change failure rate
 */
type GetChangeFailureRateParams = {
  projects: string;
  dateFrom: string;
  dateTo: string;
};

export type ChangeFailureRateChartData = {
  changeFailureRate: number;
  month: string;
};

export type GetChangeFailureRateResponse = {
  changeFailureRate: string;
  chartData: ChangeFailureRateChartData[];
};

export const GetChangeFailureRate = (params: GetChangeFailureRateParams) =>
  request('/change-failure-rate', { data: params });
