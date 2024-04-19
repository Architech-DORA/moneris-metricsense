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

import { Card } from '@blueprintjs/core';
import styled from 'styled-components';

export const StyledH1 = styled.h1`
  font-size: 20px;
  margin-bottom: 16px;
  color: #111111;
`;

export const Description = styled.p`
  font-size: 14px;
  margin: 0 0 16px 0;
  color: #646464;
`;

export const StyledH2 = styled.h2`
  font-size: 16px;
`;

export const CardWrapper = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 32px;
`;

export const StyledCard = styled(Card)<{ disabled: boolean }>`
  border-radius: 4px;
  width: 200px;
  min-height: 100px;
  padding: 16px;
  background-color: ${(props) => (props.disabled ? 'lightgrey' : '#ffffff')};
  box-shadow: ${(props) => (props.disabled ? 'none' : '0px 2px 4px 0px #0000001a')};
  pointer-events: ${(props) => (props.disabled ? 'none' : 'all')};
`;

export const TitleWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  align-items: baseline;
`;

export const LeftTitleWrapper = styled.div`
  display: flex;
`;

export const StyledLeftSvgIcon = styled.div<{ url?: string }>`
  background-color: #4ac1bf;
  mask: url(${(props) => props.url}) no-repeat center / contain;
  -webkit-mask: url(${(props) => props.url}) no-repeat center / contain;
  margin-right: 12px;

  img {
    opacity: 0;
    width: 21px;
    height: 18px;
  }
`;

export const StyledRightSvgIcon = styled.div<{ url: string }>`
  background-color: #1c1c1c;
  mask: url(${(props) => props.url}) no-repeat center / contain;
  -webkit-mask: url(${(props) => props.url}) no-repeat center / contain;
  padding: 0 10px;
  height: 15px;
  img {
    opacity: 0;
    width: 3px;
    height: 3px;
  }
`;

export const StyledP = styled.p`
  color: #8a8a8a;
  font-size: 12px;
  margin: 0;
`;
