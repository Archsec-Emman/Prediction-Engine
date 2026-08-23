package engine

import "fmt"

var DepGit = Dep{Bin: "git", Desc: "clone upstream repositories"}
var DepPython = Dep{Bin: "python", Args: []string{"--version"}, Desc: "run Python-based engines"}
var DepPip = Dep{Bin: "pip", Args: []string{"--version"}, Desc: "install Python packages"}
var DepNode = Dep{Bin: "node", Args: []string{"--version"}, Desc: "run the Polyseer web UI"}
var DepNpm = Dep{Bin: "npm", Args: []string{"--version"}, Desc: "install Polyseer dependencies"}

func Registry() []Engine {
	return []Engine{
		{
			Name:        "polyseer",
			Description: "AI-driven prediction-market research: multi-agent analysis, evidence grading, Bayesian probability aggregation. Web UI.",
			Upstream:    "https://github.com/yorkeccak/Polyseer",
			Credit:      "Polyseer by yorkeccak",
			DirName:     "polyseer",
			Requires:    []Dep{DepGit, DepNode, DepNpm},
			Install: []Step{
				{Cmd: []string{"npm", "install"}},
				{Cmd: []string{"npm", "run", "build"}},
			},
			DefaultRun: []string{"npm", "run", "dev"},
			Mode:       RunServer,
		},
		{
			Name:        "papertrader",
			Description: "Polymarket paper-trading simulator for AI agents: MCP server, level-by-level order book execution, fee and slippage model.",
			Upstream:    "https://github.com/agent-next/polymarket-paper-trader",
			Credit:      "polymarket-paper-trader by agent-next",
			DirName:     "papertrader",
			Requires:    []Dep{DepGit, DepPython, DepPip},
			Install: []Step{
				{Cmd: []string{"pip", "install", "-e", "."}},
			},
			DefaultRun: []string{"python", "-m", "pe_trader.cli"},
			Mode:       RunCLI,
		},
		{
			Name:        "backtest",
			Description: "Production-grade strategy backtesting on NautilusTrader: Polymarket adapter, order book replay, Optuna optimisation.",
			Upstream:    "https://github.com/evan-kolberg/prediction-market-backtesting",
			Credit:      "prediction-market-backtesting by evan-kolberg (MIT + LGPL-3.0-or-later)",
			DirName:     "backtest",
			Requires:    []Dep{DepGit, DepPython, DepPip},
			Install: []Step{
				{Cmd: []string{"pip", "install", "-e", "."}},
			},
			DefaultRun: nil,
			Mode:       RunCLI,
		},
	}
}

func ByName(name string) (*Engine, error) {
	for i := range Registry() {
		if Registry()[i].Name == name {
			return &Registry()[i], nil
		}
	}
	return nil, fmt.Errorf("unknown engine %q (known: polyseer, papertrader, backtest)", name)
}
