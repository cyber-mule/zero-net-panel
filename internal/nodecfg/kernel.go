package nodecfg

const (
	DefaultKernelProtocol                     = "http"
	DefaultKernelHTTPTimeoutSeconds           = 5
	DefaultKernelStatusPollIntervalSeconds    = 30
	DefaultKernelStatusPollBackoffEnabled     = true
	DefaultKernelStatusPollBackoffMaxIntervalSeconds = 300
	DefaultKernelStatusPollBackoffMultiplier  = 2
	DefaultKernelStatusPollBackoffJitter      = 0.2
	DefaultKernelOfflineProbeMaxIntervalSeconds = 0
)

// KernelBackoffConfig describes per-node poll backoff behavior.
type KernelBackoffConfig struct {
	Enabled             bool
	MaxIntervalSeconds  int
	Multiplier          float64
	Jitter              float64
}
