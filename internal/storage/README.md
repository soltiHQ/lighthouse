# internal/storage
Persistence contracts for control-plane domain entities.  
Backend-agnostic interfaces: domain and application layers depend only on these, never on concrete implementations.

## Package map
```text
storage/
├── storage.go      store interfaces + aggregate Storage
├── error.go        sentinel errors (ErrNotFound, ErrConflict …)
├── pagination.go   ListResult[T], ListOptions, limits
├── filter.go       backend-agnostic filter markers (AgentFilter, RolloutFilter …)
│
└── inmemory/
    ├── storage.go   Store — aggregates GenericStore instances, implements Storage
    ├── generic.go   GenericStore[T] — thread-safe CRUD for any domain.Entity[T]
    ├── filter.go    concrete filters with builder API (ByLabel, ByStatus, Query …)
    └── cursor.go    opaque base64 cursor encoding / decoding
```

## Store interfaces
```text
  Storage (aggregate)
  ├── AgentStore        Upsert / Get / List / Delete
  ├── UserStore         Upsert / Get / GetBySubject / List / Delete
  ├── CredentialStore   Upsert / Get / GetByUserAndAuth / ListByUser / Delete
  ├── VerifierStore     Upsert / Get / GetByCredential / Delete / DeleteByCredential
  ├── SessionStore      Create / Get / ListByUser / RotateRefresh / Revoke / Delete / DeleteByUser
  ├── RoleStore         Upsert / Get / GetMany / GetByName / List / Delete
  ├── SpecStore         Upsert / Get / List / Delete
  └── RolloutStore      Upsert / Get / List / Delete / DeleteBySpec
```
Every method documents sentinel errors it may return.

## Error model
```text
  Error               When                                       Retryable?
  ─────               ────                                       ──────────
  ErrNotFound         entity does not exist                      no
  ErrAlreadyExists    create conflict (duplicate ID)             no
  ErrConflict         concurrent modification / version mismatch no
  ErrInvalidArgument  bad input, malformed cursor, wrong filter  no
  ErrUnavailable      temporary backend failure                  yes
  ErrInternal         unexpected / invariant-violating failure   no
```
All errors are compatible with `errors.Is()`.

## Pagination
```text
  caller                          storage
  ──────                          ───────
  ListAgents(filter, ListOptions{   ──→  snapshot + sort by (UpdatedAt DESC, ID ASC)
      Limit:  50,                        slice [cursor … cursor+limit]
      Cursor: "…",                       encode last item as NextCursor
  })                               ◀──  ListResult{ Items, NextCursor }
        │
        ▼
  next page: ListOptions{ Cursor: result.NextCursor }
```

- **Ordering**: `(UpdatedAt DESC, ID ASC)` — deterministic, no gaps or duplicates
- **Cursor**: opaque base64 JSON token, backend-validated
- **Limits**: 1 .. 500 (default 100)

## Filter pattern
Filters are **backend-specific**. The storage package declares empty marker interfaces (`AgentFilter`, `RolloutFilter` …).  
Each backend provides concrete types with builder methods:
```text
  inmemory.NewAgentFilter().
      ByPlatform("linux").
      ByLabel("env", "prod").
      Query("web-")                ──→  implements storage.AgentFilter

  inmemory.NewRolloutFilter().
      BySpecID("spec-42").
      ByStatus(kind.SyncStatusPending)  ──→  implements storage.RolloutFilter
```
Passing a filter from a wrong backend returns `ErrInvalidArgument`.

## In-memory implementation

### GenericStore[T]
Thread-safe CRUD for any `domain.Entity[T]`:
```text
  GenericStore[T domain.Entity[T]]
  ├── mu   sync.RWMutex
  └── data map[string]T

  Create(entity)           insert, fail if exists
  Upsert(entity)           insert or replace
  Update(id, fn(T) T)      load clone → apply fn → store (under lock)
  Get(id)        → clone   retrieve by ID
  GetMany(ids)   → clones  batch get, preserve order
  List(pred, opts) → page  snapshot → sort → cursor → slice
  Delete(id)               remove by ID
```

- Entities are **cloned** on writing and on read — no shared mutable state
- Long scans check `ctx.Done()` every 1000 iterations
- `List` releases the read lock before sorting (snapshot isolation)
