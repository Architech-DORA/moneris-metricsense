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

import { Fragment } from 'react';
// import { Tooltip2 } from '@blueprintjs/popover2';

import { useHome } from './use-home';
import * as S from './styled';
import { history } from '@/utils/history';

export const HomePage = () => {
  const { homeItems } = useHome();

  const onCardClickHandler = (path: string, isTarget?: boolean) => {
    if (isTarget) {
      window.open(path, '_blank');
    } else {
      history.push(path);
    }
  };

  return (
    <>
      {homeItems.map(({ key, title, description, items }) => (
        <Fragment key={key}>
          <S.StyledH1>{title}</S.StyledH1>
          <S.Description>{description}</S.Description>
          <S.CardWrapper>
            {items.map(({ key, iconUrl, title, description, isTarget, path, tooltipText, disabled = false }) => (
              <S.StyledCard
                key={key}
                interactive={true}
                disabled={disabled}
                onClick={(e) => {
                  if (!disabled) {
                    e.preventDefault();
                    e.stopPropagation();
                    onCardClickHandler(path, isTarget);
                  }
                }}
              >
                <S.TitleWrapper>
                  <S.LeftTitleWrapper>
                    {iconUrl && (
                      <S.StyledLeftSvgIcon url={iconUrl}>
                        <img src={iconUrl} alt={title} />
                      </S.StyledLeftSvgIcon>
                    )}
                    <h2>{title}</h2>
                  </S.LeftTitleWrapper>
                  {/* TODO: move this back once there are more information */}
                  {/* <Tooltip2 content={<>{tooltipText}</>}>
                    <S.StyledRightSvgIcon onClick={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                    }} url="/src/pages/home/assets/info.svg">
                      <img src="/src/pages/home/assets/info.svg" alt="Info" />
                    </S.StyledRightSvgIcon>
                  </Tooltip2> */}
                </S.TitleWrapper>
                <S.StyledP>{description}</S.StyledP>
              </S.StyledCard>
            ))}
          </S.CardWrapper>
        </Fragment>
      ))}
    </>
  );
};
