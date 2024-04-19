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
const AWS_COGNITO_DOMAIN = import.meta.env.VITE_AWS_COGNITO_DOMAIN;
const AWS_CLIENT_ID = import.meta.env.VITE_AWS_COGNITO_CLIENT_ID;
const AWS_LOGIN_REDIRECT_URI = import.meta.env.VITE_AWS_COGNITO_LOGIN_REDIRECT;

export const getLogoutUrl = () => {
  return `${AWS_COGNITO_DOMAIN}/logout?response_type=code&client_id=${AWS_CLIENT_ID}&redirect_uri=${encodeURIComponent(
    AWS_LOGIN_REDIRECT_URI,
  )}&scope=openid`;
};
