// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azurekv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/banzaicloud/bank-vaults/pkg/kv"
)

// azureKeyVault is an implementation of the kv.Service interface, that encrypts
// and decrypts and stores data using Azure Key Vault.
type azureKeyVault struct {
	client       *keyvault.BaseClient
	vaultBaseURL string
}

var _ kv.Service = &azureKeyVault{}

// New creates a new kv.Service backed by Azure Key Vault
func New(name string) (kv.Service, error) {
	keyClient := keyvault.New()
	authorizer, err := GetKeyvaultAuthorizer()
	if err != nil {
		return nil, err
	}
	keyClient.Authorizer = authorizer
	return &azureKeyVault{
		client:       &keyClient,
		vaultBaseURL: fmt.Sprintf("https://%s.vault.azure.net", name),
	}, nil
}

func (a *azureKeyVault) Get(key string) ([]byte, error) {

	bundle, err := a.client.GetSecret(context.Background(), a.vaultBaseURL, key, "")

	if err != nil {
		err := err.(autorest.DetailedError)
		if err.StatusCode == http.StatusNotFound {
			return nil, kv.NewNotFoundError("error getting secret for key '%s': %s", key, err.Error())
		}
		return nil, err
	}

	return []byte(*bundle.Value), nil
}

func (a *azureKeyVault) Set(key string, val []byte) error {

	value := string(val)
	parameters := keyvault.SecretSetParameters{
		Value: &value,
	}

	_, err := a.client.SetSecret(context.Background(), a.vaultBaseURL, key, parameters)

	return err
}

func (a *azureKeyVault) Test(key string) error {
	// TODO: Implement me properly
	return nil
}
