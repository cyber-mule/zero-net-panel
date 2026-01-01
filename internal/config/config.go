package config

import (
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest"

	"github.com/zero-net-panel/zero-net-panel/pkg/cache"
	"github.com/zero-net-panel/zero-net-panel/pkg/database"
)

type Config struct {
	rest.RestConf

	Project     ProjectConfig    `json:"project" yaml:"Project"`
	Site        SiteConfig       `json:"site" yaml:"Site"`
	Database    database.Config  `json:"database" yaml:"Database"`
	Cache       cache.Config     `json:"cache" yaml:"Cache"`
	Kernel      KernelConfig     `json:"kernel" yaml:"Kernel"`
	CORS        CORSConfig       `json:"cors" yaml:"CORS"`
	Auth        AuthConfig       `json:"auth" yaml:"Auth"`
	Credentials CredentialConfig `json:"credentials" yaml:"Credentials"`
	Metrics     MetricsConfig    `json:"metrics" yaml:"Metrics"`
	Admin       AdminConfig      `json:"admin" yaml:"Admin"`
	Webhook     WebhookConfig    `json:"webhook" yaml:"Webhook"`
	GRPC        GRPCServerConfig `json:"grpcServer" yaml:"GRPCServer"`
}

type ProjectConfig struct {
	Name        string `json:"name" yaml:"Name"`
	Description string `json:"description" yaml:"Description"`
	Version     string `json:"version" yaml:"Version"`
}

type SiteConfig struct {
	Name    string `json:"name" yaml:"Name"`
	LogoURL string `json:"logoUrl" yaml:"LogoURL"`
}

type KernelConfig struct {
	DefaultProtocol    string           `json:"defaultProtocol" yaml:"DefaultProtocol"`
	HTTP               KernelHTTPConfig `json:"http" yaml:"HTTP"`
	GRPC               KernelGRPCConfig `json:"grpc" yaml:"GRPC"`
	StatusPollInterval time.Duration    `json:"statusPollInterval" yaml:"StatusPollInterval"`
	StatusPollBackoff  KernelBackoff    `json:"statusPollBackoff" yaml:"StatusPollBackoff"`
}

type KernelHTTPConfig struct {
	BaseURL string        `json:"baseUrl" yaml:"BaseURL"`
	Token   string        `json:"token" yaml:"Token"`
	Timeout time.Duration `json:"timeout" yaml:"Timeout"`
}

type KernelGRPCConfig struct {
	Endpoint string        `json:"endpoint" yaml:"Endpoint"`
	TLSCert  string        `json:"tlsCert" yaml:"TLSCert"`
	Timeout  time.Duration `json:"timeout" yaml:"Timeout"`
}

// Normalize applies defaults for kernel configuration.
func (k *KernelConfig) Normalize() {
	k.DefaultProtocol = strings.TrimSpace(k.DefaultProtocol)
	if k.StatusPollInterval < 0 {
		k.StatusPollInterval = 0
	}
	k.StatusPollBackoff.Normalize(k.StatusPollInterval)
}

// CORSConfig configures cross-origin access for the HTTP API.
type CORSConfig struct {
	Enabled      bool     `json:"enabled" yaml:"Enabled"`
	AllowOrigins []string `json:"allowOrigins" yaml:"AllowOrigins"`
	AllowHeaders []string `json:"allowHeaders" yaml:"AllowHeaders"`
}

// Normalize applies defaults for CORS config.
func (c *CORSConfig) Normalize() {
	c.AllowOrigins = normalizeStringList(c.AllowOrigins, false)
	c.AllowHeaders = normalizeStringList(c.AllowHeaders, true)
	if c.Enabled && len(c.AllowOrigins) == 0 {
		c.AllowOrigins = []string{"*"}
	}
}

// KernelBackoff defines retry backoff behavior for kernel status polling.
type KernelBackoff struct {
	Enabled     bool          `json:"enabled" yaml:"Enabled"`
	MaxInterval time.Duration `json:"maxInterval" yaml:"MaxInterval"`
	Multiplier  float64       `json:"multiplier" yaml:"Multiplier"`
	Jitter      float64       `json:"jitter" yaml:"Jitter"`
}

// Normalize applies defaults for backoff settings.
func (b *KernelBackoff) Normalize(base time.Duration) {
	if base <= 0 {
		b.Enabled = false
		return
	}
	if !b.Enabled {
		return
	}
	if b.Multiplier <= 1 {
		b.Multiplier = 2
	}
	if b.Jitter < 0 {
		b.Jitter = 0
	}
	if b.Jitter > 1 {
		b.Jitter = 1
	}
	if b.MaxInterval <= 0 {
		b.MaxInterval = base * 8
	}
	if b.MaxInterval < base {
		b.MaxInterval = base
	}
}

