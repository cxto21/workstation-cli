# AGENTS.md - Workstation Development Guide

This file defines how we build Workstation as a local-first, agentized environment for teams.

Workstation is the control plane.
Mato is one runtime resource inside Workstation.

## 1. Product Scope

### Vision
Build an Agentized Business Environment where:
- Humans supervise
- Agents execute
- Decisions and learnings are captured as organizational wisdom

### Platform Components
- Workstation (UI/CLI): operator surface and control plane
- Gentle AI: agent orchestration and tool execution backend
- Engram: long-term memory and retrieval layer
- Mato: terminal workspace runtime integrated into Workstation

## 2. Core Domain Model

### Org
Container for access, policies, and shared resources.
Start simple: low-friction team membership.

### Member
Human identity authenticated into an Org.
A Member can operate from multiple Seats.

### Seat
Execution context on an OS user environment.
Seat is not identity.
Seat is where workloads run.

### SSOT
Single Source of Truth for shared, declarative state.
Initially backed by GitHub repository sync.

### KBs and KB
- KBs: the full set of knowledge bases for an Org or Project
- KB: one concrete knowledge base directory, markdown-first, Obsidian-compatible

### Project
Collaborative unit of work tracked by the team.
Project is business/domain scope, not UI-only grouping.

### Agent
Automated worker process with profile, tools, memory access, and execution lifecycle.

### Subagent
Specialized or delegated agent under another agent/task context.

### Logs
Structured audit and observability events across core entities.

### Workflow
Reusable process definition, human-readable and agent-executable.

### Memory
Persistent organizational and agent memory, curated from outcomes and decisions.

## 3. Architectural Boundaries

To avoid complexity creep, keep this hard boundary:

### Shared Declarative State (SSOT)
Store in Git-backed SSOT:
- org/project metadata
- KB structure and markdown content
- workflow definitions
- agent profiles and capabilities metadata
- decisions, ADR-like records, task summaries

### Local Runtime State
Do not sync as SSOT truth:
- live PTY buffers and terminal cursor state
- transient daemon/session internals
- temporary UI state and local process handles

## 4. Integration Contracts

### Engram Integration (Memory Layer)
- Workstation consumes Engram through MCP interface
- KBs can map to Engram-backed indexed stores
- Retrieval feeds agent context windows
- Optional embedded terminal tab can run Engram TUI/bin for memory ops

### Gentle AI Integration (Agent Layer)
- Workstation dispatches execution requests to Gentle AI runtime
- Agent profiles map to capability/tool bundles
- Access policy eventually maps Member/Seat to allowed profiles
- Agent run metadata and outcomes are emitted to Logs and Memory

### Workstation Integration Role (Control Plane)
- Unified UX over Mato, Engram, Gentle AI
- Access boundary and policy point
- Sync orchestration between local state and SSOT
- Audit surface for requests, retrievals, tool actions, and decisions

## 5. Security Baseline

Initial baseline:
- Org-scoped membership
- Seat-aware authentication context
- explicit sync commands (no silent background writes)
- append-only critical logs where feasible

Future baseline:
- RBAC by Org/Project/KB/Agent capability
- mTLS between control plane and execution services
- per-tool policy gates and approvals

## 6. Local-First Delivery Plan

### Phase A - Workstation Shell
- Rebrand and UI structure alignment
- Chat, Alerts, Logs, Docs views as first-class surfaces
- Keep runtime stable

### Phase B - Docs and KB Operations
- local markdown KB indexing
- read-only previews in Workstation
- controlled write workflows

### Phase C - Agent Presence and Orchestration
- connect Workstation actions to Gentle AI runtime
- display capability metadata
- log agent execution lifecycle

### Phase D - Memory Integration
- connect Engram retrieval via MCP
- persist approved memory outcomes
- expose memory lookup and provenance in UI

### Phase E - SSOT Sync
- explicit pull/push sync with conflict policy
- multi-member collaboration on shared declarative artifacts
- cross-device bootstrap from repo

## 7. Engineering Rules

### Incremental Delivery
- Prefer small, reversible changes
- Keep runtime behavior stable while adding new layers
- Avoid broad refactors unless required by model boundaries

### Validation Discipline
- Build after each substantial change
- Run focused tests around touched modules
- Never claim integration until wired and observable

### Docs-First for Domain Changes
When changing domain concepts (Org, Seat, Project, KB, Memory):
- update this file
- update corresponding docs under docs/
- add a short changelog entry for rationale

## 8. Collaboration Workflow (Human + Agent)

### Preferred Cycle
1. Clarify intent and boundaries
2. Create a branch from `main` using `feat/*`, `fix/*`, or `ci/*`
3. Implement the minimal slice
4. Validate build/tests
5. Commit and push the branch
6. Open a pull request to `main` and request review
7. Document behavior and trade-offs
8. Iterate

### Contribution Flow
- Do not push directly to `main`
- Use a dedicated branch per change
- Keep PRs small and reviewable
- Prefer the repository helper script for opening PRs when available
- Merge only after another team member approves

### Hotfix Flow
- Hotfixes should use a `fix/*` branch and a PR to `main`
- If a fast-track emergency lands directly on `main`, the release workflow will tag it automatically from `Cargo.toml`
- A version bump is still required before that release can publish a new tag

### Session Changelog Rule
For each AI agent chat session:
- maintain exactly one session changelog file in docs/changelog/
- keep appending to the same file during the session
- only split if explicitly requested

### Agent Branching Rule
- Agents should work on `feat/*`, `fix/*`, or `ci/*` branches only
- Agents should push the branch and open a PR to `main`
- Agents should not skip review or merge directly to `main`

## 9. Repository Conventions

- Keep docs/changelog/ for dated development records
- Keep docs/todos/ for active execution plans
- Keep docs/release-notes/ for versioned release communication
- Keep AGENTS.md as the source for domain and process principles

## 10. Definition of Done for Workstation Features

A feature is done when:
- behavior is implemented in code
- UX is coherent with Workstation model
- build/test validation passes for touched scope
- docs are updated for user and developer understanding
- no contradiction with domain boundaries in this guide

## 11. Current Positioning Statement

Workstation is the team operating environment.
Gentle AI is the agent execution substrate.
Engram is the persistent memory substrate.
Mato remains a key terminal runtime resource within Workstation.
