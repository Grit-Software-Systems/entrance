package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Grit-Software-Systems/entrance"
)

const (
	defaultScopes = "openid profile offline_access"

	exitCodeInvalidArguments = 1
	exitCodeSuccess          = 0

	flagNameClientIdentifier = "client-identifier"
	flagNameEmail            = "email"
	flagNameScopes           = "scopes"
	flagNameShowFullTokens   = "show-full-tokens"
	flagNameTenantIdentifier = "tenant-identifier"
	flagNameTenantSubdomain  = "tenant-subdomain"

	menuChoiceExit                     = 5
	menuChoiceOneTimePasscode          = 1
	menuChoicePassword                 = 2
	menuChoicePasswordWithMultifactor  = 3
	menuChoicePasswordWithRegistration = 4

	menuPrompt = "Select a choice: "
)

func main() {
	var tenantSubdomain string
	flag.StringVar(&tenantSubdomain, flagNameTenantSubdomain, "", "Tenant subdomain, e.g. contoso")

	var tenantIdentifier string
	flag.StringVar(&tenantIdentifier, flagNameTenantIdentifier, "", "Tenant GUID or {tenant}.onmicrosoft.com")

	var clientIdentifier string
	flag.StringVar(&clientIdentifier, flagNameClientIdentifier, "", "Application (client) ID")

	var scopes string
	flag.StringVar(&scopes, flagNameScopes, defaultScopes, "Space-separated OAuth scopes")

	var emailAddress string
	flag.StringVar(&emailAddress, flagNameEmail, "", "Email address of the user to authenticate")

	var showFullTokens bool
	flag.BoolVar(&showFullTokens, flagNameShowFullTokens, false, "Print full token values instead of truncating")

	flag.Parse()

	var requiredFlagsMissing bool = validateRequiredFlags(tenantSubdomain, tenantIdentifier, clientIdentifier, emailAddress)
	if requiredFlagsMissing {
		flag.Usage()
		os.Exit(exitCodeInvalidArguments)
	}

	var configuration entrance.Configuration = entrance.Configuration{
		TenantSubdomain:  tenantSubdomain,
		TenantIdentifier: tenantIdentifier,
		ClientIdentifier: clientIdentifier,
		Scopes:           scopes,
	}

	var authenticationClient entrance.Client = entrance.NewClient(configuration)

	runMenuLoop(authenticationClient, emailAddress, showFullTokens)

	os.Exit(exitCodeSuccess)
}

func printMenu() {
	fmt.Println("Select an authentication flow:")
	fmt.Println("  1. Email one-time passcode sign-in")
	fmt.Println("  2. Password sign-in")
	fmt.Println("  3. Password sign-in with multifactor authentication")
	fmt.Println("  4. Password sign-in with registration")
	fmt.Println("  5. Exit")
}

func runMenuLoop(authenticationClient entrance.Client, emailAddress string, showFullTokens bool) {
	for {
		printMenu()
		var choice int = readChoice(menuPrompt, menuChoiceOneTimePasscode, menuChoiceExit)
		if choice == menuChoiceExit {
			return
		}
		if choice == menuChoiceOneTimePasscode {
			signInWithOneTimePasscode(authenticationClient, emailAddress, showFullTokens)
		}
		if choice == menuChoicePassword {
			signInWithPassword(authenticationClient, emailAddress, showFullTokens)
		}
		if choice == menuChoicePasswordWithMultifactor {
			signInWithPasswordAndMultifactor(authenticationClient, emailAddress, showFullTokens)
		}
		if choice == menuChoicePasswordWithRegistration {
			signInWithPasswordAndRegistration(authenticationClient, emailAddress, showFullTokens)
		}
		fmt.Println()
	}
}

func validateRequiredFlags(tenantSubdomain string, tenantIdentifier string, clientIdentifier string, emailAddress string) bool {
	if tenantSubdomain == "" {
		result := true
		return result
	}
	if tenantIdentifier == "" {
		result := true
		return result
	}
	if clientIdentifier == "" {
		result := true
		return result
	}
	if emailAddress == "" {
		result := true
		return result
	}
	result := false
	return result
}
