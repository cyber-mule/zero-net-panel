package client

import "strings"

// SourceZeroCore marks built-in rules aligned with the Zero Core integration.
const SourceZeroCore = "zero-core"

// Rule describes a client type and the User-Agent tokens used to detect it.
type Rule struct {
	Type   string
	Label  string
	Tokens []string
	Source string
}

var defaultRules = []Rule{
	{
		Type:   "zero-core",
		Label:  "Zero Core",
		Tokens: []string{"zero-core", "zero core", "zerocore"},
		Source: SourceZeroCore,
	},
	{
		Type:   "sing-box",
		Label:  "Sing-box",
		Tokens: []string{"sing-box", "singbox"},
		Source: SourceZeroCore,
	},
	{
		Type:   "clash",
		Label:  "Clash",
		Tokens: []string{"clash", "clash.meta", "clashmeta", "clash-verge", "clash verge", "clashx", "clashforwindows", "mihomo", "mihomo-party"},
		Source: SourceZeroCore,
	},
	{
		Type:   "surge",
		Label:  "Surge",
		Tokens: []string{"surge"},
		Source: SourceZeroCore,
	},
	{
		Type:   "quantumult",
		Label:  "Quantumult",
		Tokens: []string{"quantumult", "quantumult x", "quantumultx"},
		Source: SourceZeroCore,
	},
	{
		Type:   "stash",
		Label:  "Stash",
		Tokens: []string{"stash"},
		Source: SourceZeroCore,
	},
	{
		Type:   "shadowrocket",
		Label:  "Shadowrocket",
		Tokens: []string{"shadowrocket"},
		Source: SourceZeroCore,
	},
	{
		Type:   "loon",
		Label:  "Loon",
		Tokens: []string{"loon"},
		Source: SourceZeroCore,
	},
	{
		Type:   "hiddify",
		Label:  "Hiddify",
		Tokens: []string{"hiddify"},
		Source: SourceZeroCore,
	},
	{
		Type:   "v2rayn",
		Label:  "v2rayN",
		Tokens: []string{"v2rayn"},
		Source: SourceZeroCore,
	},
	{
		Type:   "v2rayng",
		Label:  "v2rayNG",
		Tokens: []string{"v2rayng"},
		Source: SourceZeroCore,
	},
	{
		Type:   "nekobox",
		Label:  "NekoBox",
		Tokens: []string{"nekobox", "neko box"},
		Source: SourceZeroCore,
	},
	{
		Type:   "kitsunebi",
		Label:  "Kitsunebi",
		Tokens: []string{"kitsunebi"},
		Source: SourceZeroCore,
	},
	{
		Type:   "potatso",
		Label:  "Potatso",
		Tokens: []string{"potatso"},
		Source: SourceZeroCore,
	},
	{
		Type:   "surfboard",
		Label:  "Surfboard",
		Tokens: []string{"surfboard"},
		Source: SourceZeroCore,
	},
}

// Rules returns a copy of the built-in client rules.
func Rules() []Rule {
	result := make([]Rule, 0, len(defaultRules))
	for _, item := range defaultRules {
		result = append(result, Rule{
			Type:   item.Type,
			Label:  item.Label,
			Tokens: append([]string(nil), item.Tokens...),
			Source: item.Source,
		})
	}
	return result
}

// DetectClientType resolves the client type from User-Agent using built-in rules.
func DetectClientType(userAgent string) string {
	return DetectClientTypeWithRules(userAgent, defaultRules)
}

// DetectClientTypeWithRules resolves the client type from User-Agent using provided rules.
func DetectClientTypeWithRules(userAgent string, rules []Rule) string {
	ua := strings.ToLower(strings.TrimSpace(userAgent))
	if ua == "" {
		return ""
	}

	for _, item := range rules {
		for _, token := range item.Tokens {
			if token == "" {
				continue
			}
			if strings.Contains(ua, token) {
				return item.Type
			}
		}
	}
	return ""
}
