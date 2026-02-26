# internal/ui

Shared vocabulary for the UI layer — three independent packages that keep
handlers and templates free from hard-coded strings, scattered permission
checks, and HTMX magic values.

## policy
Translates an authenticated `*identity.Identity` into small structs of `bool` flags. 
Templates receive these structs and conditionally render buttons, links, and forms: 
**no auth logic leaks into markup**.

```text
               Identity
                  │
         ┌─────── │ ───────┐──────────┐
         ▼        ▼        ▼          ▼
      BuildNav  BuildUser  BuildSpec  BuildAgent
         │      Detail     Detail     Detail
         ▼        │          │          │
       Nav {      ▼          ▼          ▼
        ShowUsers   UserDetail  SpecDetail  AgentDetail
        ShowTasks   { CanEdit    { CanEdit    { CanEditLabels }
        ShowAgents    CanDelete    CanDeploy
        CanAddUser    CanEditRoles CanDelete }
        CanAddSpec    IsSelf }
       }
```

`BuildNav` runs on every page request; the detail builders run when rendering entity-specific HTMX fragments.

## routepath
Every URL the system knows about lives here: page paths (`Page*`) and API endpoints (`Api*`). 
Path-builder functions append entity IDs:
```text
const                              var (path builders)
─────                              ────
PageUsers   = "/users"             PageUserInfoByID = func(id) → "/users/info/{id}"
ApiUsers    = "/api/v1/users"      ApiUserCrudOp    = func(id) → "/api/v1/users/{id}"
ApiUser     = "/api/v1/users/"     ApiUserEnable    = func(id) → "/api/v1/users/{id}/enable"
...                                ...
```

Handlers register routes with constants; templates build `hx-get` / `hx-post` URLs with builders. 
One place to change if a path ever moves.

## trigger
Two concerns in one tiny file:

**1. Mutation events** after a successful writing, the API handler calls `trigger.Set(w, trigger.SpecUpdate)` 
which sets the `HX-Trigger` response header. All HTMX fragments that listen for that event name auto-refresh.
```text
  API handler                        Browser (HTMX)
  ───────────                        ──────────────
  specUpdate(w,r)
    ├─ save to storage
    ├─ trigger.Set(w, SpecUpdate)  ──→  HX-Trigger: spec_update
    └─ 200 OK                           │
                                        ▼
                                    panels with
                                    hx-trigger="spec_update from:body"
                                    re-fetch their content
```

**2. Polling presets** named intervals (`SpecsRefresh = Every5s`,`UserSessionsRefresh = Every1m`) 
keep polling cadences consistent and adjustable from one place.

## Typical request flow
```text
  Browser
    │
    │  GET /specs/info/42
    ▼
  handler/ui.go
    ├─ id := auth(r)
    ├─ nav := policy.BuildNav(id)          ← policy
    └─ render page/taskspec/Detail(nav, "42")
         │
         │  HTMX auto-fires hx-get on load
         ▼
  handler/api.go
    ├─ id := auth(r)
    ├─ pol := policy.BuildSpecDetail(id)   ← policy
    ├─ spec := svc.Get("42")
    └─ render content/taskspec/Detail(spec, pol)
         │
         │  User clicks "Deploy"
         ▼
  handler/api.go  specDeploy(w, r)
    ├─ svc.Deploy("42")
    ├─ trigger.Set(w, trigger.SpecUpdate)  ← trigger
    └─ 204 No Content
         │
         ▼
  Browser refreshes spec panel via HX-Trigger
```
