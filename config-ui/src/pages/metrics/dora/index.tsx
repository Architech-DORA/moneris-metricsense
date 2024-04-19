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
import { useEffect, useState } from 'react';
import { Elevation } from '@blueprintjs/core';

import { AreaChart, Area, YAxis, XAxis, Tooltip, ResponsiveContainer } from 'recharts';

import {
  useGetDeploymentFrequency,
  useGetLeadTimeForChanges,
  useGetTimeToRestoreService,
  useChangeFailureRate,
  FetchParams,
} from './use-metrics';

import { BENCH_MARK_COLOR_TO_TEXT } from './use-metrics';

import * as Styled from './styled';
import DataFilter from './data-filter';
import { Project } from './data-filter/project-selection';
import * as ProjectAPI from '@/pages/project/home/api';

const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip">
        <div
          style={{ color: `${payload.color}` }}
        >{`${payload[0].payload?.month} - ${payload[0].payload?.dataCount}`}</div>
      </div>
    );
  }
  return null;
};

const MedianLeadTimeForChangesChartCustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip">
        <div
          style={{ color: `${payload.color}` }}
        >{`${payload[0].payload?.month} - ${payload[0].payload?.medianChangeLeadTimeInHour} hours`}</div>
      </div>
    );
  }
  return null;
};

const MedianTimeToRestoreServiceChartCustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip">
        <div
          style={{ color: `${payload.color}` }}
        >{`${payload[0].payload?.month} - ${payload[0].payload?.medianTimeToResolveInHour} hours`}</div>
      </div>
    );
  }
  return null;
};

const ChangeFailureRateChartCustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip">
        <div
          style={{ color: `${payload.color}` }}
        >{`${payload[0].payload?.month} - ${payload[0].payload?.changeFailureRate}`}</div>
      </div>
    );
  }
  return null;
};

