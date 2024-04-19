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

import { Button, Icon, Menu, MenuItem, Navbar } from '@blueprintjs/core';
import styled from 'styled-components';

export const Wrapper = styled.div`
  display: flex;
  height: 100vh;
  background-color: #f3f3f3;
  overflow: hidden;
  width: 100%;
`;

export const Sider = styled.div`
  flex: 0 0 250px;
  position: relative;
  padding: 32px 16px;
  width: 250px;
  background-color: #202226;
  overflow-y: scroll;
  -ms-overflow-style: none;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }

  .menu {
    margin: 20px 4px 0;
    color: #fff;
    background-color: transparent;

    .menu-item,
    .sub-menu-item {
      display: flex;
      align-items: center;
      margin: 2px 0;
      line-height: 26px;
      transition: all 0.3s ease;
      border-radius: 8px;
      outline: none;
      cursor: pointer;

      &:hover {
        background-color: rgba(167, 182, 194, 0.3);
      }

      .bp4-icon {
        svg {
          width: 12px;
          height: 12px;
        }
      }
    }

    .sub-menu-item {
      border-radius: 3px;
    }
  }

  .copyright {
    position: absolute;
    right: 0;
    bottom: 30px;
    left: 0;
    text-align: center;
    color: rgba(124, 124, 124, 0.7);
    padding: 0 20px;
    .version {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
`;

export const Main = styled.div`
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex: auto;
`;

export const Header = styled(Navbar)`
  flex: 0 0 50px;
  background-color: #e6f5f5;
  box-shadow: none;

  a {
    display: flex;
    align-items: center;

    img {
      margin-right: 4px;
      width: 16px;
    }

    span {
      font-size: 12px;
    }
  }
`;

export const Inner = styled.div`
  flex: auto;
  overflow: auto;
  padding: 32px;
`;

export const Content = styled.div`
  /* margin: 0 auto; */
  max-width: 1200px;
  min-width: 900px;
`;

export const SiderMenuItem = styled.div`
  display: flex;
  align-items: center;

  & > .bp4-tag {
    margin-left: 8px;
  }
`;

export const DashboardIcon = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 8px 12px;
  border: 1px solid #7497f7;
  border-radius: 4px;
`;

export const StyledMenu = styled(Menu)`
  background-color: transparent;
  padding: 0;
  margin-top: 32px;
`;

export const MenuItemWrapper = styled.div`
  margin-bottom: 20px;
`;

export const StyledMenuItem = styled(MenuItem)<{ active: boolean }>`
  color: ${(props) => (props.active ? '#4AC1BF' : '#FFFFFF')} !important;
  background-color: transparent !important;
  text-transform: uppercase;
  font-size: 10px;
  outline: none;
  margin-bottom: 8px;
  padding: 8px;
  border-radius: 4px;

  &:hover {
    background-color: #434750 !important;
  }
`;

export const StyledSubMenuItem = styled(MenuItem)<{ active: boolean }>`
  color: ${(props) => (props.active ? '#4AC1BF' : '#FFFFFF')} !important;
  background-color: transparent !important;
  font-size: 14px;
  outline: none;
  margin-bottom: 8px;
  padding: 8px;
  border-radius: 4px;

  &:hover {
    background-color: #434750 !important;
  }
`;

export const StyledIcon = styled(Icon)<{ active: boolean }>`
  color: ${(props) => (props.active ? '#4AC1BF' : '#FFFFFF')} !important;
`;

export const StyledSvgIcon = styled.div<{ active: boolean; url?: string }>`
  background-color: ${(props) => (props.active ? '#4AC1BF' : '#FFFFFF')};
  mask: url(${(props) => props.url}) no-repeat center / contain;
  -webkit-mask: url(${(props) => props.url}) no-repeat center / contain;
  margin-right: 5px;

  img {
    opacity: 0;
    width: 16px;
    height: 14px;
  }
`;

export const StyledButton = styled(Button)`
  color: #ffffff !important;
  background-color: transparent !important;
  border: 1px solid #404652 !important;
  border-radius: 4px;
  padding: 12px;
  width: 180px;
  outline: none;
`;

export const StyledPopoverMenu = styled(Menu)`
  background-color: #42464e;
`;

export const StyledPopoverMenuItem = styled(MenuItem)`
  color: #ffffff;
  background-color: #42464e;
  padding: 10px;
  outline: none;

  &:hover {
    color: #ffffff;
    background-color: #31343a;
  }
`;

export const StyledProjectMenuItem = styled(MenuItem)`
  display: flex;
  align-items: center;
  border-radius: 0;
  outline: none;

  &:hover {
    background-color: #42464e;
  }
`;

export const StyledProjectMenuItemIcon = styled.div<{ url?: string }>`
  background-color: #ffffff;
  mask: url(${(props) => props.url}) no-repeat center / contain;
  -webkit-mask: url(${(props) => props.url}) no-repeat center / contain;

  img {
    opacity: 0;
  }
`;
export const StyledProjectTextWrapper = styled.div`
  display: flex;
  flex-direction: column;
  margin-left: 8px;

  span {
    display: inline-block;
  }

  & span:first-child {
    font-size: 10px;
  }

  & span:last-child {
    font-size: 16px;
    font-weight: bold;
  }
`;
