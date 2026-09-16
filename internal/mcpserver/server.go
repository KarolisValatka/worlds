package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/KarolisValatka/worlds/internal/vault"
)

// New builds an MCP server bound to vaultDir.
func New(vaultDir string) (*mcp.Server, error) {
	abs, err := vault.Abs(vaultDir)
	if err != nil {
		return nil, err
	}
	if err := vault.Validate(abs); err != nil {
		return nil, err
	}
	s := mcp.NewServer(&mcp.Implementation{Name: "worlds", Version: "0.2.0"}, nil)
	registerTools(s, abs)
	return s, nil
}

// RunStdio serves MCP over stdin/stdout until the client disconnects.
func RunStdio(ctx context.Context, vaultDir string) error {
	s, err := New(vaultDir)
	if err != nil {
		return err
	}
	return s.Run(ctx, &mcp.StdioTransport{})
}

func registerTools(s *mcp.Server, vaultDir string) {
	type emptyArgs struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_worlds",
		Description: "List world directories in the vault",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ emptyArgs) (*mcp.CallToolResult, any, error) {
		worlds, err := vault.ListWorlds(vaultDir)
		if err != nil {
			return textErr(err), nil, nil
		}
		return textOK(strings.Join(worlds, "\n")), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "read_brain",
		Description: "Read the vault root BRAIN.md hot index",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ emptyArgs) (*mcp.CallToolResult, any, error) {
		body, err := vault.ReadBrain(vaultDir)
		if err != nil {
			return textErr(err), nil, nil
		}
		return textOK(body), nil, nil
	})

	type readWorldArgs struct {
		World    string `json:"world" jsonschema:"world directory name (kebab-case)"`
		ListWiki bool   `json:"list_wiki,omitempty" jsonschema:"if true, also list wiki page titles"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "read_world",
		Description: "Read a world BRAIN.md; optionally list wiki page titles",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args readWorldArgs) (*mcp.CallToolResult, any, error) {
		if args.World == "" {
			return textErr(fmt.Errorf("world is required")), nil, nil
		}
		body, err := vault.ReadWorldBrain(vaultDir, args.World)
		if err != nil {
			return textErr(err), nil, nil
		}
		if args.ListWiki {
			pages, err := vault.ListWikiPages(vaultDir, args.World)
			if err != nil {
				return textErr(err), nil, nil
			}
			body += "\n\n## wiki pages\n"
			if len(pages) == 0 {
				body += "(none)\n"
			} else {
				for _, p := range pages {
					body += "- " + p + "\n"
				}
			}
		}
		return textOK(body), nil, nil
	})

	type readWikiArgs struct {
		World string `json:"world" jsonschema:"world directory name"`
		Page  string `json:"page" jsonschema:"wiki page name without .md"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "read_wiki",
		Description: "Read a wiki page including YAML frontmatter",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args readWikiArgs) (*mcp.CallToolResult, any, error) {
		if args.World == "" || args.Page == "" {
			return textErr(fmt.Errorf("world and page are required")), nil, nil
		}
		page := strings.TrimSuffix(args.Page, ".md")
		body, err := vault.ReadWiki(vaultDir, args.World, page)
		if err != nil {
			return textErr(err), nil, nil
		}
		return textOK(body), nil, nil
	})

	type writeWikiArgs struct {
		World   string   `json:"world" jsonschema:"world directory name"`
		Page    string   `json:"page" jsonschema:"wiki page name without .md"`
		Content string   `json:"content,omitempty" jsonschema:"markdown body after frontmatter"`
		Title   string   `json:"title,omitempty" jsonschema:"page title (defaults to page name)"`
		Tags    []string `json:"tags,omitempty" jsonschema:"frontmatter tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "write_wiki",
		Description: "Create or update a wiki page; maintains created/updated frontmatter",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args writeWikiArgs) (*mcp.CallToolResult, any, error) {
		if args.World == "" || args.Page == "" {
			return textErr(fmt.Errorf("world and page are required")), nil, nil
		}
		page := strings.TrimSuffix(args.Page, ".md")
		err := vault.WriteWiki(vaultDir, args.World, page, vault.WriteWikiOptions{
			Title:   args.Title,
			Content: args.Content,
			Tags:    args.Tags,
			Source:  "mcp",
		})
		if err != nil {
			return textErr(err), nil, nil
		}
		body, err := vault.ReadWiki(vaultDir, args.World, page)
		if err != nil {
			return textErr(err), nil, nil
		}
		return textOK(body), nil, nil
	})
}

func textOK(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}

func textErr(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		IsError: true,
	}
}
