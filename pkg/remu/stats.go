package remu

import (
	"context"
	"strings"
)

// InfoBasic defines redis info basic stats
type InfoBasic struct {
	Version          string
	UptimeInDays     string
	UptimeInSeconds  string
	MemoryUsage      string
	MemoryUsagePeak  string
	ConnectedClients string
}

// RedisInfoBasic return redis info basic stats
func (r *Remu) RedisInfoBasic() (*InfoBasic, error) {
	res, err := r.Client.Info(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return parseInfoBasic(res), nil
}

// Helper function to parsing redis info string to struct
func parseInfoBasic(s string) *InfoBasic {
	infoBasic := &InfoBasic{}
	lines := strings.Split(s, "\r\n")

	for _, line := range lines {
		kv := strings.Split(line, ":")
		if len(kv) == 2 {
			switch true {
			case kv[0] == "redis_version":
				infoBasic.Version = kv[1]

			case kv[0] == "uptime_in_days":
				infoBasic.UptimeInDays = kv[1]

			case kv[0] == "uptime_in_seconds":
				infoBasic.UptimeInSeconds = kv[1]

			case kv[0] == "used_memory_human":
				infoBasic.MemoryUsage = kv[1]

			case kv[0] == "used_memory_peak_human":
				infoBasic.MemoryUsagePeak = kv[1]

			case kv[0] == "connected_clients":
				infoBasic.ConnectedClients = kv[1]
			}
		}
	}

	return infoBasic
}
