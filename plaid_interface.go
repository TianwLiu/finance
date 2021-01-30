package main

import (
	"github.com/plaid/plaid-go/plaid"
	"net/http"
)

const ENV_DEVELOPMENT="development"
const ENV_SANDBOX="sandbox"
var plaidClient *plaid.Client


func plaidReady(environment string,plaidClientId string,plaidClientSecret string){
	var clientOptions plaid.ClientOptions
	//if debug, use sandbox environment
	//else, use development environment
	if environment==ENV_DEVELOPMENT {
		clientOptions = plaid.ClientOptions{
			plaidClientId,
			plaidClientSecret,
			plaid.Development, // Available environments are Sandbox, Development, and Production
			&http.Client{}, // This parameter is optional
		}
	}else if environment==ENV_SANDBOX {
		clientOptions = plaid.ClientOptions{
			plaidClientId,
			plaidClientSecret,
			plaid.Sandbox, // Available environments are Sandbox, Development, and Production
			&http.Client{}, // This parameter is optional
		}

	}
	plaidClient, _ = plaid.NewClient(clientOptions)

}

func getAccessToken(publicToken string) (string,error){
	accessTokenResponse,err := plaidClient.ExchangePublicToken(publicToken)
	return accessTokenResponse.AccessToken,err
}
