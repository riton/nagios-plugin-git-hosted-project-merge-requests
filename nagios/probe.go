package nagios

import (
	"time"

	"github.com/pkg/errors"
	"github.com/riton/nagiosplugin/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CheckerContextKey = "checker"
)

func ProbeCobraAdapter(cmd *cobra.Command, args []string) {
	checker := cmd.Context().Value(CheckerContextKey).(*nagiosplugin.Check)

	cfg := probeConfig{
		// TODO:
	}

	probe := nagiosProbe{
		Hostname: viper.GetString("host"),
		cfg:      cfg,
		nagCheck: checker,
	}

	defer probe.Finish()
	probe.Run()
}

type probeConfig struct {
	// TODO:
}

type nagiosProbe struct {
	Hostname string
	cfg      probeConfig
	nagCheck *nagiosplugin.Check
}

func (c *nagiosProbe) init() error {
	return nil
}

func (c nagiosProbe) Run() {
	if err := c.init(); err != nil {
		c.nagCheck.Exitf(nagiosplugin.UNKNOWN, errors.Wrap(err, "initializing nagios probe").Error())
	}

	c.checkHost(c.Hostname)
}

func (c nagiosProbe) Finish() {
	c.nagCheck.Finish()
}

func (c nagiosProbe) checkHost(host string) {
	start := time.Now()

	// TODO:

	log.WithFields(log.Fields{
		"host": host,
	}).Info("host is reachable")

	c.nagCheck.AddResultf(nagiosplugin.OK, "Host is reachable")

	durationValue, err := nagiosplugin.NewFloatPerfDatumValue(time.Since(start).Seconds())
	if err != nil {
		c.nagCheck.Exitf(nagiosplugin.UNKNOWN, errors.Wrap(err, "creating perfdata").Error())
	}

	c.nagCheck.AddPerfDatum("total_duration", "s", durationValue, nil, nil, nil, nil)
}
