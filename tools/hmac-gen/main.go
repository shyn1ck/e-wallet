package main

import (
	"bufio"
	"e-wallet/pkg/crypto"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Color constants for terminal output
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorRed    = "\033[31m"
)

// API endpoint constants
const (
	endpointWalletCheck        = "/wallet/check"
	endpointWalletBalance      = "/wallet/balance"
	endpointWalletDeposit      = "/wallet/deposit"
	endpointWalletMonthlyStats = "/wallet/monthly-stats"
)

// Default credentials
const (
	defaultUserID    = "alif_partner"
	defaultSecretKey = "alif_secret_2025"
)

// API client credentials
const (
	userIDAlifPartner         = "alif_partner"
	secretKeyAlifPartner      = "alif_secret_2025"
	userIDMegafonAPI          = "megafon_api"
	secretKeyMegafonAPI       = "megafon_key_secure"
	userIDTcellIntegration    = "tcell_integration"
	secretKeyTcellIntegration = "tcell_hmac_key"
)

// API base URL
const apiBaseURL = "http://localhost:8080/api/v1"

func main() {
	interactive := flag.Bool("i", false, "Interactive mode")
	body := flag.String("body", "", "JSON request body")
	secret := flag.String("secret", defaultSecretKey, "Secret key for HMAC")
	userID := flag.String("user", defaultUserID, "User ID for X-UserId header")
	endpoint := flag.String("endpoint", endpointWalletBalance, "API endpoint")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *interactive {
		runInteractive()
		return
	}

	if *body == "" {
		runInteractive()
		return
	}

	generateAndPrint(*body, *secret, *userID, *endpoint)
}

func runInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s==========================================%s\n", colorBlue, colorReset)
	fmt.Printf("%sE-Wallet API - HMAC Generator (Interactive)%s\n", colorBlue, colorReset)
	fmt.Printf("%s==========================================%s\n\n", colorBlue, colorReset)

	// Select operation
	fmt.Printf("%sSelect operation:%s\n", colorYellow, colorReset)
	fmt.Println("1. Check wallet")
	fmt.Println("2. Get balance")
	fmt.Println("3. Deposit")
	fmt.Println("4. Monthly stats")
	fmt.Print("\nEnter choice (1-4): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var endpoint string
	var needAmount bool

	switch choice {
	case "1":
		endpoint = endpointWalletCheck
	case "2":
		endpoint = endpointWalletBalance
	case "3":
		endpoint = endpointWalletDeposit
		needAmount = true
	case "4":
		endpoint = endpointWalletMonthlyStats
	default:
		fmt.Printf("%sInvalid choice%s\n", colorRed, colorReset)
		return
	}

	// Enter account ID
	fmt.Print("\nEnter account_id (e.g., 992900123456): ")
	accountID, _ := reader.ReadString('\n')
	accountID = strings.TrimSpace(accountID)

	if accountID == "" {
		fmt.Printf("%sAccount ID cannot be empty%s\n", colorRed, colorReset)
		return
	}

	// Build JSON body
	var body string
	if needAmount {
		fmt.Print("Enter amount (e.g., 10000): ")
		amount, _ := reader.ReadString('\n')
		amount = strings.TrimSpace(amount)
		body = fmt.Sprintf(`{"account_id":"%s","amount":%s}`, accountID, amount)
	} else {
		body = fmt.Sprintf(`{"account_id":"%s"}`, accountID)
	}

	// Select credentials
	fmt.Printf("\n%sSelect credentials:%s\n", colorYellow, colorReset)
	fmt.Println("1. alif_partner / alif_secret_2025")
	fmt.Println("2. megafon_api / megafon_key_secure")
	fmt.Println("3. tcell_integration / tcell_hmac_key")
	fmt.Print("\nEnter choice (1-3, default 1): ")

	credChoice, _ := reader.ReadString('\n')
	credChoice = strings.TrimSpace(credChoice)

	var userID, secret string
	switch credChoice {
	case "2":
		userID = userIDMegafonAPI
		secret = secretKeyMegafonAPI
	case "3":
		userID = userIDTcellIntegration
		secret = secretKeyTcellIntegration
	default:
		userID = userIDAlifPartner
		secret = secretKeyAlifPartner
	}

	fmt.Println()
	generateAndPrint(body, secret, userID, endpoint)
}

