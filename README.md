# AI Companion Backend

Go REST API powering an AI companion social app with Stories, Relationship States, Memories, and Insights.

**Stack:** Go 1.24 &middot; Chi v5 &middot; PostgreSQL (Supabase) &middot; Supabase Storage &middot; OpenAI GPT-4o-mini

---

## Architecture & Technical Decisions

### Why Go? (First-Time Go Developer)

Go was a requirement for this assessment, and this was my first project in the language. I used AI tools extensively to learn Go idioms, best practices, and the standard library as I built. Rather than fighting the unfamiliarity, I treated it as a learning accelerator — studying each generated pattern, understanding _why_ Go does things differently (explicit error handling, interfaces over inheritance, composition over hierarchy), and iterating until the code felt idiomatic rather than "translated from another language."

What I learned and why Go was the right choice for this domain:

- **Low latency:** Sub-millisecond overhead per request. Critical for a real-time social app where users expect instant story loads and chat responses.
- **Concurrency model:** Goroutines handle thousands of concurrent connections without the thread-per-request overhead of traditional servers. Important for a mass-audience social product.
- **Small memory footprint:** A compiled Go binary uses ~10MB of RAM at idle vs ~100MB+ for Node.js. Lower infrastructure costs at scale.
- **No runtime dependencies:** Single static binary deployment. The Docker image is 15MB (Alpine-based).
- **Explicit error handling:** Go forces you to handle every error at the call site. Coming from languages with try/catch, this felt verbose initially, but it produces more predictable, debuggable code — every failure path is visible.

### Why PostgreSQL (Supabase)?

- **Relational model fits the domain:** Users, companions, stories, reactions, and relationships have clear entity-relationship structures with foreign key constraints. A document store would require duplicating relationship data or sacrificing referential integrity.
- **Row Level Security (RLS):** Supabase's RLS provides database-enforced multi-tenancy. Even if the application layer has a bug, users can never see each other's messages, reactions, or relationship states. This is defense-in-depth — critical for a social app handling personal conversations.
- **Connection pooling built-in:** Supabase provides PgBouncer in transaction mode, which the backend is specifically configured to support (disabling prepared statement caching, using simple protocol mode).
- **UPSERT support:** Story reactions use `ON CONFLICT DO UPDATE` for atomic insert-or-update without race conditions.
- **Mature indexing:** Composite indexes, partial indexes, and hash indexes provide the query optimization flexibility needed for the various access patterns (cursor pagination, batch loading, multi-column sorts).

### Why Raw SQL + pgx Over an ORM or Supabase SDK?

I deliberately chose raw SQL with the `pgx` driver instead of an ORM (GORM) or the Supabase client SDK:

- **ORMs can't express the queries this project needs.** Cursor-based pagination (`WHERE created_at < $cursor`), atomic UPSERT (`ON CONFLICT DO UPDATE`), and batch loading (`WHERE story_id = ANY($1::uuid[])`) either aren't supported or require falling back to raw SQL, defeating the ORM's purpose.
- **ORMs hide N+1 problems.** I solved N+1 explicitly with batch loading. GORM's lazy association loading would reintroduce it silently.
- **PgBouncer compatibility.** GORM uses prepared statements by default, which break in Supabase's transaction-mode pooler. pgx lets me switch to simple protocol (`QueryExecModeExec`) for pooler mode while keeping extended protocol for direct connections.
- **No official Supabase Go SDK.** The community library routes queries through the PostgREST HTTP API, adding an extra network hop per query. A direct Postgres connection via pgx avoids this overhead entirely. The Supabase SDK is designed for client-side use (browsers/mobile), not server-to-database communication.
- **Repository interfaces provide the same abstraction.** Services depend on interfaces like `StoryRepository`, not on SQL. This gives the same decoupling benefit an ORM would, without the performance trade-offs.

### Clean Architecture (Layered Separation)

