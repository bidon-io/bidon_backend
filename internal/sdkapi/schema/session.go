package schema

type Session struct {
	ID                        string   `json:"id" validate:"required,uuid"`
	LaunchTS                  int      `json:"launch_ts" validate:"required"`
	LaunchMonotonicTS         int      `json:"launch_monotonic_ts" validate:"required"`
	StartTS                   int      `json:"start_ts" validate:"required"`
	StartMonotonicTS          int      `json:"start_monotonic_ts" validate:"required"`
	TS                        int      `json:"ts" validate:"required"`
	MonotonicTS               int      `json:"monotonic_ts" validate:"required"`
	MemoryWarningsTS          []int    `json:"memory_warnings_ts" validate:"required"`
	MemoryWarningsMonotonicTS []int    `json:"memory_warnings_monotonic_ts" validate:"required"`
	RAMUsed                   int      `json:"ram_used" validate:"required"`
	RAMSize                   int      `json:"ram_size" validate:"required"`
	StorageFree               int      `json:"storage_free"`
	StorageUsed               int      `json:"storage_used"`
	Battery                   float64  `json:"battery" validate:"required"`
	CPUUsage                  *float64 `json:"cpu_usage" validate:"required"`
}

func (s Session) Uptime() int {
	return s.MonotonicTS - s.StartMonotonicTS
}
