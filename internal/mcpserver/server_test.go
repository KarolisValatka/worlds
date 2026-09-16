package mcpserver

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/KarolisValatka/worlds/internal/vault"
)

func TestMCPTools(t *testing.T) {
	dir := t.TempDir()
	v := filepath.Join(dir, "v")
	if err := vault.Init(v); err != nil {
		t.Fatal(err)
	}
	server, err := New(v)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{"list_worlds", "read_brain", "read_world", "read_wiki", "write_wiki"} {
		if !names[want] {
			t.Fatalf("missing tool %s in %#v", want, names)
		}
	}

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list_worlds", Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	text := toolText(res)
	if !strings.Contains(text, "learn") {
		t.Fatalf("list_worlds=%q", text)
	}

	res, err = cs.CallTool(ctx, &mcp.CallToolParams{Name: "read_brain", Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(toolText(res), "Brain index") {
		t.Fatalf("read_brain=%q", toolText(res))
	}

	res, err = cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "read_world",
		Arguments: map[string]any{
			"world":     "learn",
			"list_wiki": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(toolText(res), "hello") {
		t.Fatalf("read_world=%q", toolText(res))
	}

	res, err = cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "write_wiki",
		Arguments: map[string]any{
			"world":   "learn",
			"page":    "from-mcp",
			"title":   "from mcp",
			"content": "# from mcp\n\ncreated via tool\n",
			"tags":    []any{"mcp"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("write_wiki error: %s", toolText(res))
	}
	body, err := vault.ReadWiki(v, "learn", "from-mcp")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "source: mcp") || !strings.Contains(body, "created via tool") {
		t.Fatalf("wiki=%s", body)
	}

	res, err = cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "read_wiki",
		Arguments: map[string]any{
			"world": "learn",
			"page":  "from-mcp",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(toolText(res), "from mcp") {
		t.Fatalf("read_wiki=%q", toolText(res))
	}
}

func toolText(res *mcp.CallToolResult) string {
	var b strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			b.WriteString(tc.Text)
		}
	}
	return b.String()
}
