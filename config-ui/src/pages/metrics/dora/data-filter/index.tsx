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
import { Button } from '@blueprintjs/core';
import DateRangeFilter, { DateRangeFilterProps } from './date-range-filter';
import ProjectSelection, { BaseProps as ProjectSelectionBaseProps } from './project-selection';

import * as Styled from './styled';

type DataFilterProps = {
  onApply: () => void;
} & DateRangeFilterProps &
  ProjectSelectionBaseProps;

const DataFilter = (props: DataFilterProps) => {
  return (
    <Styled.Container>
      {/*
        Project Filter
      */}
      <Styled.ContentSection>
        <Styled.SectionLabel>Projects</Styled.SectionLabel>
        <ProjectSelection items={props.items} onChangeProjects={props.onChangeProjects} />
      </Styled.ContentSection>
      {/*
        Date Range Filter
      */}
      <Styled.ContentSection>
        <Styled.SectionLabel>Date Range</Styled.SectionLabel>
        <DateRangeFilter onDateFilterChange={props.onDateFilterChange} />
      </Styled.ContentSection>
      {/* 
        Apply Button
      */}
      <Styled.ContentSection className="apply-button-container">
        <Styled.SectionLabel>&nbsp;</Styled.SectionLabel>
        <Styled.ApplyButtonContainer>
          <Button icon="refresh" onClick={props.onApply}>
            Apply
          </Button>
        </Styled.ApplyButtonContainer>
      </Styled.ContentSection>
    </Styled.Container>
  );
};

export default DataFilter;
