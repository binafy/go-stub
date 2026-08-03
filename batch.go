package stub

import (
	"fmt"
	"io/fs"
)

// Job describes a single source-to-destination generation within a batch.
//
// When FS is nil the source is read from the operating system filesystem;
// otherwise it is read from FS (for example an embed.FS). Opts are applied
// after the shared options passed to GenerateJobs, so a job can override them.
type Job struct {
	FS   fs.FS
	Src  string
	Dst  string
	Opts []Option
}

// GenerateJobs runs a slice of jobs in order, applying shared to every job
// before that job's own Opts. It stops and returns on the first failing job,
// wrapping the error with the job's index and paths.
func GenerateJobs(jobs []Job, shared ...Option) error {
	for i, j := range jobs {
		opts := make([]Option, 0, len(shared)+len(j.Opts))
		opts = append(opts, shared...)
		opts = append(opts, j.Opts...)

		var err error
		if j.FS != nil {
			err = GenerateFS(j.FS, j.Src, j.Dst, opts...)
		} else {
			err = Generate(j.Src, j.Dst, opts...)
		}
		if err != nil {
			return fmt.Errorf("stub: job %d (%q -> %q): %w", i, j.Src, j.Dst, err)
		}
	}

	return nil
}
