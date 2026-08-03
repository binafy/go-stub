package stub_test

import (
	"testing"

	stub "github.com/binafy/go-stub"
)

func TestCaseConversions(t *testing.T) {
	tests := []struct {
		in                                          string
		snake, screaming, kebab, pascal, camelCased string
	}{
		{
			in:         "UserName",
			snake:      "user_name",
			screaming:  "USER_NAME",
			kebab:      "user-name",
			pascal:     "UserName",
			camelCased: "userName",
		},
		{
			in:         "user name",
			snake:      "user_name",
			screaming:  "USER_NAME",
			kebab:      "user-name",
			pascal:     "UserName",
			camelCased: "userName",
		},
		{
			in:         "user_profile_id",
			snake:      "user_profile_id",
			screaming:  "USER_PROFILE_ID",
			kebab:      "user-profile-id",
			pascal:     "UserProfileId",
			camelCased: "userProfileId",
		},
		{
			in:         "HTTPServer",
			snake:      "http_server",
			screaming:  "HTTP_SERVER",
			kebab:      "http-server",
			pascal:     "HttpServer",
			camelCased: "httpServer",
		},
		{
			in:         "order-item-v2",
			snake:      "order_item_v2",
			screaming:  "ORDER_ITEM_V2",
			kebab:      "order-item-v2",
			pascal:     "OrderItemV2",
			camelCased: "orderItemV2",
		},
		{
			in:         "",
			snake:      "",
			screaming:  "",
			kebab:      "",
			pascal:     "",
			camelCased: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := stub.ToSnake(tt.in); got != tt.snake {
				t.Errorf("ToSnake(%q) = %q, want %q", tt.in, got, tt.snake)
			}
			if got := stub.ToScreamingSnake(tt.in); got != tt.screaming {
				t.Errorf("ToScreamingSnake(%q) = %q, want %q", tt.in, got, tt.screaming)
			}
			if got := stub.ToKebab(tt.in); got != tt.kebab {
				t.Errorf("ToKebab(%q) = %q, want %q", tt.in, got, tt.kebab)
			}
			if got := stub.ToPascal(tt.in); got != tt.pascal {
				t.Errorf("ToPascal(%q) = %q, want %q", tt.in, got, tt.pascal)
			}
			if got := stub.ToCamel(tt.in); got != tt.camelCased {
				t.Errorf("ToCamel(%q) = %q, want %q", tt.in, got, tt.camelCased)
			}
		})
	}
}

// TestCaseHelpersInReplaceMap shows the intended use: deriving placeholder
// values from a single base name.
func TestCaseHelpersInReplaceMap(t *testing.T) {
	base := "user profile"
	path := writeTemp(t, "{{ Pascal }}/{{ snake }}")
	got, err := stub.Render(path, stub.WithReplaces(map[string]any{
		"Pascal": stub.ToPascal(base),
		"snake":  stub.ToSnake(base),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "UserProfile/user_profile" {
		t.Errorf("got %q", got)
	}
}
