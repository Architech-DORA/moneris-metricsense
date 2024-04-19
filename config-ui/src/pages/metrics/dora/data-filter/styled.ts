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

export const Container = styled.div`
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
`;

export const ContentSection = styled.div`
  display: flex;
  flex-direction: column;
`;

export const SectionLabel = styled.div`
  margin-bottom: 10px;
`;

export const ApplyButtonContainer = styled.div``;

export const SelectionContainer = styled.div`
  min-width: 300px;
  max-width: 300px;
  .bp5-popover-target,
  .bp5-fill {
    height: 100%;
  }
  .selection-button {
    width: 100%;
    display: flex !important;
    justify-content: space-between !important;
  }
  .bp5-popover-content {
    background: #fff !important;
  }

  /* Selected item remove icon */
  .bp5-tag-remove {
  }
  .bp5-icon {
    display: flex;
  }
`;

/**
 * Date Range Filter
 */
export const DateRangeFilterContainer = styled.div`
  .bp5-date-range-input {
    display: flex;
  }
  .bp5-input-group {
    margin-right: 10px;
  }
`;
