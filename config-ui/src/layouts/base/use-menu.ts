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
import { IconName } from '@blueprintjs/core';
import { PluginConfig, PluginType } from '@/plugins';
import * as API from '@/pages/project/home/api';
import HouseSvg from '@/images/icons/house.svg';
import MapSvg from '@/images/icons/map.svg';
import DiagramSanKeySvg from '@/images/icons/diagram-sankey.svg';
import { request } from '@/utils';

export type MenuItemType = {
  key: string;
  title?: string;
  icon?: IconName;
  iconUrl?: string;
  rightIcon?: IconName;
  path?: string;
  children?: MenuItemType[];
  target?: boolean;
  isBeta?: boolean;
  disabled?: boolean;
  onClick?: Function;
};

type ProjectItem = {
  name: string;
};

export const useMenu = () => {
  const [projects, setProjects] = useState<ProjectItem[]>([]);

  const getProjects = async () => {
    const res = await API.getProjects({ page: 1, pageSize: 200 });
    setProjects(
      res.projects.map((it: any) => ({
        name: it.name,
      })),
    );
  };

  const handleSignOut = async () => {
    request('/logout').then((response) => {
      if (response.status === 200 || response.status === 302) {
        window.location.href = response.headers.get('Location');
      } else {
        throw new Error('Logout failed');
      }
    });
  };

  useEffect(() => {
    getProjects();
  }, []);

  return useMemo(
    (): MenuItemType[] => [
      {
        key: 'home',
        children: [
          {
            key: 'home',
            title: 'Home',
            iconUrl: HouseSvg,
            path: '/home',
          },
        ],
      },
      {
        key: 'dora-metrics',
        children: [
          {
            key: 'dora-metrics',
            title: 'DORA Metrics',
            icon: 'chart',
            path: '/dora-metrics',
          },
        ],
      },
      {
        key: 'connection',
        title: 'Connections',
        icon: 'data-connection',
        path: '/connections',
        children: PluginConfig.filter((p) =>
          [PluginType.Connection, PluginType.Incoming_Connection].includes(p.type),
        ).map((it) => ({
          key: it.plugin,
          title: it.name,
          iconUrl: it.icon,
          path: `/connections/${it.plugin}`,
          isBeta: it.isBeta,
        })),
      },
      {
        key: 'project',
        title: 'Projects',
        icon: 'home',
        path: '/projects',
        children: projects.map(({ name }) => ({
          key: name,
          title: name,
          rightIcon: 'chevron-right',
          path: `/projects/${name}`,
        })),
      },
      {
        key: 'advanced',
        title: 'Advanced',
        icon: 'pulse',
        path: '/advanced',
        children: [
          {
            key: 'blueprints',
            title: 'Blueprints',
            iconUrl: MapSvg,
            path: '/blueprints',
          },
          {
            key: 'pipelines',
            title: 'Pipelines',
            iconUrl: DiagramSanKeySvg,
            path: '/pipelines',
            disabled: true,
          },
        ],
      },
      {
        key: 'signout',
        children: [
          {
            key: 'signout',
            title: 'Sign Out',
            icon: 'log-out',
            onClick: handleSignOut,
          },
        ],
      },
    ],
    [projects],
  );
};
