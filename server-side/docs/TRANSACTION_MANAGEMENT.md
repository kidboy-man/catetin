# Transaction Management

This document explains how to use the transaction management system in the Catetin backend application.

## Overview

The transaction system provides a clean abstraction for database transactions that maintains the hexagonal architecture principle where the service layer doesn't need to know about GORM-specific implementation details.

## Architecture

### Components

1. **TransactionManager Interface** (`internal/repository/transaction.go`)
   - Defines the contract for transaction management
   - Located in the repository layer to maintain architecture boundaries
   - Infrastructure-agnostic

2. **GORM Implementation** (`internal/infrastructure/database/postgresql/transaction_manager.go`)
   - Concrete implementation using GORM
   - Handles actual database transaction operations
   - Uses context.Context for transaction propagation

3. **Repository Support** (All `*_repository_impl.go` files)
   - All repositories use `GetDB(ctx, r.db)` helper function
   - Automatically detects if a transaction is active in context
   - Falls back to regular DB if no transaction present

## Usage Patterns

### 1. Simple Transaction (Recommended)

Use `WithTransaction` for most cases. It automatically handles commit/rollback:

```go
func (s *AuthService) Register(ctx context.Context, fullName, email, password string) (*RegisterResponse, error) {
    // ... validation code ...

    var user *domain.User

    // Wrap operations in transaction
    err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // Create user
        user = domain.NewUser(fullName, email)
        if err := s.userRepo.Create(txCtx, user); err != nil {
            return err // Automatic rollback on error
        }

        // Create user auth
        userAuth := &repository.UserAuth{
            ID:               uuid.New(),
            UserID:           user.ID,
            AuthProviderID:   provider.ID,
            CredentialID:     email,
            CredentialSecret: hashedPassword,
        }
        if err := s.userAuthRepo.Create(txCtx, userAuth); err != nil {
            return err // Automatic rollback on error
        }

        return nil // Automatic commit on success
    })

    if err != nil {
        return nil, err
    }

    // Continue with non-transactional operations...
}
```

**Key Points:**
- Return `error` from the callback to trigger rollback
- Return `nil` to trigger commit
- Any panic will also trigger rollback
- Use the `txCtx` context (not the original `ctx`) for all repository operations inside the transaction

### 2. Nested Transactions

The transaction manager automatically handles nested transactions:

```go
// Outer service method with transaction
err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    // Do some work
    if err := s.userRepo.Create(txCtx, user); err != nil {
        return err
    }

    // Call another service method that also uses transactions
    // This will reuse the existing transaction (not create a nested one)
    if err := s.anotherService.DoSomething(txCtx); err != nil {
        return err
    }

    return nil
})
```

**Behavior:**
- If already in a transaction, `WithTransaction` reuses it
- No nested/savepoint transactions are created
- All operations participate in the same transaction
- A single commit/rollback at the outermost level

### 3. Manual Transaction Control (Advanced)

For complex scenarios requiring manual control:

```go
func (s *SomeService) ComplexOperation(ctx context.Context) error {
    // Start transaction
    txCtx, err := s.txManager.BeginTransaction(ctx)
    if err != nil {
        return err
    }

    // Check if we're in a transaction
    if s.txManager.IsInTransaction(txCtx) {
        log.Println("Transaction active")
    }

    // Do work
    if err := s.userRepo.Create(txCtx, user); err != nil {
        _ = s.txManager.RollbackTransaction(txCtx)
        return err
    }

    // More work with conditional logic
    if someCondition {
        if err := s.txManager.CommitTransaction(txCtx); err != nil {
            return err
        }
    } else {
        if err := s.txManager.RollbackTransaction(txCtx); err != nil {
            return err
        }
    }

    return nil
}
```

**Warning:** Manual control is error-prone. Prefer `WithTransaction` unless you have a specific need.

## Repository Implementation

All repositories automatically support transactions through the `GetDB` helper:

```go
func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
    model := r.domainToModel(user)

    // Use GetDB to support transactions
    db := GetDB(ctx, r.db)

    if err := db.Create(model).Error; err != nil {
        return err
    }

    user.ID = model.ID
    return nil
}
```

