/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package emailverification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestBackwardCompatibilityServiceWithoutCustomFunction(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(evmodels.TypeInput{
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	(*singletonInstance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildeliverymodels.EmailType{
		EmailVerification: &emaildeliverymodels.EmailVerificationType{
			User: emaildeliverymodels.User{
				ID:    "someId",
				Email: "someEmail",
			},
		},
	}, nil)

	assert.Equal(t, EmailVerificationEmailSentForTest, true)
}

func TestBackwardCompatibilityServiceWithCustomFunction(t *testing.T) {
	funcCalled := false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(evmodels.TypeInput{
				CreateAndSendCustomEmail: func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
					funcCalled = true
				},
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	(*singletonInstance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildeliverymodels.EmailType{
		EmailVerification: &emaildeliverymodels.EmailVerificationType{
			User: emaildeliverymodels.User{
				ID:    "someId",
				Email: "someEmail",
			},
		},
	}, nil)

	assert.Equal(t, EmailVerificationEmailSentForTest, false)
	assert.Equal(t, funcCalled, true)
}

func TestBackwardCompatibilityServiceWithOverride(t *testing.T) {
	funcCalled := false
	overrideCalled := false
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(evmodels.TypeInput{
				EmailDelivery: &emaildeliverymodels.TypeInput{
					Override: func(originalImplementation emaildeliverymodels.EmailDeliveryInterface) emaildeliverymodels.EmailDeliveryInterface {
						(*originalImplementation.SendEmail) = func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) error {
							overrideCalled = true
							return nil
						}
						return originalImplementation
					},
				},
				CreateAndSendCustomEmail: func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
					funcCalled = true
				},
				GetEmailForUserID: func(userID string, userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	(*singletonInstance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildeliverymodels.EmailType{
		EmailVerification: &emaildeliverymodels.EmailVerificationType{
			User: emaildeliverymodels.User{
				ID:    "someId",
				Email: "someEmail",
			},
		},
	}, nil)

	assert.Equal(t, EmailVerificationEmailSentForTest, false)
	assert.Equal(t, funcCalled, false)
	assert.Equal(t, overrideCalled, true)
}
