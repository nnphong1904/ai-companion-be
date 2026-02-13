# AI Companion Backend

Go REST API powering an AI companion social app with Stories, Relationship States, and Memories.

**Live API:** `<deployment-url>`
**Stack:** Go 1.24 &middot; Chi v5 &middot; PostgreSQL (Supabase) &middot; OpenAI GPT-4o-mini

---

## Table of Contents

- [Competitor Research](#competitor-research)
- [Product Feature Decisions](#product-feature-decisions)
- [Architecture & Technical Decisions](#architecture--technical-decisions)
- [Database Design & Scalability](#database-design--scalability)
- [API Reference](#api-reference)
- [Engineering Process & Challenges](#engineering-process--challenges)
- [Running Locally](#running-locally)

---

## Competitor Research

I studied the following products to understand what makes AI companionship compelling:

- **Nectar AI** — The primary reference. Stories feel native to the companion experience because they make AI characters feel like they have independent lives happening "off-screen." The social loop (stories -> reactions -> chat) keeps users returning. Nectar treats companions as social media presences, not just chatbots.
- **Character.ai** — Demonstrates that personality consistency is the #1 retention driver. Users return for the *character*, not the technology. This informed my decision to make personality traits directly influence AI response generation and mood progression.
- **Replika** — Pioneered the emotional bond system. Their "relationship levels" create a progression mechanic that gamifies engagement. I adapted this into my Relationship States feature, but made it bidirectional — the companion's mood toward the *user* changes, not just a level counter.
- **Instagram Stories** — The UX benchmark for Stories. Key patterns: tap-to-advance, hold-to-pause, swipe-to-next-person, progress bar segmentation, 24h expiry. These informed the backend's data model (stories with ordered media slides, expiration timestamps).

**Key insight:** The most successful companion apps create the illusion of a persistent, autonomous being. Stories are powerful because they imply the companion *did something* without the user prompting it.

---

## Product Feature Decisions

### Core Feature: Stories

Instagram-style ephemeral stories posted by AI companions. Each story contains one or more media slides (images/videos) that auto-expire. Users can react with emoji reactions, which deepens the relationship.

**Why Stories?** Stories are the heart of the "Social Loop" described in the assessment. They create passive engagement — users check in to see what their companions are up to, even without initiating a conversation. This mirrors real social media behavior and makes companions feel autonomous.

### Custom Feature 1: Relationship States (Emotional Bond System)

A dynamic emotional model where each user-companion pair has two independent scores:

- **Mood Score (0-100):** How the companion currently *feels* about the user. Affects AI response tone.
  - Distant (0-20) -> Neutral (20-50) -> Happy (50-80) -> Attached (80+)
- **Relationship Score (0-100):** The overall depth of the bond, built over time.

Scores increase through interactions — chatting (+2 mood, +1 relationship), reacting to stories (+3 mood, +2 relationship). The mood label is fed directly into the OpenAI system prompt, so a "Distant" companion gives cold, disengaged replies while an "Attached" companion is deeply caring and uses pet names.

**Why this feature?** After studying Replika and Character.ai, it became clear that static AI personalities get stale. Users need to feel like their interactions *matter* — that the companion remembers and evolves. Relationship States create a progression mechanic that rewards consistent engagement and makes every interaction feel consequential. It also solves a real product problem: giving users a reason to come back daily. Checking your companion's mood becomes habitual, like checking social media.

### Custom Feature 2: Memories (Curated Moment Timeline)

Users can save meaningful moments from their conversations as "Memories" — text snapshots with optional tags that form a shared timeline with each companion. Memories can be pinned to keep important moments at the top.

**Why this feature?** Companion apps generate enormous amounts of conversation, but most of it is ephemeral small talk. Memories let users curate the moments that matter — a funny exchange, a meaningful confession, a milestone in the relationship. This serves two purposes: (1) it gives users a sense of investment in the relationship ("look at everything we've shared"), and (2) it creates a nostalgia mechanism that deepens emotional attachment. The pinning system lets users build a personal highlight reel.

---

## Architecture & Technical Decisions

### Why Go? (First-Time Go Developer)

Go was a requirement for this assessment, and this was my first project in the language. I used AI tools extensively to learn Go idioms, best practices, and the standard library as I built. Rather than fighting the unfamiliarity, I treated it as a learning accelerator — studying each generated pattern, understanding *why* Go does things differently (explicit error handling, interfaces over inheritance, composition over hierarchy), and iterating until the code felt idiomatic rather than "translated from another language."

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

---

## Database Design & Scalability

### Schema Overview

8 tables with Row Level Security on all of them:

| Table | Purpose | Key Index Strategy |
|---|---|---|
| `users` | Authentication | Hash index on email for O(1) login lookup |
| `companions` | AI character profiles | Full table scan (5 rows, cached) |
| `stories` | Story metadata + expiry | `(companion_id, created_at DESC)` for per-companion feed; `(created_at DESC)` for global paginated feed |
| `story_media` | Ordered slides within stories | `(story_id, sort_order)` for batch loading |
| `story_reactions` | Emoji reactions (UPSERT) | `UNIQUE(user_id, media_id)` for atomic upsert |
| `messages` | Chat history | `(user_id, companion_id, created_at DESC)` for cursor pagination |
| `relationship_states` | Mood + relationship scores | `UNIQUE(user_id, companion_id)` for single-row lookup |
| `memories` | Curated moments | `(user_id, companion_id, pinned DESC, created_at DESC)` for pinned-first timeline |

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

## API Reference

### Authentication
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/auth/signup` | Register a new user |
| `POST` | `/api/auth/login` | Authenticate and receive JWT |
| `GET` | `/api/auth/me` | Get current user profile |

### Companions
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/companions` | List all companions |
| `GET` | `/api/companions/:id` | Get companion details |

### Stories
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/stories?cursor=&limit=` | Paginated active stories feed |
| `GET` | `/api/companions/:id/stories` | Stories for a specific companion |
| `POST` | `/api/stories/:id/react` | React to a story slide |

### Messages
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/companions/:id/messages?cursor=&limit=` | Paginated chat history |
| `POST` | `/api/companions/:id/messages` | Send message, receive AI reply |

### Relationships
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/relationships` | All user relationships |
| `GET` | `/api/companions/:id/relationship` | Relationship with specific companion |
| `POST` | `/api/onboarding/select-companion` | Select companion during onboarding |

### Memories
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/companions/:id/memories?limit=` | Paginated memory timeline |
| `POST` | `/api/companions/:id/memories` | Create a memory |
| `DELETE` | `/api/memories/:id` | Delete a memory |
| `PATCH` | `/api/memories/:id/pin` | Toggle pin status |

All protected endpoints require `Authorization: Bearer <jwt>`.

---

## Engineering Process & Challenges

### Process

1. **Research phase:** Studied Nectar AI, Character.ai, Replika, and Instagram to understand what makes companion apps compelling and what the "Social Loop" means in practice.
2. **Learning Go with AI:** Since this was my first Go project, I used AI tools to learn the language's patterns while building. I'd describe the architecture I wanted, study the generated code to understand Go-specific patterns (interface-based polymorphism, goroutine-safe connection pools, struct embedding), and refactor until the code felt idiomatic. Key learnings: Go's error handling philosophy, the `context.Context` propagation pattern, and how `pgx` handles connection pooling differently from ORMs I was familiar with.
3. **Schema-first design:** Designed the database schema before writing any Go code. Every table, index, and RLS policy was planned upfront to avoid retroactive migrations.
4. **Layer-by-layer implementation:** Built bottom-up — database -> repositories -> services -> handlers -> router. Each layer was tested in isolation before adding the next.
5. **Optimization pass:** After the feature set was complete, audited every query against its indexes using the Supabase performance advisor. Dropped 6 redundant indexes and added cursor-based pagination where queries were unbounded.

### Challenges & Debugging

**PgBouncer prepared statement conflicts:** Early in development, queries intermittently failed with "prepared statement does not exist" errors. This happened because Supabase's default pooler runs in transaction mode, where connections are shared between requests. pgx's default extended protocol caches prepared statements per-connection, but in transaction mode you might get a different connection on the next query. The fix was setting `QueryExecModeExec` to use simple protocol when the pooler is enabled.

**Relationship state in story reactions:** The `updateRelationshipOnReaction` function originally called `GetActiveStories()` — fetching *all* active stories just to find one story's `companion_id`. This was a hidden O(n) lookup that would degrade as stories accumulated. Fixed by adding a `GetByID` method to the story repository — a single PK lookup instead of a full table scan.

**Cursor pagination for multi-column sorts:** The memories query sorts by `(pinned DESC, created_at DESC)`, which makes cursor pagination non-trivial. A compound cursor encoding both `pinned` and `created_at` was considered but added significant complexity. Since memories per user-companion pair are naturally bounded (50-200 typical), I opted for a simple LIMIT-based approach with a `has_more` indicator — pragmatic over perfectly "correct."

**Migration idempotency:** All migrations use `IF NOT EXISTS` / `IF EXISTS` guards and `ON CONFLICT DO NOTHING` for seed data. This allows migrations to be re-run safely without failing on duplicate objects, which is important when the migration runner executes all files on every server start.

**Learning Go idioms on the fly:** Coming from other languages, my initial instinct was to reach for classes and inheritance. Go's composition model (struct embedding + interfaces) required a mindset shift. For example, the repository layer uses interfaces (`StoryRepository`, `MemoryRepository`) consumed by services — this is Go's version of dependency inversion, and it was unintuitive at first. AI tools helped me understand the *why* behind patterns like accepting interfaces and returning structs, which I now recognize as key to Go's simplicity.

---

## Running Locally

### Prerequisites
- Go 1.24+
- PostgreSQL (or a Supabase project)

### Setup

```bash
# Clone and install dependencies
git clone <repo-url>
cd ai-companion-be
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your database URL, JWT secret, and OpenAI key

# Run
make run
```

The server starts on `:8080` and automatically runs all migrations.

### Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `JWT_SECRET` | Yes | — | Secret for JWT signing |
| `OPENAI_KEY` | Yes | — | OpenAI API key |
| `OPENAI_MODEL` | No | `gpt-4o-mini` | Model for companion responses |
| `SERVER_PORT` | No | `8080` | HTTP server port |
| `DB_USE_POOLER` | No | `true` | Enable PgBouncer compatibility |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Frontend origin |
