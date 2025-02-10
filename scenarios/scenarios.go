package scenarios

import (
	deploydestruct "github.com/theQRL/tx-spammer/scenarios/deploy-destruct"
	"github.com/theQRL/tx-spammer/scenarios/deploytx"
	"github.com/theQRL/tx-spammer/scenarios/eoatx"
	"github.com/theQRL/tx-spammer/scenarios/erctx"
	"github.com/theQRL/tx-spammer/scenarios/gasburnertx"
	"github.com/theQRL/tx-spammer/scenarios/setcodetx"
	"github.com/theQRL/tx-spammer/scenarios/wallets"
	"github.com/theQRL/tx-spammer/scenariotypes"
)

var Scenarios map[string]func() scenariotypes.Scenario = map[string]func() scenariotypes.Scenario{
	"eoatx":           eoatx.NewScenario,
	"erctx":           erctx.NewScenario,
	"deploy-destruct": deploydestruct.NewScenario,
	"deploytx":        deploytx.NewScenario,
	"gasburnertx":     gasburnertx.NewScenario,
	"setcodetx":       setcodetx.NewScenario,

	"wallets": wallets.NewScenario,
}
