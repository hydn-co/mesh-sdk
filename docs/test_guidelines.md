## 🧪 Testing Guidelines

- All code must be covered with meaningful tests.
- Use Go’s standard `testing` package with `testify/assert` for assertions and `testify/require` to short circuit when arrange sections fail.
- Each test should only test one thing.
- Use table driven tests when a test behavior is operating on variations of the inputs
- Write tests in **behavioral style**, using names like:

  ```go
  func TestShouldReturnErrorWhenUserIsInvalid(t *testing.T)
  func TestShouldStoreResultGivenValidInput(t *testing.T)
  ```

 - Inside tests, follow this structure with clear comments:

  ```go
  // Arrange
  require.NoError(t, err)
  ...

  // Act
  ...

  // Assert
  assert.Equal(t, expected, actual)

  ```

Require vs. Assert semantics
----------------------------

Use `require` for setup (Arrange) validations and preconditions that must be true for the rest of the test to run. `require` fails the test immediately on failure and prevents a cascade of nil-pointer or misleading errors.

Use `assert` for post-condition checks in the Assert section so multiple checks can run and provide aggregated failure output. `assert` lets you make multiple assertions about the expected state after the Act step.

Example:

```go
func TestShouldCreateThingGivenValidInput(t *testing.T) {
    // Arrange
    db, err := OpenTestDB()
    require.NoError(t, err, "failed to open test DB") // stop if DB cannot be opened

    repo := NewRepo(db)

    // Act
    got, err := repo.Create(context.Background(), input)

    // Assert
    assert.NoError(t, err)          // check error but continue to check result
    assert.NotNil(t, got)          // multiple asserts allowed
    assert.Equal(t, expected.Name, got.Name)
}
```

Rationale:
- `require` protects the test from invalid setup and makes failures easier to diagnose.
- `assert` improves test signal by running multiple checks and reporting all differences in a single run.

When to prefer `require` in Assert:
- If a postcondition is a hard precondition for later assertions (e.g., you expect a non-nil value before inspecting fields), use `require` to avoid panics and to focus the failure.

When to prefer `assert` in Arrange:
- Avoid using `assert` in Arrange; if a setup step can continue without failing the test, consider restructuring the test so that `require` is used when necessary.