type AuthConfig struct {
	AccessSecret   string                   `json:"accessSecret" yaml:"AccessSecret"`
	AccessExpire   time.Duration            `json:"accessExpire" yaml:"AccessExpire"`
	RefreshSecret  string                   `json:"refreshSecret" yaml:"RefreshSecret"`
	RefreshExpire  time.Duration            `json:"refreshExpire" yaml:"RefreshExpire"`
	Registration   AuthRegistrationConfig   `json:"registration" yaml:"Registration"`
	Verification   AuthVerificationConfig   `json:"verification" yaml:"Verification"`
	PasswordReset  AuthPasswordResetConfig  `json:"passwordReset" yaml:"PasswordReset"`
	PasswordPolicy AuthPasswordPolicyConfig `json:"passwordPolicy" yaml:"PasswordPolicy"`
	Lockout        AuthLockoutConfig        `json:"lockout" yaml:"Lockout"`
	Email          AuthEmailConfig          `json:"email" yaml:"Email"`
}

// Normalize applies defaults for auth.
func (a *AuthConfig) Normalize() {
	a.AccessSecret = strings.TrimSpace(a.AccessSecret)
	a.RefreshSecret = strings.TrimSpace(a.RefreshSecret)
	a.Registration.Normalize()
	a.Verification.Normalize()
	a.PasswordReset.Normalize()
	a.PasswordPolicy.Normalize()
	a.Lockout.Normalize()
	a.Email.Normalize()
}

type CredentialConfig struct {
	MasterKey string `json:"masterKey" yaml:"MasterKey"`
}

func (c *CredentialConfig) Normalize() {
	c.MasterKey = strings.TrimSpace(c.MasterKey)
}

// AuthRegistrationConfig controls signup behavior.
type AuthRegistrationConfig struct {
	Enabled                  bool     `json:"enabled" yaml:"Enabled"`
	InviteOnly               bool     `json:"inviteOnly" yaml:"InviteOnly"`
	InviteCodes              []string `json:"inviteCodes" yaml:"InviteCodes"`
	DefaultRoles             []string `json:"defaultRoles" yaml:"DefaultRoles"`
	RequireEmailVerification bool     `json:"requireEmailVerification" yaml:"RequireEmailVerification"`
}

// Normalize applies defaults for registration.
func (r *AuthRegistrationConfig) Normalize() {
	var codes []string
	for _, code := range r.InviteCodes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		codes = append(codes, code)
	}
	r.InviteCodes = codes
	var roles []string
	for _, role := range r.DefaultRoles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role == "" {
			continue
		}
		roles = append(roles, role)
	}
	r.DefaultRoles = roles
	if len(r.DefaultRoles) == 0 {
		r.DefaultRoles = []string{"user"}
	}
}

// AuthVerificationConfig defines verification code policies.
type AuthVerificationConfig struct {
	CodeLength       int           `json:"codeLength" yaml:"CodeLength"`
	CodeTTL          time.Duration `json:"codeTTL" yaml:"CodeTTL"`
	SendCooldown     time.Duration `json:"sendCooldown" yaml:"SendCooldown"`
	SendLimitPerHour int           `json:"sendLimitPerHour" yaml:"SendLimitPerHour"`
}

// Normalize applies defaults for verification codes.
func (v *AuthVerificationConfig) Normalize() {
	if v.CodeLength <= 0 {
		v.CodeLength = 6
	}
	if v.CodeTTL <= 0 {
		v.CodeTTL = 15 * time.Minute
	}
	if v.SendCooldown < 0 {
		v.SendCooldown = 0
	}
	if v.SendCooldown == 0 {
		v.SendCooldown = time.Minute
	}
	if v.SendLimitPerHour < 0 {
		v.SendLimitPerHour = 0
	}
	if v.SendLimitPerHour == 0 {
		v.SendLimitPerHour = 5
	}
}

// AuthPasswordResetConfig defines password reset policies.
type AuthPasswordResetConfig struct {
	CodeLength       int           `json:"codeLength" yaml:"CodeLength"`
	CodeTTL          time.Duration `json:"codeTTL" yaml:"CodeTTL"`
	SendCooldown     time.Duration `json:"sendCooldown" yaml:"SendCooldown"`
	SendLimitPerHour int           `json:"sendLimitPerHour" yaml:"SendLimitPerHour"`
}

// Normalize applies defaults for reset codes.
func (p *AuthPasswordResetConfig) Normalize() {
	if p.CodeLength <= 0 {
		p.CodeLength = 6
	}
	if p.CodeTTL <= 0 {
		p.CodeTTL = 15 * time.Minute
	}
	if p.SendCooldown < 0 {
		p.SendCooldown = 0
	}
	if p.SendCooldown == 0 {
		p.SendCooldown = time.Minute
	}
	if p.SendLimitPerHour < 0 {
		p.SendLimitPerHour = 0
	}
	if p.SendLimitPerHour == 0 {
		p.SendLimitPerHour = 5
	}
}

