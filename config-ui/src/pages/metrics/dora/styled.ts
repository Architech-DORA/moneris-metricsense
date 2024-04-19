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
import styled from 'styled-components';
import { Card as bpCard } from '@blueprintjs/core';

export const MetricsContainer = styled.div`
  margin-bottom: 50px;
`;

export const Wrapper = styled.div`
  display: flex;
  flex-wrap: wrap;
`;

export const MetricHeading = styled.div`
  padding-top: 20px;
`;

export const Card = styled(bpCard)`
  cursor: unset;
  display: flex;
  flex-direction: column;
  width: 46%;
  height: 380px;
  margin: 1%;
  padding: 30px;
  border-radius: 4px;
  box-shadow: 0px 2px 4px 0px #0000001a;
  &:nth-child(odd) {
    margin-left: 1px !important;
  }
  .header {
    flex: 0 1 auto;
    .bp5-heading {
      color: #202226;
      margin-bottom: 8px;
    }
  }
  .body {
    display: flex;
    align-items: start;
    justify-content: start;
    flex: 1 1 auto;

    .summary-container {
      display: flex;
      align-items: start;
      margin-bottom: 20px;
      .chip-indicator {
        display: flex;
        padding: 4px 12px;
        border-radius: 50px;
        color: #fff;
        font-weight: 500;
      }
      .summary {
        padding: 1px 0;
        display: flex;
        align-items: center;
        margin-left: 8px;
        color: #46484b;
      }
    }
  }
  .custom-tooltip {
    background-color: #fff;
    padding: 5px;
    border-radius: 3px;
  }
`;
