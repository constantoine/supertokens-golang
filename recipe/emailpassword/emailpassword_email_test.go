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

package emailpassword

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
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
			Init(&epmodels.TypeInput{
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{},
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						getUserByID := func(id string, userContext supertokens.UserContext) (*epmodels.User, error) {
							return &epmodels.User{
								ID:    "someId",
								Email: "someEmail",
							}, nil
						}
						originalImplementation.GetUserByID = &getUserByID
						return originalImplementation
					},
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

	SendEmail(emaildelivery.EmailType{
		PasswordReset: &emaildelivery.PasswordResetType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "someEmail",
			},
			PasswordResetLink: "someLink",
		},
	})

	assert.Equal(t, PasswordResetEmailSentForTest, true)
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
			Init(&epmodels.TypeInput{
				ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
						funcCalled = true
					},
				},
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						getUserByID := func(id string, userContext supertokens.UserContext) (*epmodels.User, error) {
							return &epmodels.User{
								ID:    "someId",
								Email: "someEmail",
							}, nil
						}
						originalImplementation.GetUserByID = &getUserByID
						return originalImplementation
					},
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

	SendEmail(emaildelivery.EmailType{
		PasswordReset: &emaildelivery.PasswordResetType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "someEmail",
			},
			PasswordResetLink: "someLink",
		},
	})

	assert.Equal(t, PasswordResetEmailSentForTest, false)
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
						(*originalImplementation.SendEmail) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
							overrideCalled = true
							return nil
						}
						return originalImplementation
					},
				},
				ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
						funcCalled = true
					},
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

	SendEmail(emaildelivery.EmailType{
		PasswordReset: &emaildelivery.PasswordResetType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "someEmail",
			},
			PasswordResetLink: "someLink",
		},
	})

	assert.Equal(t, PasswordResetEmailSentForTest, false)
	assert.Equal(t, funcCalled, false)
	assert.Equal(t, overrideCalled, true)
}

func TestSMTPServiceOverride(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
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

	err = SendEmail(emaildelivery.EmailType{
		PasswordReset: &emaildelivery.PasswordResetType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "",
			},
			PasswordResetLink: "someLink",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestEmailVerificationSMTPOverride(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
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

	err = SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID:    "someId",
				Email: "",
			},
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestEmailVerificationTokenThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				if input.EmailVerification != nil {
					getContentCalled = true
					return emaildelivery.SMTPGetContentResult{Body: "EmailVerify"}, nil
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				if input.Body == "EmailVerify" {
					sendRawEmailCalled = true
				}
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
			}),
			session.Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	signUpBodyResponse := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &signUpBodyResponse)
	assert.NoError(t, err)
	cookies := resp.Cookies()

	unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.Equal(t, PasswordResetEmailSentForTest, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestPasswordResetTokenThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
			}),
			session.Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	_, err = unittesting.SignupRequest("random@gmail.com", "validpass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)

	unittesting.PasswordResetTokenRequest("random@gmail.com", testServer.URL)
	assert.Equal(t, PasswordResetEmailSentForTest, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestPasswordResetTokenForNonExistantUserThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawEmailCalled := false
	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				getContentCalled = true
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
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
			Init(&epmodels.TypeInput{
				EmailDelivery: &emaildelivery.TypeInput{
					Service: &smtpService,
				},
			}),
			session.Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	unittesting.PasswordResetTokenRequest("random@gmail.com", testServer.URL)
	assert.Equal(t, PasswordResetEmailSentForTest, false)
	assert.Equal(t, getContentCalled, false)
	assert.Equal(t, sendRawEmailCalled, false)
}