// AuthPasswordPolicyConfig defines password complexity requirements.
type AuthPasswordPolicyConfig struct {
	MinLength      int  `json:"minLength" yaml:"MinLength"`
	MaxLength      int  `json:"maxLength" yaml:"MaxLength"`
	RequireUpper   bool `json:"requireUpper" yaml:"RequireUpper"`
	RequireLower   bool `json:"requireLower" yaml:"RequireLower"`
	RequireDigit   bool `json:"requireDigit" yaml:"RequireDigit"`
	RequireSpecial bool `json:"requireSpecial" yaml:"RequireSpecial"`
}

// Normalize applies defaults for password policy.
func (p *AuthPasswordPolicyConfig) Normalize() {
	if p.MinLength <= 0 {
		p.MinLength = 8
	}
	if p.MaxLength < 0 {
		p.MaxLength = 0
	}
	if p.MaxLength > 0 && p.MaxLength < p.MinLength {
		p.MaxLength = p.MinLength
	}
}

// AuthLockoutConfig controls login lockouts.
type AuthLockoutConfig struct {
	MaxAttempts  int           `json:"maxAttempts" yaml:"MaxAttempts"`
	LockDuration time.Duration `json:"lockDuration" yaml:"LockDuration"`
}

// Normalize applies defaults for lockout.
func (l *AuthLockoutConfig) Normalize() {
	if l.MaxAttempts <= 0 {
		l.MaxAttempts = 5
	}
	if l.LockDuration <= 0 {
		l.LockDuration = 15 * time.Minute
	}
}

// AuthEmailConfig controls email delivery.
type AuthEmailConfig struct {
	Provider string         `json:"provider" yaml:"Provider"`
	From     string         `json:"from" yaml:"From"`
	SMTP     AuthSMTPConfig `json:"smtp" yaml:"SMTP"`
}

// Normalize applies defaults for email.
func (e *AuthEmailConfig) Normalize() {
	e.Provider = strings.ToLower(strings.TrimSpace(e.Provider))
	if e.Provider == "" {
		e.Provider = "log"
	}
	e.From = strings.TrimSpace(e.From)
	if e.From == "" {
		e.From = "no-reply@localhost"
	}
	e.SMTP.Normalize()
}

// AuthSMTPConfig defines SMTP delivery settings.
type AuthSMTPConfig struct {
	Host     string `json:"host" yaml:"Host"`
	Port     int    `json:"port" yaml:"Port"`
	Username string `json:"username" yaml:"Username"`
	Password string `json:"password" yaml:"Password"`
	UseTLS   bool   `json:"useTLS" yaml:"UseTLS"`
}

// Normalize applies defaults for SMTP.
func (s *AuthSMTPConfig) Normalize() {
	s.Host = strings.TrimSpace(s.Host)
	s.Username = strings.TrimSpace(s.Username)
	s.Password = strings.TrimSpace(s.Password)
	if s.Port < 0 {
		s.Port = 0
	}
	if s.Port == 0 {
		s.Port = 587
	}
}

type MetricsConfig struct {
	Enable   bool   `json:"enable" yaml:"Enable"`
	Path     string `json:"path" yaml:"Path"`
	ListenOn string `json:"listenOn" yaml:"ListenOn"`
}

// Normalize trims the path/listener and applies defaults.
func (m *MetricsConfig) Normalize() {
	path := strings.TrimSpace(m.Path)
	if path == "" {
		path = "/metrics"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	m.Path = path

	m.ListenOn = strings.TrimSpace(m.ListenOn)
	if !m.Enable {
		m.ListenOn = ""
	}
}

// Enabled returns whether metrics export is enabled.
func (m MetricsConfig) Enabled() bool {
	return m.Enable
}

// Standalone reports whether metrics should be served on an independent listener.
func (m MetricsConfig) Standalone() bool {
	return m.Enable && m.ListenOn != ""
}

// Normalize trims user-facing site values.
func (s *SiteConfig) Normalize() {
	s.Name = strings.TrimSpace(s.Name)
	s.LogoURL = strings.TrimSpace(s.LogoURL)
}

// AdminConfig 控制管理端路由相关配置。
type AdminConfig struct {
	RoutePrefix string            `json:"routePrefix" yaml:"RoutePrefix"`
	Access      AdminAccessConfig `json:"access" yaml:"Access"`
}

// Normalize 统一前缀写法并设置默认值。
func (a *AdminConfig) Normalize() {
	prefix := strings.TrimSpace(a.RoutePrefix)
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		prefix = "admin"
	}
	a.RoutePrefix = prefix
	a.Access.Normalize()
}

