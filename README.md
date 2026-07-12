# quell-cli

The command-line tool for [Quell](https://qubitlabs.magnobit.com) — a backend-agnostic quantum circuit language. Write a circuit once in Quell, then compile or run it against IBM Quantum, AWS Braket, Google Quantum Engine, IonQ, Rigetti, or Azure Quantum, using your own credentials for whichever provider you choose.

## Install

Download the binary for your platform from [Releases](https://github.com/magnobit/quell-cli/releases) and put it on your `PATH`.

> This repo builds against Magnobit's private Quell compiler/runtime, so `go build`/`go install` from a fresh clone won't work outside Magnobit's own CI — grab a prebuilt release binary instead. The CLI wrapper source here is public so you can read exactly what it does and how it talks to each provider before you hand it credentials.

## Usage

```
quell run <file.quell>                Run a circuit (local sim or a configured backend)
quell compile <file.quell>            Compile to OpenQASM, Qiskit, Cirq, or Braket
quell serve                           Start a local HTTP compile server
quell ask "<question>"                AI assistant (needs ANTHROPIC_API_KEY)
quell convert <file.py>               Convert Python/Qiskit code to Quell
```

### Running on real hardware

Credentials and per-backend parameters can come from `quell.config.yml`, from environment variables, or straight from the command line — a flag always wins.

```sh
quell run bell.quell --backend ibm --ibm-token $IBM_TOKEN --ibm-device ibm_brisbane

quell run bell.quell --backend azure \
  --azure-tenant-id $TID --azure-client-id $CID --azure-client-secret $SECRET \
  --azure-subscription-id $SUB --azure-resource-group $RG --azure-workspace $WS \
  --azure-target ionq.simulator
```

A parameter without a dedicated flag yet can still be sent through, without waiting on a new CLI release:

```sh
quell run bell.quell --backend ionq --ionq-api-key $KEY --set ionq.error_mitigation=true
```

Run `quell run --help` for the full per-backend flag list, or use a config file:

```yaml
# quell.config.yml
backend: ibm
ibm:
  token: ${IBM_QUANTUM_TOKEN}
  device: ibm_brisbane
```

## License

Apache 2.0 — see [LICENSE](LICENSE). This covers the CLI wrapper source in this repository; it does not grant rights to Magnobit's private Quell compiler/runtime that the released binaries link against.