**How GetDB Works:**
1. Checks if context contains an active transaction
2. If yes, returns the transaction DB
3. If no, returns regular DB with context
4. Repository code doesn't need to know which it's using

## Context Propagation

Transactions are propagated through `context.Context`:

```go
// Context helpers in internal/repository/transaction.go
func SetTransactionInContext(ctx context.Context, tx interface{}) context.Context
func GetTransactionFromContext(ctx context.Context) interface{}
```

**Important:**
- Always use the transaction context (`txCtx`) returned by `WithTransaction`
- Pass this context to all repository operations
- Don't mix transaction and non-transaction contexts

## Error Handling

### Within Transactions

```go
err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    if err := s.userRepo.Create(txCtx, user); err != nil {
        // Wrap error with context but still trigger rollback
        return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to create user", 500)
    }
    return nil
})

if err != nil {
    // Transaction already rolled back at this point
    return nil, err
}
```

### Transaction-Level Errors

```go
err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    // Your operations...
    return nil
})

if err != nil {
    // Check if it's a transaction error
    if errors.Is(err, gorm.ErrInvalidTransaction) {
        // Handle transaction error
    }
    return err
}
```

## Best Practices

### ✅ DO

1. **Use WithTransaction for simple cases**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       // operations
       return nil
   })
   ```

2. **Keep transactions short**
   - Only include database operations
   - Avoid external API calls, file I/O, long computations

3. **Use txCtx consistently**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       s.repo1.Create(txCtx, data1) // ✅ Use txCtx
       s.repo2.Create(txCtx, data2) // ✅ Use txCtx
       return nil
   })
   ```

4. **Return errors to trigger rollback**
   ```go
   return appErrors.Wrap(err, code, msg, status) // ✅ Rolls back
   ```

5. **Perform token generation outside transactions**
   ```go
   // Inside transaction: database operations only
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       return s.userRepo.Create(txCtx, user)
   })

   // Outside transaction: token generation, external calls
   accessToken, err := s.jwtManager.GenerateAccessToken(...)
   ```

### ❌ DON'T

1. **Don't use original context inside transaction**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       s.repo.Create(ctx, data) // ❌ Wrong! Use txCtx
       return nil
   })
   ```

2. **Don't ignore transaction context**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       // ❌ Don't ignore txCtx parameter
       _, _ = txCtx, ctx
       return nil
   })
   ```

3. **Don't perform slow operations in transactions**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       s.userRepo.Create(txCtx, user)

       // ❌ Don't make external API calls in transactions
       response, _ := http.Get("https://api.example.com/verify")

       // ❌ Don't generate tokens in transactions
       token, _ := s.jwtManager.GenerateAccessToken(...)

       return nil
   })
   ```

4. **Don't manually commit/rollback with WithTransaction**
   ```go
   err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
       s.userRepo.Create(txCtx, user)
       s.txManager.CommitTransaction(txCtx) // ❌ WithTransaction handles this
       return nil
   })
   ```

## Testing

### Mocking Transactions

For unit tests, mock the TransactionManager interface:

```go
type mockTransactionManager struct {
    mock.Mock
}

func (m *mockTransactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
    // Just execute the function without actual transaction
    return fn(ctx)
}

func (m *mockTransactionManager) IsInTransaction(ctx context.Context) bool {
    return false
}

