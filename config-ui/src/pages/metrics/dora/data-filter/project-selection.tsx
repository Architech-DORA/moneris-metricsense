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
import { useState } from 'react';
import { Intent, MenuItem } from '@blueprintjs/core';
import { ItemPredicate, ItemRenderer, MultiSelect } from '@blueprintjs/select';
import * as Styled from './styled';

import './style.css';

export interface Project {
  name: string;
}

export type BaseProps = {
  items: Project[];
  onChangeProjects: (projects: Project[]) => void;
};

const filterProject: ItemPredicate<Project> = (query, project, _index, exactMatch) => {
  const normalizedTitle = project.name.toLowerCase();
  const normalizedQuery = query.toLowerCase();

  if (exactMatch) {
    return normalizedTitle === normalizedQuery;
  } else {
    return `${normalizedTitle}`.indexOf(normalizedQuery) >= 0;
  }
};

const renderProject: ItemRenderer<Project> = (project, { handleClick, handleFocus, modifiers, query }) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      key={project.name}
      label={project.name.toString()}
      onClick={handleClick}
      onFocus={handleFocus}
      roleStructure="listoption"
      text={`${project.name}`}
      shouldDismissPopover={false}
    />
  );
};

type ProjectSelectProps = {} & BaseProps;

const ProjectSelect = (props: ProjectSelectProps) => {
  const [selectedProject, setSelectedProject] = useState<Project[]>([]);

  const handleDeselect = (index: number) => {
    const _selectedProject = selectedProject.filter((_, i) => (i === index ? false : true));

    setSelectedProject(_selectedProject);
    props.onChangeProjects(_selectedProject);
  };

  const handleSelect = (item: Project) => {
    const indexOfSelectedItem = selectedProject.findIndex((project: Project) => {
      return project.name === item.name;
    });
    const isAlreadySelected = indexOfSelectedItem > -1;
    if (isAlreadySelected) {
      handleDeselect(indexOfSelectedItem);
    } else {
      const _selectedProject = selectedProject;
      _selectedProject.push(item);
      setSelectedProject(_selectedProject);
      props.onChangeProjects(_selectedProject);
    }
  };

  return (
    <MultiSelect<Project>
      items={props.items}
      menuProps={{ 'aria-label': 'films' }}
      itemPredicate={filterProject}
      itemRenderer={renderProject}
      noResults={<MenuItem disabled={true} text="No results." roleStructure="listoption" />}
      onItemSelect={(item) => {
        handleSelect(item);
      }}
      fill={true}
      popoverProps={{ matchTargetWidth: true, minimal: true }}
      selectedItems={selectedProject}
      tagRenderer={(project: Project) => project.name}
      tagInputProps={{
        onRemove: (_, index) => {
          handleDeselect(index);
        },
        tagProps: {
          intent: Intent.NONE,
          minimal: false,
        },
      }}
    />
  );
};

type ProjectSelectionProps = {} & BaseProps;

const ProjectSelection = (props: ProjectSelectionProps) => {
  return (
    <Styled.SelectionContainer>
      <ProjectSelect {...props} />
    </Styled.SelectionContainer>
  );
};

export default ProjectSelection;
