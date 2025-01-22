package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chrlesur/aiyou.golib"
	"github.com/spf13/cobra"
)

var (
	schemaFile string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configuration file",
	Long:  `Generate a YAML configuration file based on machine type and context`,
	RunE:  runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&machineType, "type", "", "Type of machine (router, switch, firewall, server)")
	generateCmd.Flags().StringVar(&userContext, "context", "", "Context description in natural language")
	generateCmd.Flags().StringVar(&schemaFile, "schema", "", "Path to JSON schema file")
	generateCmd.MarkFlagRequired("type")
	generateCmd.MarkFlagRequired("context")
	generateCmd.MarkFlagRequired("schema")
}

func createClient() (*aiyou.Client, error) {
	if userEmail == "" || userPass == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	logger := aiyou.NewDefaultLogger(os.Stderr)
	if debug {
		logger.SetLevel(aiyou.DEBUG)
	} else if quietMode {
		logger.SetLevel(aiyou.ERROR)
	} else {
		logger.SetLevel(aiyou.INFO)
	}

	return aiyou.NewClient(
		aiyou.WithEmailPassword(userEmail, userPass),
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL("https://ai.dragonflygroup.fr"),
	)
}

func loadSchema(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %v", err)
	}
	return string(content), nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Load schema
	schema, err := loadSchema(schemaFile)
	if err != nil {
		return err
	}

	// Prepare the prompt
	prompt := fmt.Sprintf(`Using the following JSON schema as validation reference:

%s

Generate a YAML configuration for a %s with the following context: %s.

Important instructions:
1. Response must be ONLY the YAML configuration
2. No explanations or comments
3. Strictly follow the schema structure
4. No markdown formatting or code block markers`, schema, machineType, userContext)

	// Prepare the request
	req := aiyou.ChatCompletionRequest{
		Messages: []aiyou.Message{
			{
				Role: "user",
				Content: []aiyou.ContentPart{
					{Type: "text", Text: prompt},
				},
			},
		},
		AssistantID: assistantID,
		Stream:      true,
		Temperature: 0.7,
		TopP:        0.95,
	}

	// Get completion using streaming
	ctx := context.Background()
	stream, err := client.ChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start completion stream: %v", err)
	}

	var fullResponse strings.Builder
	if !quietMode {
		fmt.Println("Generating configuration...")
	}

	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading stream: %v", err)
		}

		if chunk != nil && len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			if choice.Delta != nil && choice.Delta.Content != "" {
				fullResponse.WriteString(choice.Delta.Content)
			}
		}
	}

	// Clean up the response
	yamlContent := fullResponse.String()
	yamlContent = cleanYAMLResponse(yamlContent)

	if yamlContent == "" {
		return fmt.Errorf("received empty YAML content from AI.YOU")
	}

	// Generate output filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	outputFile := fmt.Sprintf("%s_config_%s.yaml", machineType, timestamp)
	outputPath := filepath.Join(".", outputFile)

	// Write to file
	err = os.WriteFile(outputPath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write configuration file: %v", err)
	}

	if !quietMode {
		fmt.Printf("Configuration generated: %s\n", outputPath)
	}
	return nil
}

func cleanYAMLResponse(response string) string {
	// Remove markdown code block markers if present
	response = strings.TrimPrefix(response, "```yaml")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	
	// Remove any explanatory text before or after the YAML
	lines := strings.Split(response, "\n")
	var yamlLines []string
	inYAML := false
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "#") {
			continue // Skip comments
		}
		
		if !inYAML {
			if strings.Contains(line, ":") {
				inYAML = true
				yamlLines = append(yamlLines, line)
			}
		} else {
			if trimmedLine == "" {
				continue
			}
			yamlLines = append(yamlLines, line)
		}
	}
	
	return strings.TrimSpace(strings.Join(yamlLines, "\n"))
}