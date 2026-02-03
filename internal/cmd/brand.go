package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var brandCmd = &cobra.Command{
	Use:   "brand",
	Short: "Display your brand DNA",
	Long: `Display the active brand DNA for your organization.

This includes your color palette, fonts, and brand voice guidelines
that were extracted from your website.`,
	RunE: runBrand,
}

func runBrand(cmd *cobra.Command, args []string) error {
	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)
	brand, err := client.GetBrand()
	if err != nil {
		return fmt.Errorf("failed to fetch brand: %w", err)
	}

	if brand == nil {
		fmt.Println("No active brand DNA found.")
		fmt.Println("Visit the web app to set up your brand.")
		return nil
	}

	if jsonOutput {
		return outputJSON(brand)
	}

	return outputBrandText(brand)
}

func outputBrandText(brand *api.BrandDNA) error {
	fmt.Printf("Brand DNA (v%d)\n", brand.Version)
	fmt.Println("================")
	fmt.Println()

	if brand.CompanyName != "" {
		fmt.Printf("Company: %s\n", brand.CompanyName)
	}
	if brand.SourceURL != "" {
		fmt.Printf("Source: %s\n", brand.SourceURL)
	}
	fmt.Println()

	// Palette
	fmt.Println("Color Palette:")
	if palette, ok := brand.Palette["colors"].([]interface{}); ok {
		for _, c := range palette {
			if color, ok := c.(map[string]interface{}); ok {
				name := color["name"]
				hex := color["hex"]
				role := color["role"]
				fmt.Printf("  - %s (%s): %s\n", name, role, hex)
			}
		}
	} else {
		fmt.Println("  (not available)")
	}
	fmt.Println()

	// Fonts
	fmt.Println("Typography:")
	if fonts, ok := brand.Fonts["fonts"].([]interface{}); ok {
		for _, f := range fonts {
			if font, ok := f.(map[string]interface{}); ok {
				name := font["name"]
				usage := font["usage"]
				fmt.Printf("  - %s: %s\n", usage, name)
			}
		}
	} else {
		fmt.Println("  (not available)")
	}
	fmt.Println()

	// Voice
	fmt.Println("Brand Voice:")
	if voice, ok := brand.Voice["voice"].(map[string]interface{}); ok {
		if tone, ok := voice["tone"].(string); ok {
			fmt.Printf("  Tone: %s\n", tone)
		}
		if personality, ok := voice["personality"].([]interface{}); ok {
			fmt.Print("  Personality: ")
			for i, p := range personality {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(p)
			}
			fmt.Println()
		}
		if audience, ok := voice["target_audience"].(string); ok {
			fmt.Printf("  Target Audience: %s\n", audience)
		}
	} else {
		fmt.Println("  (not available)")
	}

	return nil
}

func outputJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
