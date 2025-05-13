package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/anerdwithaknife/gh-app/internal/store"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var tokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Initiate OAuth flow for GitHub app",
	Long: `Starts a local server to handle the OAuth flow for a GitHub app.
This command will:
1. Start a local server to handle the OAuth callback
2. Open your browser to the GitHub OAuth authorization page
3. Exchange the received code for access and refresh tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, _ := cmd.Flags().GetString("slug")
		port, _ := cmd.Flags().GetInt("port")

		db, err := store.NewDefaultStore(false)
		if err != nil {
			log.Fatalf("Error loading store: %v", err)
		}
		defer db.Close()

		app, err := db.GetAppBySlug(slug)
		if err != nil {
			log.Fatalf("Error getting app details: %v", err)
		}

		clientId := app.ClientID
		clientSecret := app.ClientSecret

		srv := &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		}

		tokenChan := make(chan bool)

		http.HandleFunc("/api/auth/callback/github", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Error: 'code' parameter not provided", http.StatusBadRequest)
				return
			}

			//log.Printf("Received code: %s", code)

			params := url.Values{}
			params.Add("client_id", clientId)
			params.Add("client_secret", clientSecret)
			params.Add("code", code)

			req, err := http.NewRequest(
				"POST",
				"https://github.com/login/oauth/access_token",
				nil,
			)
			if err != nil {
				log.Printf("Error creating token request: %v", err)
				http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
				return
			}

			req.URL.RawQuery = params.Encode()
			req.Header.Set("Accept", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error exchanging code for token: %v", err)
				http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
				log.Printf("Error decoding token response: %v", err)
				http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Access Token: %s\nRefresh Token: %s\n",
				tokenData.AccessToken, tokenData.RefreshToken)

			drawTokenTable(tokenData.AccessToken, tokenData.RefreshToken)
			cmd.Println()

			tokenChan <- true
		})

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Printf("Error starting server: %v", err)
			}
		}()

		authURL := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=http://localhost:%d/api/auth/callback/github",
			clientId,
			port,
		)
		boxPrint(authURL)
		cmd.Println("")

		if err := openBrowser(authURL); err != nil {
			log.Printf("Error opening browser: %v", err)
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		select {
		case <-tokenChan:
			//log.Println("Successfully received tokens")
		case <-interrupt:
			log.Println("Received interrupt signal")
		case <-time.After(5 * time.Minute):
			log.Println("Timeout waiting for authorization")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(oauthCmd)

	oauthCmd.Flags().StringP("slug", "s", "", "Slug of the GitHub app")
	oauthCmd.Flags().IntP("port", "p", 3000, "Localhost port to listen on for OAuth callback")

	oauthCmd.MarkFlagRequired("slug")
}

func boxPrint(text string) {
	length := len(text) + 2
	border := strings.Repeat("─", length)
	box := fmt.Sprintf("┌%s┐\n│ %s │\n└%s┘", border, text, border)
	fmt.Println(box)
}

func drawTokenTable(accessToken string, refreshToken string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Type", "Token")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPrintHeaders(false)

	tbl.AddRow("Access Token", accessToken)
	tbl.AddRow("Refresh Token", refreshToken)

	tbl.Print()
}

// openBrowser opens the specified URL in the default browser of the user's system
func openBrowser(url string) error {
	var cmd string
	var args []string

	browser := os.Getenv("BROWSER")
	if browser != "" {
		cmd = browser
	} else {
		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default:
			cmd = "xdg-open"
		}
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
