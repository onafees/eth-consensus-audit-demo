# Ethereum Consensus Layer Heartbeat Auditor

### Overview
This is a lightweight protocol-auditing tool written in **Go**, designed to validate the synchronization integrity of an Ethereum Beacon Node. It deterministically calculates the expected global slot based on the Ethereum Mainnet Genesis time and compares it against the live head reported by a Beacon Node API.

### Objective
In a distributed Proof-of-Stake system, clock synchronization and state consistency are critical for consensus. This tool serves as a primitive for **Consensus Layer (CL) Health Auditing**, specifically focusing on:
* **Temporal Accuracy:** Ensuring the node's local clock and sync logic align with the network's 12-second slot interval.
* **State Integrity:** Retrieving and surfacing the latest `state_root` for cross-node verification.



### Methodology

```
       EPOCH (6.4 Minutes)

|-------------------------------------------------------------------|
| Slot 0 | Slot 1 | Slot 2 |  ...  | Slot 30 | Slot 31 | (32 Total) |
|---|----|---|----|---|----|  ...  |----|----|----|----|------------|
    |        |        |                    |         |
 [Block]  [Empty]  [Block]              [Block]   [Block]
 (12s)    (12s)    (12s)                (12s)     (12s)
```

The auditor calculates the "Expected Slot" ($S_{expected}$) using the following formula:

$$S_{expected} = \lfloor \frac{T_{current} - T_{genesis}}{12} \rfloor$$

Where:
* **$T_{genesis}$** = 1606824023 (Dec 1, 2020, 12:00:23 UTC).
* The interval is a constant **12 seconds** per slot.

The script queries the `/eth/v1/beacon/headers/head` endpoint to retrieve the `actual_slot` and `state_root` currently recognized by the node.

### Example Output
Below are two successful audit runs demonstrating first a synchronization lag, and then shortly after, perfect synchronization between the calculated global clock and the Beacon Node:

```
PS C:\Users\OmiNafees\Downloads\eth-consensus-audit-demo> go run .\consensus_audit.go
--- Ethereum Consensus Layer Heartbeat Audit ---
Current Unix Time:    1772377166
Expected Global Slot: 13796095
Node Reported Slot:   13796094
Latest State Root:    0xfd523b0367917c7f47105b036a002cd662b51f26b5fad5cae7f89f632a0eeb17

--- Audit Summary ---
[WARN] Node slot mismatch. This could indicate latency or a synchronization lag.

PS C:\Users\OmiNafees\Downloads\eth-consensus-audit-demo> go run .\consensus_audit.go
--- Ethereum Consensus Layer Heartbeat Audit ---
Current Unix Time:    1772377174
Expected Global Slot: 13796095
Node Reported Slot:   13796095
Latest State Root:    0x79200ff470f3140efcaed75401ad45c8f10f1cccf53dd4aca670c3bd97efd38f

--- Audit Summary ---
[PASS] Node is perfectly synchronized with the global slot clock.
```

### Why This Matters for Protocol Testing
For a Protocol Tester, this script provides a foundation for detecting several critical failure modes:

1. Clock Drift: Identifying nodes that are consistently reporting slots behind the expected global time.
1. Synchronization Lag: Quantifying the delay between a block's propagation and a node's internal state update.
1. Inconsistent State Roots: When scaled to multiple endpoints, this logic can detect consensus deviations or "forking" between different client implementations (e.g., Prysm vs. Lighthouse).

### Usage
1. Ensure Go 1.21+ is installed.
1. Initialize the module: `go mod init eth-consensus-audit`
1. Execute the auditor: `go run consensus_audit.go`

### Future Work: Scaling to a Fuzzing Harness
* Multi-Client Consistency Check: Comparing `state_root` values across multiple RPC endpoints simultaneously.
* Latency Benchmarking: Measuring the time delta between slot-start and `head` availability.
* Integration Testing: Packaging this logic into a CI/CD pipeline for verifying local devnets during protocol upgrades.

### License
This project is licensed under the GPL v3 License.
