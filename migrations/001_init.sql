-- ============================================================================
-- AI Companion Backend — Supabase Migration (idempotent)
--
-- Best practices applied:
--   [4.1]  text over varchar, timestamptz, boolean, numeric
--   [4.2]  All FK columns indexed
--   [4.4]  UUID PKs via gen_random_uuid()
--   [4.5]  Lowercase snake_case identifiers
--   [1.3]  Composite indexes for multi-column queries
--   [3.2]  RLS enabled on ALL public tables
--   [3.3]  RLS policies use (select ...) wrapper for per-query caching
--   [6.3]  Cursor-based pagination via (created_at) indexes
--   [6.4]  UPSERT pattern for story reactions
-- ============================================================================

-- ===================
-- Users
-- ===================
CREATE TABLE IF NOT EXISTS users (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email       text UNIQUE NOT NULL,
    password    text NOT NULL,
    name        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email_hash ON users USING hash (email);
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'users' AND policyname = 'users_self_access') THEN
        CREATE POLICY users_self_access ON users FOR ALL
            USING (id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;

-- ===================
-- Companions
-- ===================
CREATE TABLE IF NOT EXISTS companions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL,
    description text NOT NULL DEFAULT '',
    avatar_url  text NOT NULL DEFAULT '',
    personality text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE companions ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'companions' AND policyname = 'companions_read_all') THEN
        CREATE POLICY companions_read_all ON companions FOR SELECT USING (true);
    END IF;
END $$;

-- ===================
-- Stories
-- ===================
CREATE TABLE IF NOT EXISTS stories (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    companion_id  uuid NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
    created_at    timestamptz NOT NULL DEFAULT now(),
    expires_at    timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_stories_companion_id ON stories (companion_id);
CREATE INDEX IF NOT EXISTS idx_stories_companion_created ON stories (companion_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_stories_expires_at ON stories (expires_at);

ALTER TABLE stories ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'stories' AND policyname = 'stories_read_all') THEN
        CREATE POLICY stories_read_all ON stories FOR SELECT USING (true);
    END IF;
END $$;

-- ===================
-- Story Media
-- ===================
CREATE TABLE IF NOT EXISTS story_media (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    story_id    uuid NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    media_url   text NOT NULL,
    media_type  text NOT NULL CHECK (media_type IN ('image', 'video')),
    duration    int NOT NULL DEFAULT 5,
    sort_order  int NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_story_media_story_id ON story_media (story_id, sort_order);

ALTER TABLE story_media ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'story_media' AND policyname = 'story_media_read_all') THEN
        CREATE POLICY story_media_read_all ON story_media FOR SELECT USING (true);
    END IF;
END $$;

-- ===================
-- Story Reactions
-- ===================
CREATE TABLE IF NOT EXISTS story_reactions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    story_id    uuid NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    media_id    uuid NOT NULL REFERENCES story_media(id) ON DELETE CASCADE,
    reaction    text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE(user_id, media_id)
);

CREATE INDEX IF NOT EXISTS idx_story_reactions_user_id ON story_reactions (user_id);
CREATE INDEX IF NOT EXISTS idx_story_reactions_story_id ON story_reactions (story_id);
CREATE INDEX IF NOT EXISTS idx_story_reactions_media_id ON story_reactions (media_id);

ALTER TABLE story_reactions ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'story_reactions' AND policyname = 'reactions_own_access') THEN
        CREATE POLICY reactions_own_access ON story_reactions FOR ALL
            USING (user_id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;

-- ===================
-- Messages
-- ===================
CREATE TABLE IF NOT EXISTS messages (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    companion_id  uuid NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
    content       text NOT NULL,
    role          text NOT NULL CHECK (role IN ('user', 'companion')),
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages (user_id, companion_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_companion_id ON messages (companion_id);

ALTER TABLE messages ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'messages' AND policyname = 'messages_own_access') THEN
        CREATE POLICY messages_own_access ON messages FOR ALL
            USING (user_id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;

-- ===================
-- Relationship State
-- ===================
CREATE TABLE IF NOT EXISTS relationship_states (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    companion_id        uuid NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
    mood_score          numeric(5,2) NOT NULL DEFAULT 50.0,
    relationship_score  numeric(5,2) NOT NULL DEFAULT 0.0,
    last_interaction    timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    UNIQUE(user_id, companion_id)
);

CREATE INDEX IF NOT EXISTS idx_relationship_states_companion_id ON relationship_states (companion_id);

ALTER TABLE relationship_states ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'relationship_states' AND policyname = 'relationships_own_access') THEN
        CREATE POLICY relationships_own_access ON relationship_states FOR ALL
            USING (user_id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;

-- ===================
-- Memories
-- ===================
CREATE TABLE IF NOT EXISTS memories (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    companion_id  uuid NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
    content       text NOT NULL,
    tag           text,
    pinned        boolean NOT NULL DEFAULT false,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_memories_user_companion ON memories (user_id, companion_id, pinned DESC, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_memories_companion_id ON memories (companion_id);

ALTER TABLE memories ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'memories' AND policyname = 'memories_own_access') THEN
        CREATE POLICY memories_own_access ON memories FOR ALL
            USING (user_id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;
