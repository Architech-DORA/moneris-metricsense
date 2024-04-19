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
import { DateRangeInput3 } from '@blueprintjs/datetime2';
// import { type DateRange } from "@blueprintjs/datetime";
import { enCA } from 'date-fns/locale';

import * as Styled from './styled';

import './date-range-filter.css';

export type DateRangeFilterProps = {
  onDateFilterChange: (timestampFrom: string, timeStampTo: string) => void;
};

const DateRangeFilter = (props: DateRangeFilterProps) => {
  return (
    <Styled.DateRangeFilterContainer>
      <DateRangeInput3
        singleMonthOnly={true}
        dateFnsFormat={'yyyy-MM-dd'}
        onChange={(selectedRange) => {
          const dateFrom = new Date(`${selectedRange[0]}`);
          const dateTo = new Date(`${selectedRange[1]}`);
          dateFrom.setHours(0, 0, 0);
          dateTo.setHours(23, 59, 59);
          const timestampFrom = parseInt((dateFrom.getTime() / 1000).toFixed(2));
          const timestampTo = parseInt((dateTo.getTime() / 1000).toFixed(2));
          props.onDateFilterChange(`${timestampFrom}`, `${timestampTo}`);
        }}
        locale={enCA}
      />
    </Styled.DateRangeFilterContainer>
  );
};

export default DateRangeFilter;
