# Unit Testing Guide — Backend

## Overview

This document defines the unit testing standard for the API-5SEM backend.
As a DevOps practice, one example has been implemented. Developers must replicate this pattern for each new use case they implement.

---

## Rules

- Tests must be in the **same folder** as the code being tested
- Test file name: `<filename>_test.go`
- Every use case must have **at least 2 scenarios**: success + error
- Tests must follow the **AAA pattern**: Arrange → Act → Assert
- Tests must be submitted in the **same PR** as the feature
- If CI fails, the **PR author** is responsible for fixing it

---

## Pattern

### 1. Create the test file

Inside `internal/usecase/`, create `yourfile_test.go`.

### 2. Create a mock

The mock simulates the database without connecting to it:

```go
type mockProjetoRepo struct {
    projetos []entity.DimProjeto
    err      error
}

func (m *mockProjetoRepo) FindAll() ([]entity.DimProjeto, error) {
    return m.projetos, m.err
}
```

### 3. Write the tests

```go
// Success scenario
func TestXxxUseCase_GetAll_Success(t *testing.T) {
    // Arrange
    mock := &mockXxxRepo{data: fakeData, err: nil}
    uc := NewXxxUseCase(mock)

    // Act
    result, err := uc.GetAll()

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got: %v", err)
    }
}

// Error scenario
func TestXxxUseCase_GetAll_Error(t *testing.T) {
    // Arrange
    mock := &mockXxxRepo{data: nil, err: errors.New("db error")}
    uc := NewXxxUseCase(mock)

    // Act
    result, err := uc.GetAll()

    // Assert
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}
```

---

## Running tests

```bash
# Run all tests
go test ./...

# Run only use case tests
go test ./internal/usecase/

# Run with details
go test -v ./internal/usecase/
```

---

## Reference example

See [`internal/usecase/project_test.go`](./internal/usecase/project_test.go) for a complete working example.
All use cases in `internal/usecase/` are fully tested following this standard.

