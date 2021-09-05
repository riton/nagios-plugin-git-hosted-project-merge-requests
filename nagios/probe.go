package nagios

import (
	"time"

	"github.com/pkg/errors"
	"github.com/riton/nagiosplugin/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	CheckerContextKey = "checker"
)

func ProbeCobraAdapter(cmd *cobra.Command, args []string, cfg ProbeConfig) {
	checker := cmd.Context().Value(CheckerContextKey).(*nagiosplugin.Check)

	probe := nagiosProbe{
		cfg:      cfg,
		nagCheck: checker,
	}

	timer := time.NewTimer(cfg.Timeout)
	doneChan := make(chan struct{})

	// Do the real work in a dedicated goroutine
	// so that the main one will handle the global timeout
	go func() {
		defer func() {
			doneChan <- struct{}{}
			probe.Finish()
		}()
		probe.Run()
	}()

	for {
		select {
		case <-timer.C:
			checker.Unknownf("timeout after %s", cfg.Timeout)
		case <-doneChan:
			timer.Stop()
			break
		}
	}
}

type nagiosProbe struct {
	Hostname string
	cfg      ProbeConfig
	nagCheck *nagiosplugin.Check
}

func (c *nagiosProbe) init() error {
	return nil
}

func (c nagiosProbe) Run() {
	if err := c.init(); err != nil {
		c.nagCheck.Exitf(nagiosplugin.UNKNOWN, errors.Wrap(err, "initializing nagios probe").Error())
	}

	// TODO:
	if c.cfg.GitProvider != GitlabGitProvider {
		c.nagCheck.Criticalf("git provider %s is not supported yet", c.cfg.GitProvider)
	}

	mrChecker, err := newGitlabProjectMRChecker(c.cfg.APIEndpoint, c.cfg.APIToken)
	if err != nil {
		c.nagCheck.Unknownf("fail to initialize gitlab checker: %s", err)
	}

	c.checkMergeRequests(mrChecker)
}

func (c nagiosProbe) Finish() {
	c.nagCheck.Finish()
}

func (c nagiosProbe) checkMergeRequests(mrChecker GitMergeRequestChecker) {
	start := time.Now()

	mr, err := mrChecker.CheckMergeRequests(c.cfg.Project, c.cfg.TargetBranch)
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"project":       c.cfg.Project,
			"api-endpoint":  c.cfg.APIEndpoint,
			"target-branch": c.cfg.TargetBranch,
		}).Error("fail to check for merge requests")
		c.nagCheck.Criticalf("fail to check for merge requests: %s", err)
	}

	log.WithFields(log.Fields{
		"merge-requests": mr,
	}).Debug("merge requests fetched successfully")

	durationValue, err := nagiosplugin.NewFloatPerfDatumValue(time.Since(start).Seconds())
	if err != nil {
		c.nagCheck.Exitf(nagiosplugin.UNKNOWN, errors.Wrap(err, "creating perfdata").Error())
	}
	c.nagCheck.AddPerfDatum("total_duration", "s", durationValue, nil, nil, nil, nil)

	totalMr, _ := nagiosplugin.NewFloatPerfDatumValue(float64(len(mr)))
	c.nagCheck.AddPerfDatum("opened_merge_requests", "", totalMr, nil, nil, nil, nil)

	if len(mr) == 0 {
		c.nagCheck.Exitf(nagiosplugin.OK, "No opened merge requests")
		return
	}

	c.nagCheck.AddResult(nagiosplugin.OK, "No merge requests too old")

	var oldestMrDuration time.Duration
	for _, cmr := range mr {
		tSinceLastUpdate := time.Since(cmr.UpdatedAt)
		if tSinceLastUpdate >= c.cfg.CriticalLastUpdateDelay {
			c.nagCheck.AddResultf(nagiosplugin.CRITICAL, "Merge request %d last activity was %s ago", cmr.ID, tSinceLastUpdate)
		} else if tSinceLastUpdate >= c.cfg.WarningLastUpdateDelay {
			c.nagCheck.AddResultf(nagiosplugin.WARNING, "Merge request %d last activity was %s ago", cmr.ID, tSinceLastUpdate)
		}

		// keep track of our oldest merge-request for perfdata
		if tSinceLastUpdate > oldestMrDuration {
			oldestMrDuration = tSinceLastUpdate
		}
	}

	oldestMrPerfDatum, _ := nagiosplugin.NewFloatPerfDatumValue(oldestMrDuration.Seconds())
	c.nagCheck.AddPerfDatum("oldest_merge_request", "s", oldestMrPerfDatum, nil, nil, nil, nil)
}