// Use in tests
func TestSomeService(t *testing.T) {
    mockTxMgr := new(mockTransactionManager)
    service := NewSomeService(repo, mockTxMgr)

    // Test service methods
}
```

### Integration Tests

For integration tests with real database:

```go
func TestTransactionRollback(t *testing.T) {
    db := setupTestDB(t)
    txManager := postgresql.NewTransactionManager(db)

    err := txManager.WithTransaction(context.Background(), func(txCtx context.Context) error {
        // Create user
        user := domain.NewUser("Test", "test@example.com")
        userRepo.Create(txCtx, user)

        // Force rollback by returning error
        return errors.New("force rollback")
    })

    assert.Error(t, err)

    // Verify user was not created (transaction rolled back)
    _, err = userRepo.FindByEmail(context.Background(), "test@example.com")
    assert.Error(t, err)
}
```

## Performance Considerations

### Transaction Scope

Keep transaction scope as small as possible:

```go
// ✅ Good: Minimal transaction scope
func (s *Service) ProcessOrder(ctx context.Context, order *Order) error {
    // Non-transactional validation
    if err := s.validateOrder(order); err != nil {
        return err
    }

    // Transaction only for database writes
    err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        if err := s.orderRepo.Create(txCtx, order); err != nil {
            return err
        }
        if err := s.inventoryRepo.UpdateStock(txCtx, order.Items); err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return err
    }

    // Non-transactional notification
    s.notifyCustomer(order)
    return nil
}
```

```go
// ❌ Bad: Transaction spans entire operation
func (s *Service) ProcessOrder(ctx context.Context, order *Order) error {
    return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // Slow validation inside transaction
        if err := s.validateOrder(order); err != nil {
            return err
        }

        s.orderRepo.Create(txCtx, order)
        s.inventoryRepo.UpdateStock(txCtx, order.Items)

        // External call inside transaction - very slow!
        s.notifyCustomer(order)
        return nil
    })
}
```

### Connection Pooling

Transactions hold database connections. Keep them short to avoid exhausting the connection pool.

## Troubleshooting

### Issue: "Transaction has already been committed or rolled back"

**Cause:** Trying to use a transaction context after it's been closed.

**Solution:** Don't store txCtx for later use. Use it only within the callback.

### Issue: Changes not visible to other operations

**Cause:** Reading data outside the transaction that was written inside.

**Solution:** Move reads inside the transaction or commit before reading:

```go
var user *domain.User

// Write inside transaction
err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    user = domain.NewUser(...)
    return s.userRepo.Create(txCtx, user)
})

// Now you can read (transaction committed)
foundUser, _ := s.userRepo.FindByID(ctx, user.ID)
```

### Issue: Deadlocks

**Cause:** Two transactions waiting for each other's locks.

**Solution:**
- Always access tables in the same order
- Keep transactions short
- Use appropriate isolation levels

## Examples

### Example 1: User Registration with Auth

```go
func (s *AuthService) Register(ctx context.Context, fullName, email, password string) (*RegisterResponse, error) {
    // Validation (outside transaction)
    provider, err := s.authProviderRepo.FindByName(ctx, EmailPasswordProviderName)
    if err != nil {
        return nil, err
    }

    hashedPassword, err := s.passwordHasher.Hash(password)
    if err != nil {
        return nil, err
    }

    var user *domain.User

    // Atomic user + auth creation
    err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // Create user
        user = domain.NewUser(fullName, email)
        if err := s.userRepo.Create(txCtx, user); err != nil {
            return err
        }

        // Create auth credentials
        userAuth := &repository.UserAuth{
            ID:               uuid.New(),
            UserID:           user.ID,
            AuthProviderID:   provider.ID,
            CredentialID:     email,
            CredentialSecret: hashedPassword,
        }
        if err := s.userAuthRepo.Create(txCtx, userAuth); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Token generation (outside transaction)
    accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(user.ID, email, fullName)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &RegisterResponse{
        User:         user,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    expiresIn,
    }, nil
}
```

### Example 2: Order Processing with Multiple Tables

```go
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
    return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // Create order
        if err := s.orderRepo.Create(txCtx, order); err != nil {
            return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to create order", 500)
        }

        // Update inventory for each item
        for _, item := range order.Items {
            if err := s.inventoryRepo.DecrementStock(txCtx, item.ProductID, item.Quantity); err != nil {
                return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to update inventory", 500)
            }
        }

        // Create payment record
        payment := &domain.Payment{
            OrderID: order.ID,
            Amount:  order.Total,
            Status:  "pending",
        }
        if err := s.paymentRepo.Create(txCtx, payment); err != nil {
            return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to create payment", 500)
        }

        return nil
    })
}
```

## Summary

The transaction management system provides:
- ✅ Clean abstraction that respects hexagonal architecture
- ✅ Automatic commit/rollback with `WithTransaction`
- ✅ Nested transaction support (transaction reuse)
- ✅ Context-based transaction propagation
- ✅ Simple repository integration via `GetDB` helper
- ✅ Compatible with existing GORM code

For most use cases, simply wrap your multi-step operations in `WithTransaction` and use the transaction context (`txCtx`) for all repository calls.
