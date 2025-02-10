package scenariotypes

import (
	"github.com/spf13/pflag"
	"github.com/theQRL/tx-spammer/tester"
)

type Scenario interface {
	Flags(flags *pflag.FlagSet) error
	Init(testerCfg *tester.TesterConfig) error
	Run(tester *tester.Tester) error
}
