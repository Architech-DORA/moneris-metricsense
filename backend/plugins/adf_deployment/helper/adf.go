/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datafactory/mgmt/datafactory"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/apache/incubator-devlake/core/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type AdfApiClient struct {
	Client *datafactory.FactoriesClient
}

func (k *AdfApiClient) TestConnection() error {
	fmt.Println("TestConnection from ADF API Client")
	return nil
}

type AdfApiClientSet struct {
	ClientSet *kubernetes.Clientset
}

func NewAdfApiClient(credentials map[string]interface{}) (*AdfApiClient, errors.Error) {
	fmt.Println("credentials: ", credentials)
	providerType, ok := credentials["providerType"].(string)

	if providerType == "" || !ok {
		err := errors.BadInput.New("providerType is not defined")
		return nil, err
	}

	var adfApiClient *datafactory.FactoriesClient
	var err errors.Error

	if providerType == "azure" {
		fmt.Println("Creating Azure ADF Client")
		adfApiClient, err = createAzureClientConfig(credentials)

		if err != nil {
			return nil, err
		}
	}

	return &AdfApiClient{
		Client: adfApiClient,
	}, nil
}

func NewKubeApiClient(credentials map[string]interface{}) (*AdfApiClientSet, errors.Error) {
	println("credentials: ", credentials)
	providerType, ok := credentials["providerType"].(string)

	if providerType == "" || !ok {
		err := errors.BadInput.New("providerType is not defined")
		return nil, err
	}

	var adfApiClient *kubernetes.Clientset
	var err errors.Error

	if providerType == "azure" {
		fmt.Println("Creating Azure ADF Client")
		adfApiClient, err = createKubeClientConfig(credentials)

		if err != nil {
			return nil, err
		}
	} else if providerType == "local" {
		fmt.Println("Creating Local Kube Client")
		adfApiClient, err = createClientConfig()

		if err != nil {
			return nil, err
		}
	}

	return &AdfApiClientSet{
		ClientSet: adfApiClient,
	}, nil
}

func createAzureClientConfig(creds map[string]interface{}) (*datafactory.FactoriesClient, errors.Error) {
	clientID := creds["clientID"].(string)
	clientSecret := creds["clientSecret"].(string)
	tenantID := creds["tenantID"].(string)
	subscriptionID := creds["subscriptionID"].(string)
	factoryName := creds["factoryName"].(string)
	resourceGroupName := creds["resourceGroupName"].(string)

	authorizer, err := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID).Authorizer()
	if err != nil {
		panic(err.Error())
	}

	dataFactoryClient := datafactory.NewFactoriesClient(subscriptionID)
	dataFactoryClient.Authorizer = authorizer


	df, err := dataFactoryClient.Get(context.Background(), resourceGroupName, factoryName, "")
	if err != nil {
		return nil, errors.Default.New("Unable to establish connection to Azure Data Factory Resource Group")
	}


	fmt.Printf("Data Factory Name: %v\n", *df.Name)
	fmt.Printf("Data Factory ID: %v\n", *df.ID)

	

	return &dataFactoryClient, nil
}


func createKubeClientConfig(creds map[string]interface{}) (*kubernetes.Clientset, errors.Error) {
	clientID := creds["clientID"].(string)
	clientSecret := creds["clientSecret"].(string)
	tenantID := creds["tenantID"].(string)
	subscriptionID := creds["subscriptionID"].(string)
	clusterName := creds["clusterName"].(string)
	resourceGroupName := creds["resourceGroupName"].(string)

	authorizer, err := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID).Authorizer()
	if err != nil {
		panic(err.Error())
	}

	client := containerservice.NewManagedClustersClient(subscriptionID)

	client.Authorizer = authorizer

	credentials, err := client.ListClusterAdminCredentials(context.Background(), resourceGroupName, clusterName, "")
	if err != nil {
		return nil, errors.Default.New("Unable to establish connection to Azure Kubernetes Cluster")
	}

	kubeConfigs := *credentials.Kubeconfigs

	if len(kubeConfigs) == 0 {
		return nil, errors.Default.New("Empty kube config")
	}

	kubeConfigAZ := *kubeConfigs[0].Value

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfigAZ)

	if err != nil {
		return nil, errors.BadInput.New("Unable to establish connection to create client config")
	}

	config, err := clientConfig.ClientConfig()

	if err != nil {
		return nil, errors.BadInput.New("Unable to establish connection to create client config")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientSet, nil
}


func createClientConfig() (*kubernetes.Clientset, errors.Error) {
	kubeConfig := "~/.kube/config"

	// Build the client config from the kubeConfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, errors.Default.New("Unable to establish connection to local Kubernetes Cluster")
	}

	// Create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Default.New("Unable to establish connection to create client config")
	}

	return clientSet, nil
}