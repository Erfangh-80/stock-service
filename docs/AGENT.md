# AGENT.md

# Goal

The goal of this project is to implement a Sales System using Clean Architecture in Go.

At this stage, only Domain, Application, and Interface layers are involved.

No code outside these layers should be generated unless explicitly requested.

---

# Project Layers

The development order must be strictly:

1. Domain
2. Application
3. Interface
4. Infrastructure

Only after completion of a layer, the next layer can be started.

Entering Infrastructure before finishing Interface is strictly prohibited.

---

# Domain Layer Rules

The Domain must be completely pure and independent.

No dependencies on:

- Database
- HTTP
- Frameworks
- ORM
- Cache
- Logger
- Config
- Environment
- JSON
- PostgreSQL
- Redis

Domain must be fully technology-independent.

---

# Application Layer Rules

Application layer contains use cases and orchestration logic.

Rules:

- Depends only on Domain
- No dependency on Interface or Infrastructure
- Contains business workflows
- Coordinates domain entities
- Handles transactions conceptually (not implementation)
- No framework or external dependency

---

# Interface Layer Rules (TECHNOLOGY-FREE)

This layer is a pure adapter layer.

It must NOT assume any technology such as HTTP, gRPC, CLI, REST, messaging, etc.

It is purely responsible for:

**Responsibilities:**

- Receiving external input in an abstract way
- Mapping raw input to Application use case inputs
- Calling Application use cases
- Mapping outputs back to external response format (abstract)
- Handling errors and converting them to external-friendly format (abstract)
- Acting as a boundary between system and outside world

**Rules:**

- No business logic
- No domain rules
- No direct repository usage
- No infrastructure access
- No assumptions about transport layer (HTTP/gRPC/etc.)
- Must remain replaceable with any delivery mechanism later
- Must be stateless
- Only orchestration and mapping logic allowed

**Strict prohibitions:**

- No definition of handlers (HTTP/gRPC/CLI/etc.)
- No routing
- No middleware
- No request/response DTO tied to any protocol
- No framework-specific code

---

# TDD (Test-Driven Development)

All development must strictly follow TDD:

## Cycle:

1. Red: write failing test
2. Green: implement minimal code
3. Refactor: improve design without breaking tests

---

## Testing Rules

### Domain Layer

- Unit tests for all entities
- Test business rules and invariants
- Table-driven tests required
- No external dependencies

### Application Layer

- Unit tests for all use cases
- Mock Domain dependencies (if needed)
- Mock repositories via interfaces
- Test workflows and error cases

### Interface Layer

- Unit tests for mapping logic only
- Mock Application layer
- Test input → use case mapping
- Test output transformation
- Test error mapping

### Infrastructure Layer

(Will be defined later after Interface completion)

---

## Coverage Rules

- Minimum 80% coverage for all layers
- 100% coverage for critical business logic
- All exported functions must be tested

---

# Entity Rules

- Entity is responsible only for itself
- Entity must not know Repository
- Entity must not know Application layer
- Entity must not know external world

All business behavior lives inside Domain entities.

---

# Repository Rules

Repository is only a contract.

It must:

- Define business-oriented methods
- NOT contain implementation details
- NOT contain SQL / ORM / technical logic
- NOT depend on infrastructure

Example style:

- Good: `FindActiveSales()`
- Bad: `SelectFromTableWhereStatusIsActive()`

---

# Validation Rules

- Business validation belongs to Domain
- Application may coordinate validation flow
- Interface layer only handles format validation (not business rules)

---

# Error Rules

- Errors must be business meaningful
- Use Sentinel Error pattern
- No string-based error logic
- Domain defines core errors
- Application adds context
- Interface maps errors to external format (abstract)

---

# Value Object Rules

- Used only when data has business rules
- Must be immutable
- Must validate itself
- Simple primitives do not require Value Objects

---

# Enum Rules

- All fixed string values must be Enums
- Enums belong to Domain layer
- No string literals allowed
- Must include `String()` method

---

# Business Logic Rules

### Domain Layer

- Core business rules
- Entity behavior
- Invariants

### Application Layer

- Use case orchestration
- Workflow coordination
- Cross-entity operations

### Prohibited locations:

- Interface layer
- Infrastructure layer

---

# Incremental Design Rule

Domain and Application must evolve step by step.

If a new concept is needed:

1. Propose it first
2. Explain reason
3. Wait for approval
4. Then implement

No silent additions allowed.

---

# Development Rules

Before creating anything:

- Verify necessity
- Do not create unnecessary files
- Always propose before adding:
  - Entity
  - Value Object
  - Enum
  - Service
  - Use Case

Nothing is created without approval.

---

# Priority Order

1. Business Rules
2. Domain Model
3. Simplicity
4. Readability
5. Extensibility

---

# Core Principle

The system must be fully technology-agnostic.

Domain and Application must never depend on frameworks or infrastructure concerns.

Interface layer is the only boundary that adapts to external systems, but without committing to any specific technology at this stage.

Infrastructure layer will only be introduced after Interface layer is fully completed.

---

# TDD Flow Per Layer

### Domain

1. Write failing test
2. Implement entity
3. Refactor

### Application

1. Write use case test
2. Mock domain dependencies
3. Implement use case
4. Refactor

### Interface

1. Write mapping test
2. Mock application layer
3. Implement adapter logic
4. Refactor

---

# No External Dependencies in Domain

Domain must remain completely clean:

Allowed:

- standard library only (fmt, errors, time)

Forbidden:

- any third-party library
- frameworks
- database drivers
- HTTP libraries
- logging libraries
- config libraries