export const DORAPage = () => {
  const [chartFilter, setChartFilter] = useState<{
    projects?: string[];
    dateFrom?: string;
    dateTo?: string;
  }>({
    projects: undefined,
    dateFrom: undefined,
    dateTo: undefined,
  });

  const [projects, setProjects] = useState<Project[]>([]);

  const deploymentFrequency = useGetDeploymentFrequency();
  const leadTimeForChanges = useGetLeadTimeForChanges();
  const timeToRestoreService = useGetTimeToRestoreService();
  const changeFailureRate = useChangeFailureRate();

  const onDateFilterChange = (timestampFrom: string, timeStampTo: string) => {
    setChartFilter({
      ...chartFilter,
      dateFrom: timestampFrom,
      dateTo: timeStampTo,
    });
  };

  const onChangeProjects = (projects: Project[]) => {
    setChartFilter({
      ...chartFilter,
      projects: projects.map((item) => item.name),
    });
  };

  const onApply = async () => {
    if (
      typeof chartFilter.projects !== 'undefined' &&
      typeof chartFilter.dateFrom !== 'undefined' &&
      typeof chartFilter.dateTo !== 'undefined'
    ) {
      const params = {
        ...chartFilter,
        projects: (chartFilter.projects || []).map((i) => `"${i}"`).join(),
      } as FetchParams;

      await deploymentFrequency.fetchData(params);
      await leadTimeForChanges.fetchData(params);
      await timeToRestoreService.fetchData(params);
      await changeFailureRate.fetchData(params);
    }
  };

  useEffect(() => {
    (async () => {
      const res = await ProjectAPI.getProjects({ page: 1, pageSize: 200 });
      setProjects(
        res.projects.map((it: any) => ({
          name: it.name,
        })),
      );
    })();
  }, []);

  return (
    <Styled.MetricsContainer>
      <DataFilter
        onDateFilterChange={onDateFilterChange}
        items={projects}
        onChangeProjects={onChangeProjects}
        onApply={onApply}
      />
      <Styled.MetricHeading>
        <h1>DORA Metrics</h1>
      </Styled.MetricHeading>
      <div>
        <Styled.Wrapper>
          <Styled.Card interactive={false} elevation={Elevation.TWO}>
            <div className="header">
              <h3 className="bp5-heading">Deployment frequency</h3>
            </div>
            <div className="body">
              {deploymentFrequency.loading && 'Loading...'}
              <div className="summary-container bp5-ui-text">
                <div
                  className="chip-indicator bp5-text-small"
                  style={{ backgroundColor: deploymentFrequency.benchMarkColor }}
                >
                  {BENCH_MARK_COLOR_TO_TEXT[deploymentFrequency.benchMarkColor]}
                </div>
                <div className="summary">{deploymentFrequency.metricData?.deploymentFrequency}</div>
              </div>
            </div>
            <ResponsiveContainer style={{ marginLeft: '-25px' }} height={225}>
              <AreaChart
                data={deploymentFrequency.metricData.chartData}
                margin={{
                  top: 0,
                  right: 0,
                  left: 0,
                  bottom: 30,
                }}
              >
                <defs>
                  <linearGradient id="colorUv1" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor={`${deploymentFrequency.benchMarkColor}`} stopOpacity={0.8} />
                    <stop offset="95%" stopColor={`${deploymentFrequency.benchMarkColor}`} stopOpacity={0} />
                  </linearGradient>
                </defs>
                <YAxis tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <XAxis dataKey="month" tick={{ fill: '#7E7F81', fontSize: '12px' }} />

                <Tooltip content={<CustomTooltip />} />
                <Area
                  type="monotone"
                  dataKey="dataCount"
                  stroke={`${deploymentFrequency.benchMarkColor}`}
                  fillOpacity={1}
                  fill="url(#colorUv1)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Styled.Card>
          <Styled.Card interactive={false} elevation={Elevation.TWO}>
            <div className="header">
              <h3 className="bp5-heading">Lead time for changes</h3>
            </div>
            <div className="body">
              {leadTimeForChanges.loading && 'Loading...'}
              <div className="summary-container bp5-ui-text">
                <div
                  className="chip-indicator bp5-text-small"
                  style={{ backgroundColor: leadTimeForChanges.benchMarkColor }}
                >
                  {BENCH_MARK_COLOR_TO_TEXT[leadTimeForChanges.benchMarkColor]}
                </div>
                <div className="summary">{leadTimeForChanges.metricData?.leadTimeForChanges}</div>
              </div>
            </div>
            <ResponsiveContainer style={{ marginLeft: '-25px' }} height={225}>
              <AreaChart
                data={leadTimeForChanges.metricData.chartData}
                margin={{
                  top: 0,
                  right: 0,
                  left: 0,
                  bottom: 30,
                }}
              >
                <defs>
                  <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor={`${leadTimeForChanges.benchMarkColor}`} stopOpacity={0.8} />
                    <stop offset="95%" stopColor={`${leadTimeForChanges.benchMarkColor}`} stopOpacity={0} />
                  </linearGradient>
                </defs>
                <YAxis tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <XAxis dataKey="month" tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <Tooltip content={<MedianLeadTimeForChangesChartCustomTooltip />} />
                <Area
                  type="monotone"
                  dataKey="medianChangeLeadTimeInHour"
                  stroke={`${leadTimeForChanges.benchMarkColor}`}
                  fillOpacity={1}
                  fill="url(#colorUv)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Styled.Card>
          <Styled.Card interactive={false} elevation={Elevation.TWO}>
            <div className="header">
              <h3 className="bp5-heading">Median time to restore service</h3>
            </div>
            <div className="body">
              {timeToRestoreService.loading && 'Loading...'}
              <div className="summary-container bp5-ui-text">
                <div
                  className="chip-indicator bp5-text-small"
                  style={{ backgroundColor: timeToRestoreService.benchMarkColor }}
                >
                  {BENCH_MARK_COLOR_TO_TEXT[timeToRestoreService.benchMarkColor]}
                </div>
                <div className="summary">{timeToRestoreService.metricData?.timeToRestoreService}</div>
              </div>
            </div>
            <ResponsiveContainer style={{ marginLeft: '-25px' }} height={225}>
              <AreaChart
                data={timeToRestoreService.metricData.chartData}
                margin={{
                  top: 0,
                  right: 0,
                  left: 0,
                  bottom: 30,
                }}
              >
                <defs>
                  <linearGradient id="colorUv3" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor={`${timeToRestoreService.benchMarkColor}`} stopOpacity={0.8} />
                    <stop offset="95%" stopColor={`${timeToRestoreService.benchMarkColor}`} stopOpacity={0} />
                  </linearGradient>
                </defs>
                <YAxis tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <XAxis dataKey="month" tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <Tooltip content={<MedianTimeToRestoreServiceChartCustomTooltip />} />
                <Area
                  type="monotone"
                  dataKey="medianTimeToResolveInHour"
                  stroke={`${timeToRestoreService.benchMarkColor}`}
                  fillOpacity={1}
                  fill="url(#colorUv1)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Styled.Card>
          <Styled.Card interactive={false} elevation={Elevation.TWO}>
            <div className="header">
              <h3 className="bp5-heading">Change failure rate</h3>
            </div>
            <div className="body">
              {changeFailureRate.loading && 'Loading...'}
              <div className="summary-container bp5-ui-text">
                <div
                  className="chip-indicator bp5-text-small"
                  style={{ backgroundColor: changeFailureRate.benchMarkColor }}
                >
                  {BENCH_MARK_COLOR_TO_TEXT[changeFailureRate.benchMarkColor]}
                </div>
                <div className="summary">{changeFailureRate.metricData?.changeFailureRate}</div>
              </div>
            </div>
            <ResponsiveContainer style={{ marginLeft: '-25px' }} height={225}>
              <AreaChart
                data={changeFailureRate.metricData.chartData}
                margin={{
                  top: 0,
                  right: 0,
                  left: 0,
                  bottom: 30,
                }}
              >
                <defs>
                  <linearGradient id="colorUv3" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor={`${changeFailureRate.benchMarkColor}`} stopOpacity={0.8} />
                    <stop offset="95%" stopColor={`${changeFailureRate.benchMarkColor}`} stopOpacity={0} />
                  </linearGradient>
                </defs>
                <YAxis tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <XAxis dataKey="month" tick={{ fill: '#7E7F81', fontSize: '12px' }} />
                <Tooltip content={<ChangeFailureRateChartCustomTooltip />} />
                <Area
                  type="monotone"
                  dataKey="changeFailureRate"
                  stroke={`${changeFailureRate.benchMarkColor}`}
                  fillOpacity={1}
                  fill="url(#colorUv1)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Styled.Card>
        </Styled.Wrapper>
      </div>
    </Styled.MetricsContainer>
  );
};
