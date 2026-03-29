package types

import "time"

type VerbChecksum struct {
	Verb       RunnerVerb `json:"verb"`
	Checksum   string     `json:"checksum"`
	ExecutedAt time.Time  `json:"executed_at"`
}

type ChecksumCache map[string]VerbChecksum
