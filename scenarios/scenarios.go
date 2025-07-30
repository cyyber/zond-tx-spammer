package scenarios

import (
	"github.com/theQRL/qrl-tx-spammer/scenarios/deploytx"
	"github.com/theQRL/qrl-tx-spammer/scenarios/eoatx"
	"github.com/theQRL/qrl-tx-spammer/scenarios/gasburnertx"
	"github.com/theQRL/qrl-tx-spammer/scenarios/sqrctx"
	"github.com/theQRL/qrl-tx-spammer/scenarios/wallets"
	"github.com/theQRL/qrl-tx-spammer/scenariotypes"
)

var Scenarios map[string]func() scenariotypes.Scenario = map[string]func() scenariotypes.Scenario{
	"eoatx":       eoatx.NewScenario,
	"sqrctx":      sqrctx.NewScenario,
	"deploytx":    deploytx.NewScenario,
	"gasburnertx": gasburnertx.NewScenario,

	"wallets": wallets.NewScenario,
}
