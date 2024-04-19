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

import { useEffect, useMemo, useState } from 'react';
import { PluginConfig, PluginType } from '@/plugins';
import { IconName } from '@blueprintjs/core';
import * as API from '@/pages/project/home/api';
import ChartMixedSvg from './assets/chart-mixed.svg';
import MapSvg from '@/images/icons/map.svg';
import DiagramSankey from '@/images/icons/diagram-sankey.svg';

type HomeItemType = {
  key: string;
  title: string;
  description: string;
  items: {
    key: string;
    iconUrl?: string | IconName;
    title: string;
    description: string;
    path: string;
    disabled?: boolean;
    isTarget?: boolean;
    tooltipText: string;
  }[];
};

type ProjectItem = {
  name: string;
};

export const useHome = () => {
  const [projects, setProjects] = useState<ProjectItem[]>([]);

  const getProjects = async () => {
    const res = await API.getProjects({ page: 1, pageSize: 200 });
    setProjects(
      res.projects.map((it: any) => ({
        name: it.name,
      })),
    );
  };

  useEffect(() => {
    getProjects();
  }, []);

  const getGrafanaUrl = () => {
    const suffix = '/d/lCO8w-pVk/homepage?orgId=1';
    const { protocol, hostname } = window.location;

    return import.meta.env.DEV ? `${protocol}//${hostname}:3002${suffix}` : `/grafana${suffix}`;
  };
  const homeItems: HomeItemType[] = [
    {
      key: 'dashboards',
      title: 'Dashboards',
      description: 'Access your Grafana dashboards based on your projects',
      items: [
        {
          key: 'dashboards',
          iconUrl: ChartMixedSvg,
          title: 'Dashboards',
          description: 'Access your Dashboards and Metrics.',
          path: getGrafanaUrl(),
          isTarget: true,
          tooltipText: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
        },
      ],
    },
    {
      key: 'connection',
      title: 'Connections',
      description:
        'Connections are available for data collection and you can use webhooks to import deployments and incidents from the unsupported data integrations to calculate DORA metrics, etc. Please note: webhooks cannot be created or managed in Blueprints.',
      items: PluginConfig.filter((p) => [PluginType.Connection, PluginType.Incoming_Connection].includes(p.type)).map(
        (it) => ({
          key: it.plugin,
          iconUrl: it.icon,
          title: it.name,
          description: it.description || '',
          path: `/connections/${it.plugin}`,
          tooltipText: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
        }),
      ),
    },
    {
      key: 'projects',
      title: 'Projects',
      description: 'Access your projects and manage them.',
      items: projects.map(({ name }) => ({
        key: name,
        title: name,
        description: 'Access and manage your project settings.',
        path: `/projects/${name}`,
        tooltipText: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
      })),
    },
    {
      key: 'advanced',
      title: 'Advanced',
      description: 'Is possible to have advanced features with MetricSense',
      items: [
        {
          key: 'blueprints',
          iconUrl: MapSvg,
          title: 'Blueprints',
          description: 'Configure your blueprint settings.',
          path: '/blueprints',
          tooltipText: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
        },
        {
          key: 'pipelines',
          iconUrl: DiagramSankey,
          title: 'Pipelines',
          description: 'Create and Manage pipeline.',
          path: '/pipelines',
          disabled: true,
          tooltipText: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
        },
      ],
    },
  ];

  return useMemo(
    () => ({
      homeItems,
    }),
    [projects],
  );
};
