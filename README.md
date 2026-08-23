<p align="center">
  <img src="https://raw.githubusercontent.com/Archsec-Emman/Prediction-Engine/main/prediction-engine.png" alt="Prediction Engine" width="600">
</p>

# Prediction Engine

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/Archsec-Emman/Prediction-Engine/ci.yml?branch=main&label=build)](https://github.com/Archsec-Emman/Prediction-Engine/actions)

**A single launcher that installs and runs three open-source prediction-market tools in parallel — with honest exit codes and full attribution.**

> **What this is NOT:** this binary contains no trading logic, no AI, and no backtesting engine of its own. It is a dependency-checking orchestrator for three upstream projects by other authors. Earlier versions of this repository vendored those projects under the same name; that was misleading and has been removed.

## The three engines

| Engine | What it does | Upstream (all credit to) |
|--------|--------------|--------------------------|
| `polyseer` | Multi-agent AI market research: evidence grading, Bayesian probability aggregation, web UI | [Polyseer](https://github.com/yorkeccak/Polyseer) by yorkeccak |
| `papertrader` | Polymarket paper trading for agents: MCP server, order-book execution, fee/slippage model | [polymarket-paper-trader](https://github.com/agent-next/polymarket-paper-trader) by agent-next |
| `backtest` | NautilusTrader-based strategy backtesting: book replay, Optuna optimisation | [prediction-market-backtesting](https://github.com/evan-kolberg/prediction-market-backtesting) by evan-kolberg |

## Install

```bash
go install github.com/Archsec-Emman/Prediction-Engine@latest
# or build from source:
git clone https://github.com/Archsec-Emman/Prediction-Engine.git
cd Prediction-Engine && go build -o prediction-engine .
```

Zero dependencies beyond the Go standard library.

## Usage

```bash
# see engines, dependency status (git/node/python/pip), install state, credits
prediction-engine status

# clone upstream repos into ~/.prediction-engine and run their installers
prediction-engine install all
prediction-engine install polyseer

# run one or all engines (output streamed live per engine, summary at end)
prediction-engine run papertrader
prediction-engine run polyseer            # starts the web UI dev server
prediction-engine run all -- --help       # args after -- are passed through
```

Exit code reflects reality: if any requested engine fails, `run` exits non-zero and the summary shows which one.

## Design notes

- **Nothing is vendored.** Engines are cloned from their canonical repositories at install time, so you always run upstream's latest code with upstream's own license.
- **Dependency checks before action.** `status` reports exactly which tools (`git`, `node`, `npm`, `python`, `pip`) are missing before any install or run attempt.
- **Windows-safe process handling.** Batch-file shims (`npm.cmd` etc.) are invoked through raw command lines so paths containing spaces work.
- **Tested.** Unit tests cover dependency resolution, output streaming, exit-code propagation and the registry contract; CI runs them on Linux, Windows and macOS.

## License

MIT — see [LICENSE](LICENSE). The engines it launches remain under their own licenses (see each upstream repository).
