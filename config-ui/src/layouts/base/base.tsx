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

import { Intent, Tag } from '@blueprintjs/core';
import React from 'react';
import { useLocation } from 'react-router-dom';

import { Logo, PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { history } from '@/utils/history';

import * as API from './api';
import * as S from './styled';
import { MenuItemType, useMenu } from './use-menu';

interface Props {
  children: React.ReactNode;
}

export const BaseLayout = ({ children }: Props) => {
  const menu = useMenu();
  const { pathname } = useLocation();

  const { ready, data } = useRefreshData<{ version: string }>(() => API.getVersion(), []);

  // const token = window.localStorage.getItem('accessToken');

  const handlePushPath = (it: MenuItemType) => {
    if (it.onClick) {
      it.onClick(it);
    } else if (!it.target && it.path) {
      history.push(it.path);
    } else {
      window.open(it.path, '_blank');
    }
  };

  // const handleSignOut = () => {
  //   localStorage.removeItem(`accessToken`);
  //   history.push('/login');
  // };

  // const getGrafanaUrl = () => {
  //   const suffix = '/d/lCO8w-pVk/homepage?orgId=1';
  //   const { protocol, hostname } = window.location;

  //   return import.meta.env.DEV ? `${protocol}//${hostname}:3002${suffix}` : `/grafana${suffix}`;
  // };

  if (!ready || !data) {
    return <PageLoading />;
  }

  return (
    <S.Wrapper>
      <S.Sider>
        <Logo />
        <S.StyledMenu>
          {menu.map((it, i) => {
            const paths = [it.path, ...(it.children ?? []).map((cit) => cit.path)];
            const active = paths.indexOf(pathname) >= 0;

            return (
              <S.MenuItemWrapper key={it.key}>
                {it.title && (
                  <S.StyledMenuItem key={it.key} text={it.title} active={active} onClick={() => handlePushPath(it)} />
                )}
                {it.children?.map((cit) => {
                  const subMenuItemActive = pathname === cit.path;

                  return (
                    <S.StyledSubMenuItem
                      key={cit.key}
                      text={
                        <S.SiderMenuItem>
                          <span>{cit.title}</span>
                          {cit.isBeta && <Tag intent={Intent.WARNING}>beta</Tag>}
                        </S.SiderMenuItem>
                      }
                      icon={
                        cit.icon ? (
                          <S.StyledIcon icon={cit.icon} active={subMenuItemActive} />
                        ) : cit.iconUrl ? (
                          <S.StyledSvgIcon active={subMenuItemActive} url={cit.iconUrl}>
                            <img src={cit.iconUrl} alt={cit.title} />
                          </S.StyledSvgIcon>
                        ) : null
                      }
                      labelElement={
                        cit.rightIcon ? <S.StyledIcon icon={cit.rightIcon} active={subMenuItemActive} /> : null
                      }
                      active={subMenuItemActive}
                      disabled={cit.disabled}
                      onClick={() => handlePushPath(cit)}
                    />
                  );
                })}
              </S.MenuItemWrapper>
            );
          })}
        </S.StyledMenu>
        {/*<div className="copyright">*/}
        {/*  <div>Apache 2.0 License</div>*/}
        {/*  <div className="version">{data.version}</div>*/}
        {/*</div>*/}
      </S.Sider>
      <S.Main>
        {/* <S.Header>
          <Navbar.Group align={Alignment.RIGHT}>
            <S.DashboardIcon>
              <ExternalLink link={getGrafanaUrl()}>
                <img src={DashboardIcon} alt="dashboards" />
                <span>Dashboards</span>
              </ExternalLink>
            </S.DashboardIcon>
            <Navbar.Divider />
            {token && (
              <>
                <Navbar.Divider />
                <Button small intent={Intent.NONE} onClick={handleSignOut}>
                  Sign Out
                </Button>
              </>
            )}
          </Navbar.Group>
        </S.Header> */}
        <S.Inner>
          <S.Content>{children}</S.Content>
        </S.Inner>
      </S.Main>
    </S.Wrapper>
  );
};
