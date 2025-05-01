package scenarios

import (
	"github.com/theQRL/zond-tx-spammer/scenarios/deploytx"
	"github.com/theQRL/zond-tx-spammer/scenarios/eoatx"
	"github.com/theQRL/zond-tx-spammer/scenarios/gasburnertx"
	"github.com/theQRL/zond-tx-spammer/scenarios/wallets"
	"github.com/theQRL/zond-tx-spammer/scenarios/zrctx"
	"github.com/theQRL/zond-tx-spammer/scenariotypes"
)

var Scenarios map[string]func() scenariotypes.Scenario = map[string]func() scenariotypes.Scenario{
	"eoatx":       eoatx.NewScenario,
	"zrctx":       zrctx.NewScenario,
	"deploytx":    deploytx.NewScenario,
	"gasburnertx": gasburnertx.NewScenario,

	"wallets": wallets.NewScenario,
}