```
cmd/server/main.go          -- Entry point, dependency injection
internal/
  handler/                   -- HTTP request/response, input validation
  service/                   -- Business logic, orchestration
  repository/                -- Data access, SQL queries
  models/                    -- Domain objects, request/response types
  middleware/                 -- JWT auth, request context
  ai/                        -- OpenAI client wrapper
  config/                    -- Environment configuration
  database/                  -- Connection pool, migrations
  router/                    -- Route definitions
migrations/                  -- Idempotent SQL migrations
```

**Why this structure?**

- **Testability:** Each layer depends on interfaces, not implementations. Repositories implement interfaces that services consume, enabling mock-based unit testing without a database.
- **Separation of concerns:** Handlers never touch SQL. Services never touch `http.Request`. This prevents the "fat handler" antipattern where business logic leaks into HTTP handlers.
- **Dependency injection:** All wiring happens in `main.go`. Services receive repository interfaces through constructors. No global state, no service locator.

### OpenAI Integration

The chat system uses GPT-4o-mini with a carefully crafted system prompt that incorporates:

1. The companion's personality traits (from the database)
2. The current mood label derived from the relationship state
3. Character guidelines that prevent the AI from breaking character

A fallback response matrix handles OpenAI outages — pre-written mood-aware responses keyed by personality trait ensure the companion never goes silent. This is a resilience pattern: the external API dependency is non-blocking for the core UX.

### Supabase Storage for Media Assets

Companion avatars and story media (images, videos) are hosted on Supabase Storage in two public buckets:

- **`avatars`** — AI-generated companion profile images
- **`stories`** — Story slides organized by companion (`stories/{companion}/{filename}`)

Public read policies allow the frontend to load media directly via CDN URLs without authentication. This avoids proxying binary data through the Go backend and leverages Supabase's edge caching.

### User-Scoped Story Feed

The stories feed (`GET /api/stories`) only returns stories from companions the authenticated user has connected with. This is achieved by joining `stories` with `relationship_states` on `(companion_id, user_id)` at the database level — no application-side filtering needed. Users who haven't selected any companions see an empty feed rather than content from strangers.

### Message–Memory Linking

Memories can optionally reference the source message via `message_id`. This enables two features:

1. **`is_memorized` flag on messages** — The chat history query uses a correlated `EXISTS` subquery against the `memories` table to annotate each message with whether it has been saved as a memory. This avoids a separate API call and keeps the chat UI in sync.
2. **Traceability** — When a user saves a message as a memory, the FK link preserves the origin. The partial index `idx_memories_message_id WHERE message_id IS NOT NULL` keeps the `EXISTS` lookup fast without indexing the majority of rows where `message_id` is NULL.

---

## Database Design & Scalability

### Schema Overview

9 tables with Row Level Security on all of them:

| Table                 | Purpose                       | Key Index Strategy                                                                                                                          |
| --------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `users`               | Authentication                | Hash index on email for O(1) login lookup                                                                                                   |
| `companions`          | AI character profiles         | Full table scan (5 rows, cached)                                                                                                            |
| `stories`             | Story metadata + expiry       | `(companion_id, created_at DESC)` for per-companion feed; joined with `relationship_states` to scope to user's connected companions         |
| `story_media`         | Ordered slides within stories | `(story_id, sort_order)` for batch loading                                                                                                  |
| `story_reactions`     | Emoji reactions (UPSERT)      | `UNIQUE(user_id, media_id)` for atomic upsert                                                                                               |
| `messages`            | Chat history                  | `(user_id, companion_id, created_at DESC)` for cursor pagination                                                                            |
| `relationship_states` | Mood + relationship scores    | `UNIQUE(user_id, companion_id)` for single-row lookup                                                                                       |
| `memories`            | Curated moments               | `(user_id, companion_id, pinned DESC, created_at DESC)` for pinned-first timeline; partial index on `message_id` for `is_memorized` lookups |
| `mood_history`        | Daily mood snapshots          | `(user_id, companion_id, recorded_date)` for trend queries                                                                                  |

### Scalability Decisions

**1. Cursor-based pagination over OFFSET (Rule 6.3)**

Messages and stories use keyset/cursor pagination (`WHERE created_at < $cursor ORDER BY created_at DESC LIMIT $N`). Unlike OFFSET, this is O(1) regardless of page depth — page 10,000 is exactly as fast as page 1. This is critical for chat history (thousands of messages) and the global story feed.