func generateAndPrint(body, secret, userID, endpoint string) {
	digest := crypto.ComputeHMAC(crypto.AlgorithmSHA1, secret, body)

	fmt.Printf("%s==========================================%s\n", colorBlue, colorReset)
	fmt.Printf("%sGenerated HMAC-SHA1 Signature%s\n", colorBlue, colorReset)
	fmt.Printf("%s==========================================%s\n\n", colorBlue, colorReset)

	fmt.Printf("%sRequest Body:%s\n", colorGreen, colorReset)
	fmt.Printf("%s\n\n", body)

	fmt.Printf("%sCredentials:%s\n", colorGreen, colorReset)
	fmt.Printf("User ID: %s\n", userID)
	fmt.Printf("Secret:  %s\n\n", secret)

	fmt.Printf("%sHMAC-SHA1 Digest:%s\n", colorGreen, colorReset)
	fmt.Printf("%s%s%s\n\n", colorYellow, digest, colorReset)

	fmt.Printf("%s==========================================%s\n", colorBlue, colorReset)
	fmt.Printf("%scURL Command:%s\n\n", colorGreen, colorReset)

	fmt.Printf("curl -X POST %s%s \\\n", apiBaseURL, endpoint)
	fmt.Printf("  -H 'Content-Type: application/json' \\\n")
	fmt.Printf("  -H 'X-UserId: %s' \\\n", userID)
	fmt.Printf("  -H 'X-Digest: %s' \\\n", digest)
	fmt.Printf("  -d '%s'\n\n", body)

	fmt.Printf("%s==========================================%s\n", colorBlue, colorReset)
}

func printHelp() {
	fmt.Printf("%s==========================================%s\n", colorBlue, colorReset)
	fmt.Printf("%sE-Wallet API - HMAC-SHA1 Generator%s\n", colorBlue, colorReset)
	fmt.Printf("%s==========================================%s\n\n", colorBlue, colorReset)

	fmt.Printf("%sUsage:%s\n", colorYellow, colorReset)
	fmt.Println("  go run tools/hmac-gen/main.go [options]")
	fmt.Println()

	fmt.Printf("%sOptions:%s\n", colorYellow, colorReset)
	fmt.Println("  -i")
	fmt.Println("        Interactive mode (recommended)")
	fmt.Println("  -body string")
	fmt.Println("        JSON request body")
	fmt.Println("  -secret string")
	fmt.Println("        Secret key for HMAC (default: alif_secret_2025)")
	fmt.Println("  -user string")
	fmt.Println("        User ID for X-UserId header (default: alif_partner)")
	fmt.Println("  -endpoint string")
	fmt.Println("        API endpoint (default: /wallet/balance)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()

	fmt.Printf("%sExamples:%s\n", colorYellow, colorReset)
	fmt.Println("  # Interactive mode (default)")
	fmt.Println("  go run tools/hmac-gen/main.go")
	fmt.Println()
	fmt.Println("  # Direct mode")
	fmt.Println(`  go run tools/hmac-gen/main.go -body '{"account_id":"992900123456"}'`)
	fmt.Println()
	fmt.Println("  # With custom credentials")
	fmt.Println(`  go run tools/hmac-gen/main.go -body '{"account_id":"992900111222"}' -user megafon_api -secret megafon_key_secure`)
	fmt.Println()

	fmt.Printf("%sAvailable Test Credentials:%s\n", colorYellow, colorReset)
	fmt.Println("  User ID: alif_partner,       Secret: alif_secret_2025")
	fmt.Println("  User ID: megafon_api,        Secret: megafon_key_secure")
	fmt.Println("  User ID: tcell_integration,  Secret: tcell_hmac_key")
	fmt.Println()

	fmt.Printf("%sAPI Endpoints:%s\n", colorYellow, colorReset)
	fmt.Println("  POST /api/v1/wallet/check         - Check if wallet exists")
	fmt.Println("  POST /api/v1/wallet/balance       - Get wallet balance")
	fmt.Println("  POST /api/v1/wallet/deposit       - Deposit to wallet")
	fmt.Println("  POST /api/v1/wallet/monthly-stats - Get monthly statistics")
	fmt.Println()
}