// APIBasePath 返回管理端挂载的完整 API 前缀。
func (a AdminConfig) APIBasePath() string {
	if a.RoutePrefix == "" {
		return "/api/v1"
	}
	return "/api/v1/" + a.RoutePrefix
}

// AdminAccessConfig controls admin ingress policies.
type AdminAccessConfig struct {
	AllowCIDRs         []string `json:"allowCidrs" yaml:"AllowCIDRs"`
	RateLimitPerMinute int      `json:"rateLimitPerMinute" yaml:"RateLimitPerMinute"`
	Burst              int      `json:"burst" yaml:"Burst"`
}

// Normalize applies sane defaults.
func (a *AdminAccessConfig) Normalize() {
	if a.RateLimitPerMinute < 0 {
		a.RateLimitPerMinute = 0
	}
	if a.Burst < 0 {
		a.Burst = 0
	}
	if a.RateLimitPerMinute > 0 && a.Burst == 0 {
		a.Burst = a.RateLimitPerMinute / 6
		if a.Burst < 1 {
			a.Burst = 1
		}
	}
}

// WebhookConfig controls external callback validation.
type WebhookConfig struct {
	AllowCIDRs  []string            `json:"allowCidrs" yaml:"AllowCIDRs"`
	SharedToken string              `json:"sharedToken" yaml:"SharedToken"`
	Stripe      StripeWebhookConfig `json:"stripe" yaml:"Stripe"`
}

// Normalize applies defaults.
func (w *WebhookConfig) Normalize() {
	w.SharedToken = strings.TrimSpace(w.SharedToken)
	w.Stripe.Normalize()
}

// StripeWebhookConfig controls Stripe webhook signature verification.
type StripeWebhookConfig struct {
	SigningSecret    string `json:"signingSecret" yaml:"SigningSecret"`
	ToleranceSeconds int    `json:"toleranceSeconds" yaml:"ToleranceSeconds"`
}

// Normalize applies defaults for Stripe.
func (s *StripeWebhookConfig) Normalize() {
	s.SigningSecret = strings.TrimSpace(s.SigningSecret)
	if s.ToleranceSeconds <= 0 {
		s.ToleranceSeconds = 300
	}
}

// GRPCServerConfig 控制内建 gRPC 服务监听配置。
type GRPCServerConfig struct {
	Enable     *bool  `json:"enable" yaml:"Enable"`
	ListenOn   string `json:"listenOn" yaml:"ListenOn"`
	Reflection *bool  `json:"reflection" yaml:"Reflection"`
}

// Normalize 设置默认监听地址与开关。
func (g *GRPCServerConfig) Normalize() {
	if g.Enable == nil {
		g.Enable = boolPtr(true)
	}
	if g.Reflection == nil {
		g.Reflection = boolPtr(true)
	}
	if g.Enabled() && strings.TrimSpace(g.ListenOn) == "" {
		g.ListenOn = "0.0.0.0:8890"
	}
}

// Enabled 返回 gRPC 服务是否启用（默认为 true）。
func (g GRPCServerConfig) Enabled() bool {
	if g.Enable == nil {
		return true
	}
	return *g.Enable
}

// SetEnabled 修改 gRPC 启用状态。
func (g *GRPCServerConfig) SetEnabled(enabled bool) {
	g.Enable = boolPtr(enabled)
}

// ReflectionEnabled 返回是否开放 gRPC reflection（默认为 true）。
func (g GRPCServerConfig) ReflectionEnabled() bool {
	if g.Reflection == nil {
		return true
	}
	return *g.Reflection
}

// Normalize 将配置补齐默认值。
func (c *Config) Normalize() {
	c.Project.Name = strings.TrimSpace(c.Project.Name)
	c.Project.Description = strings.TrimSpace(c.Project.Description)
	c.Project.Version = strings.TrimSpace(c.Project.Version)
	c.Site.Normalize()
	if c.Site.Name == "" {
		c.Site.Name = c.Project.Name
	}
	c.Kernel.Normalize()
	c.CORS.Normalize()
	c.Auth.Normalize()
	c.Credentials.Normalize()
	c.Metrics.Normalize()
	c.Admin.Normalize()
	c.Webhook.Normalize()
	c.GRPC.Normalize()
	c.Middlewares.Prometheus = c.Metrics.Enabled()
	c.Middlewares.Metrics = c.Metrics.Enabled()
}

func boolPtr(v bool) *bool {
	return &v
}

func normalizeStringList(values []string, lower bool) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if lower {
			trimmed = strings.ToLower(trimmed)
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
