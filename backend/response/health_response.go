package response

type HealthResponse struct {
	Status string                    `json:"status"`
	Uptime string                    `json:"uptime"`
	Checks map[string]ComponentCheck `json:"checks"`
	System SystemInfo                `json:"system"`
}

type ComponentCheck struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SystemInfo struct {
	Version     string `json:"version"`
	Env         string `json:"env"`
	GoRoutines  int    `json:"go_routines"`
	MemoryAlloc string `json:"memory_alloc"`
	MemorySys   string `json:"memory_sys"`
}
