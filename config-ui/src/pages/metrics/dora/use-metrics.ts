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
import { useState, useMemo } from 'react';

import * as API from './api';

export type FetchParams = {
  projects: string;
  dateFrom: string;
  dateTo: string;
};

const BENCH_MARKS = {
  DEPLOYMENT_FREQUENCY: {
    'Fewer than once per six month': 'LOW',
    'Between once per month and once every 6 months': 'MEDIUM',
    'Between once per week and once per month': 'HIGH',
    'On-demand': 'ELITE',
  },
  LEAD_TIME_FOR_CHANGES: {
    'More than six months': 'LOW',
    'Between one week and six months': 'MEDIUM',
    'Less than one week': 'HIGH',
    'Less than one hour': 'ELITE',
  },
  TIME_TO_RESTORE_SERVICE: {
    'More than one week': 'LOW',
    'Between one day and one week': 'MEDIUM',
    'Less than one day': 'HIGH',
    'Less than one hour': 'ELITE',
  },
  CHANGE_FAILURE_RATE: {
    '> 30%': 'LOW',
    '21%-30%': 'MEDIUM',
    '16%-20%': 'HIGH',
    '0-15%': 'ELITE',
  },
} as const;

const BENCH_MARK_COLOR_SCHEME = {
  LOW: '#fe4e35',
  MEDIUM: '#fbd606',
  HIGH: '#28bf58',
  ELITE: '#007bbf',
};

export const BENCH_MARK_COLOR_TO_TEXT = {
  [BENCH_MARK_COLOR_SCHEME.LOW]: 'Low',
  [BENCH_MARK_COLOR_SCHEME.MEDIUM]: 'Medium',
  [BENCH_MARK_COLOR_SCHEME.HIGH]: 'High',
  [BENCH_MARK_COLOR_SCHEME.ELITE]: 'Elite',
};

export const useGetDeploymentFrequency = () => {
  const [loading, setLoading] = useState(false);
  const [metricData, setMetricData] = useState<API.GetDeploymentFrequencyResponse>({
    deploymentFrequency: '',
    chartData: [],
  });

  const getDeploymentFrequency = async (params: FetchParams) => {
    setLoading(true);
    try {
      const res = await API.getDeploymentFrequency(params);
      setMetricData(res);
    } finally {
      setLoading(false);
    }
  };

  const benchMark =
    BENCH_MARKS.DEPLOYMENT_FREQUENCY[metricData?.deploymentFrequency as keyof typeof BENCH_MARKS.DEPLOYMENT_FREQUENCY];

  const fetchData = async (params: FetchParams) => {
    await getDeploymentFrequency(params);
  };

  return useMemo(
    () => ({
      fetchData,
      loading,
      metricData,
      benchMark,
      benchMarkColor: BENCH_MARK_COLOR_SCHEME[benchMark],
    }),
    [loading, metricData, benchMark],
  );
};

export const useGetLeadTimeForChanges = () => {
  const [loading, setLoading] = useState(false);
  const [metricData, setMetricData] = useState<API.GetLeadTimeForChangesResponse>({
    leadTimeForChanges: '',
    chartData: [],
  });

  const getData = async (params: FetchParams) => {
    setLoading(true);
    try {
      const res = await API.getLeadTimeForChanges(params);
      let chartData = [];
      if (res.chartData && res.chartData.length) {
        chartData = res.chartData.map((item: any) => {
          return {
            ...item,
            medianChangeLeadTimeInHour: +(+item.medianChangeLeadTimeInHour || 0).toFixed(2),
          };
        });
      }
      setMetricData({ ...res, chartData });
    } finally {
      setLoading(false);
    }
  };

  const benchMark =
    BENCH_MARKS.LEAD_TIME_FOR_CHANGES[metricData?.leadTimeForChanges as keyof typeof BENCH_MARKS.LEAD_TIME_FOR_CHANGES];

  const fetchData = async (params: FetchParams) => {
    await getData(params);
  };

  return useMemo(
    () => ({
      fetchData,
      loading,
      metricData,
      benchMark,
      benchMarkColor: BENCH_MARK_COLOR_SCHEME[benchMark],
    }),
    [loading, metricData, benchMark],
  );
};

export const useGetTimeToRestoreService = () => {
  const [loading, setLoading] = useState(false);
  const [metricData, setMetricData] = useState<API.GetTimeToRestoreServiceResponse>({
    timeToRestoreService: '',
    chartData: [],
  });

  const getData = async (params: FetchParams) => {
    setLoading(true);
    try {
      const res = await API.GetTimeToRestoreService(params);
      let chartData = [];
      if (res.chartData && res.chartData.length) {
        chartData = res.chartData.map((item: any) => {
          return {
            ...item,
            medianTimeToResolveInHour: +(+item.medianTimeToResolveInHour || 0).toFixed(2),
          };
        });
      }
      setMetricData({ ...res, chartData });
    } finally {
      setLoading(false);
    }
  };

  const benchMark =
    BENCH_MARKS.TIME_TO_RESTORE_SERVICE[
      metricData?.timeToRestoreService as keyof typeof BENCH_MARKS.TIME_TO_RESTORE_SERVICE
    ];

  const fetchData = async (params: FetchParams) => {
    await getData(params);
  };

  return useMemo(
    () => ({
      fetchData,
      loading,
      metricData,
      benchMark,
      benchMarkColor: BENCH_MARK_COLOR_SCHEME[benchMark],
    }),
    [loading, metricData, benchMark],
  );
};

export const useChangeFailureRate = () => {
  const [loading, setLoading] = useState(false);
  const [metricData, setMetricData] = useState<API.GetChangeFailureRateResponse>({
    changeFailureRate: '',
    chartData: [],
  });

  const getData = async (params: FetchParams) => {
    setLoading(true);
    try {
      const res = await API.GetChangeFailureRate(params);

      let chartData = [];
      if (res.chartData && res.chartData.length) {
        chartData = res.chartData.map((item: any) => {
          return {
            ...item,
            changeFailureRate: +(+item.changeFailureRate || 0).toFixed(2),
          };
        });
      }
      setMetricData({ ...res, chartData });
    } finally {
      setLoading(false);
    }
  };

  const benchMark =
    BENCH_MARKS.CHANGE_FAILURE_RATE[metricData?.changeFailureRate as keyof typeof BENCH_MARKS.CHANGE_FAILURE_RATE];

  const fetchData = async (params: FetchParams) => {
    await getData(params);
  };

  return useMemo(
    () => ({
      fetchData,
      loading,
      metricData,
      benchMark,
      benchMarkColor: BENCH_MARK_COLOR_SCHEME[benchMark],
    }),
    [loading, metricData, benchMark],
  );
};