**2. N+1 query elimination via batch loading (Rule 6.2)**

Story media is loaded in a single batch query using `WHERE story_id = ANY($1::uuid[])` instead of one query per story. Loading 20 stories with their media takes exactly 2 queries regardless of media count, not 21+.

**3. Composite indexes following the left-prefix rule (Rule 1.3)**

Every multi-column query has a matching composite index with equality columns first and range/sort columns last. For example, `idx_messages_conversation (user_id, companion_id, created_at DESC)` serves both the equality filter (`WHERE user_id = $1 AND companion_id = $2`) and the sort (`ORDER BY created_at DESC`) in a single index scan.

**4. Redundant index elimination**

After the initial schema, I audited every index against actual query patterns and dropped 6 redundant indexes:

- Single-column indexes that were already covered as left-prefixes of composite indexes (e.g., `idx_stories_companion_id` was redundant with `idx_stories_companion_created`)
- A hash index on email that was redundant with the UNIQUE btree constraint
- FK indexes on `companion_id` columns where no query ever filters by `companion_id` alone

Each redundant index slows every INSERT/UPDATE/DELETE for zero read benefit. Removing them directly improves write throughput.

**5. UPSERT for atomic operations (Rule 6.4)**

Story reactions use `INSERT ... ON CONFLICT (user_id, media_id) DO UPDATE SET reaction = $5` — a single atomic operation that handles both first-time reactions and reaction changes. No race conditions, no check-then-insert patterns.

**6. Connection pooling with PgBouncer compatibility (Rule 2.3/2.4)**

The database layer explicitly supports Supabase's transaction-mode PgBouncer:

- Disables pgx prepared statement caching (`QueryExecModeExec`)
- Uses simple protocol to avoid "prepared statement does not exist" errors
- Maintains a pool of 20 max / 2 min connections

This allows the backend to handle thousands of concurrent users through a small connection pool rather than exhausting Postgres connections.

**7. Expired story cleanup**

A `cleanup_expired_stories()` SQL function deletes stories that expired over an hour ago. Can be scheduled via pg_cron or called from an edge function. CASCADE on `story_media` ensures media is cleaned up automatically. Without this, expired stories would bloat the table indefinitely.

**8. RLS as defense-in-depth (Rule 3.2/3.3)**

Every table has Row Level Security policies. User-scoped tables (messages, reactions, memories, relationships) use `USING (user_id = (select current_setting('app.current_user_id', true))::uuid)` with the `(select ...)` wrapper for per-query caching (100x+ faster than calling the function per-row on large tables).

---

## Running Locally

### Prerequisites

- Go 1.24+ (or Docker)
- PostgreSQL (or a Supabase project)

### Setup

```bash
# Clone and configure
git clone <repo-url>
cd ai-companion-be
cp .env.example .env
# Edit .env with your database URL, JWT secret, and OpenAI key
```

### Option 1: Docker (Recommended)

```bash
docker compose up --build
```

This builds a multi-stage Alpine image (~15MB) and starts the server on `:8080` with automatic health checks. The container reads all config from your `.env` file.

### Option 2: Without Docker

```bash
go mod download
make run
```

The server starts on `:8080` and automatically runs all migrations.

### Environment Variables

| Variable               | Required | Default                 | Description                    |
| ---------------------- | -------- | ----------------------- | ------------------------------ |
| `DATABASE_URL`         | Yes      | —                       | PostgreSQL connection string   |
| `JWT_SECRET`           | Yes      | —                       | Secret for JWT signing         |
| `OPENAI_KEY`           | Yes      | —                       | OpenAI API key                 |
| `OPENAI_MODEL`         | No       | `gpt-4o-mini`           | Model for companion responses  |
| `SERVER_PORT`          | No       | `8080`                  | HTTP server port               |
| `DB_USE_POOLER`        | No       | `true`                  | Enable PgBouncer compatibility |
| `CORS_ALLOWED_ORIGINS` | No       | `http://localhost:3000` | Frontend origin                |
