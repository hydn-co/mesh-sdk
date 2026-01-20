## Coding Guidelines

This document is a short, opinionated starting point for idiomatic, maintainable Go code in this repository. It focuses on practical rules we can enforce in PR review and CI.

1. Formatting & tooling

   - Use `go fmt ./...` to format code. Enforce via CI.
   - Run `go vet ./...` as part of PR checks.

2. Imports

   - Group imports: stdlib first, blank line, external, blank line, internal packages.
   - Avoid unused imports; CI should fail the build if any remain.

3. Package design & APIs

   - Keep packages small and cohesive. One responsibility per package.
   - Exported names must be documented with godoc-style comments (full sentence, start with the name).
   - Prefer small interfaces near the callers (interface acceptance principle).
   - Avoid global package-level state where possible; prefer dependency injection.

4. Package organization by feature

   - Organize packages by feature or domain, not by technical layer (avoid `factories`, `utils`, `helpers`).
   - Feature packages should contain all related code: types, business logic, persistence, and feature-specific factories.
   - Use sub-packages within a feature for internal organization when needed (e.g., `auth/nkeys`, `auth/creds`).
   - Shared utilities and cross-cutting concerns live in separate packages (e.g., `env`, `secrets`, `testkit`).
   - Package structure should reflect business domains visible to users of the SDK.
   
   Example structure:
   ```
   pkg/
   ├── auth/              # Authentication feature
   │   ├── credentials/   # Credential management
   │   ├── nkeys/        # NATS nkeys handling
   │   └── broker/       # Broker authentication
   ├── leasekit/         # Distributed leasing feature
   ├── messaging/        # Message processing feature
   │   ├── bus/         # Bus message processors
   │   └── stream/      # Stream processors
   ├── env/             # Shared: Environment helpers
   ├── secrets/         # Shared: Crypto/secret generation
   ├── meshctx/         # Shared: Mesh context
   └── testkit/         # Shared: Testing utilities
   ```

5. Error handling

   - Prefer returning typed errors. Wrap errors using `fmt.Errorf("%w")` or `errors.Join` (Go 1.20+) with context.
   - Do not use panic for ordinary error conditions. Reserve `panic` for internal unrecoverable programmer errors.
   - Use sentinel errors sparingly; prefer error types where callers need to inspect fields.
   - Document error behavior in function comments when necessary.

6. Logging

   - Use structured logging (slog or a chosen logger). Do not log secrets.
   - Log at appropriate levels: Debug for developer info, Info for high-level events, Warn for retryable problems, Error for terminal failures.
   - Keep logs localized; callers may decide whether to log or propagate errors.

7. Concurrency

   - Prefer channels or context for cancellation and timeouts. Always accept `context.Context` as the first parameter in exported functions that do I/O or long-running work.
   - Use `sync` primitives carefully; prefer channels or higher-level abstractions when expressive.
   - Avoid sharing mutable state without synchronization.

8. Testing

   - Follow repository test guidelines (Arrange/Act/Assert). Use `require` for setup and `assert` for postconditions.
   - Write table-driven tests for variants. Keep tests deterministic and fast.
   - Use test doubles/mocks where appropriate; centralize test helpers in `pkg/testkit`.
   - Run `go test ./... -cover` in CI and aim for meaningful coverage, not arbitrary percentages.

9. Security

   - Never commit secrets. Enforce secret scanning (gitleaks) in CI and pre-commit hooks.
   - Validate external inputs and avoid unsafe string concatenation in security-sensitive contexts (URLs, SQL, shell commands).
   - Prefer HTTPS for network calls to auth and token endpoints.

10. Dependencies & modules

   - Use Go modules. Keep `go.mod` tidy and run `go mod tidy` on dependency changes.
   - Pin semver-minor versions where feasible and review transitive updates.
   - Prefer small, well-maintained libraries; avoid adding heavy dependencies for small utilities.

11. CI & PRs

    - Every PR must run tests and linters in CI. Failing CI should block merging.
    - Include focused unit tests for new behavior and small integration tests for cross-package interactions.
    - Add short PR descriptions and mention any migration or compatibility notes.

12. Documentation

    - Document exported package APIs with godoc comments and provide a short README for each top-level package where helpful.
    - Keep `README.md`, `SECURITY.md`, and `docs/*` up to date with environment variables and deployment notes.

13. Backwards compatibility

    - Follow semantic versioning. For breaking changes, increment major version and document migration steps.

14. Repo hygiene
    - Keep the repository free of large binary files. Use .gitignore for transient files.
    - Use the `docs/` directory for policy and developer guidance.

Examples and enforcement

- Add CI steps for `gofmt`, `go vet`, `golangci-lint`, `gitleaks`, and unit tests.
- Add pre-commit hooks (optional) for local checks.

This is a first draft. If you want, I can:

- Convert key rules into actionable CI steps and add a `golangci-lint` config.
- Add pre-commit hooks and a `Makefile` or `mage` targets for dev convenience.
