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

import { FormGroup, InputGroup, Radio, RadioGroup } from '@blueprintjs/core';

import { useEffect } from 'react';
import * as S from './styled';

interface Props {
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
  // this will be the original setValues function, to provide full control to the state
  setValuesDefault: React.Dispatch<React.SetStateAction<Record<string, any>>>;
}

interface AzureAuth {
  providerType: 'azure';
  subscriptionID: string;
  clientID: string;
  clientSecret: string;
  tenantID: string;
  resourceGroupName: string;
  factoryName: string;
}

export const Auth = ({ initialValues, values, errors, setValues, setErrors, setValuesDefault }: Props) => {
  useEffect(() => {
    if (!values.credentials?.providerType) {
      setErrors({
        error: 'Required',
      });
    } else {
      // all fields in the 3 auth types are required, if any fields is empty string or undefined, set error
      if (values.credentials?.providerType === 'azure') {
        if (
          values?.credentials?.subscriptionID === '' ||
          values?.credentials?.clientID === '' ||
          values?.credentials?.clientSecret === '' ||
          values?.credentials?.tenantID === '' ||
          values?.credentials?.resourceGroupName === '' ||
          values?.credentials?.factoryName === ''
        ) {
          setErrors({
            error: 'Required',
          });
          // unset errors
        } else {
          setErrors({ error: '' });
        }
      }
    }
  }, [values]);

  const defaultValues = {
    authMethod: 'AccessToken',
    endpoint: `${window.location.origin}/adf_deployment/connection/fallback-endpoint`,
    id: 23,
    proxy: '',
    rateLimitPerHour: 0,
  };

  // setInitialValues
  useEffect(() => {
    setValuesDefault((prev) => ({
      ...defaultValues,
      ...prev,
      credentials: {
        ...initialValues.credentials,
        providerType: 'azure',
        subscriptionID: '',
        clientID: '',
        clientSecret: '',
        tenantID: '',
        resourceGroupName: '',
        factoryName: '',
      }
    }));
  }, []);

  // setInitialValues
  useEffect(() => {
    if (initialValues.credentials?.providerType === 'azure') {
      setValuesDefault((prev) => ({
        ...prev,
        credentials: initialValues.credentials,
      }));
    }
  }, [initialValues.credentials]);

  return (
    <>
      {values.credentials?.providerType === 'azure' && (
        <>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Subscription ID</S.Label>}>
            <InputGroup
              placeholder="Your Subscription ID"
              value={(values.credentials as AzureAuth).subscriptionID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    subscriptionID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client ID</S.Label>}>
            <InputGroup
              placeholder="Your Client ID"
              value={(values.credentials as AzureAuth).clientID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clientID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client Secret</S.Label>}>
            <InputGroup
              placeholder="Your Client Secret"
              type={'password'}
              value={(values.credentials as AzureAuth).clientSecret}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clientSecret: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Tenant ID</S.Label>}>
            <InputGroup
              placeholder="Your Tenant ID"
              value={(values.credentials as AzureAuth).tenantID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    tenantID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Resource Group Name</S.Label>}>
            <InputGroup
              placeholder="Your Resource Group Name"
              value={(values.credentials as AzureAuth).resourceGroupName}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    resourceGroupName: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Factory Name</S.Label>}>
            <InputGroup
              placeholder="Your Factory Name"
              value={(values.credentials as AzureAuth).factoryName}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    factoryName: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
        </>
      )}
    </>
  );
};
