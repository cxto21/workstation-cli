# Gentle AI Integration Plan

## Goal

Integrate Workstation with gentle-ai as the default agent ecosystem bootstrapper, while keeping Workstation as control plane and preserving local-first runtime behavior.

## Context

- Workstation domain model already defines Gentle AI as the agent execution substrate.
- gentle-ai is an ecosystem configurator, not a single-agent runtime.
- gentle-ai already supports VS Code Copilot and multiple agent adapters, plus MCP and memory-oriented setup patterns.

## Integration Contract

Workstation should call gentle-ai for ecosystem setup and sync, and keep Workstation responsible for:

- Operator UX and workflow orchestration
- Policy and approvals
- Local runtime state and process lifecycle
- Logs and audit trail

gentle-ai should handle:

- Agent config injection
- MCP server wiring
- Skills and persona provisioning
- Idempotent re-apply operations

## MVP Scope

### 1. Tool Detection and Status

Expose a status surface in Workstation that reports:

- gentle-ai binary availability in PATH
- detected version
- basic health state (callable / not callable)

### 2. Bootstrap Command

Add a safe bootstrap operation from Workstation that runs gentle-ai install in a controlled profile:

- first use dry-run
- then apply with explicit user confirmation
- default preset should be minimal or ecosystem-only

### 3. Sync Command

Expose a sync operation that runs gentle-ai sync from Workstation and writes structured logs to the Workstation log stream.

### 4. Audit Logging

Every Workstation invocation of gentle-ai should log:

- command name and arguments
- start and end timestamps
- exit code
- sanitized stderr/stdout summary

## Phase 2 Scope

- agent profile mapping: Workstation templates to gentle-ai presets
- per-project setup command wiring
- optional onboarding step that offers gentle-ai bootstrap during first run

## Phase 3 Scope

- policy gates for allowed components and MCP servers
- explicit approval flow for config mutations
- rollback entrypoints surfaced in Workstation UI

## Proposed First Implementation Slice

1. Add Workstation CLI commands for tool status and gentle-ai bootstrap/sync wrapper.
2. Keep shell invocation minimal and explicit (no background hidden mutation).
3. Capture output into Workstation logs.
4. Add tests for command success/failure mapping.

## Risks

- Command/flag drift between gentle-ai versions
- Partial setup if dependencies are missing
- Different MCP schema conventions across agents

## Mitigations

- run dry-run first and parse result
- strict version check and friendly upgrade guidance
- keep wrappers idempotent and retry-safe
- maintain compatibility matrix in docs

## Definition of Done for MVP

- Workstation can detect gentle-ai and show status
- Workstation can run a dry-run bootstrap and an apply bootstrap
- Workstation can run sync and show structured result
- operator-visible logs include command, duration, and exit code
- docs explain behavior and limitations