-- ============================================================================
-- Optimization: drop redundant indexes, add stories pagination index,
-- add expired-story cleanup function.
--
-- Redundant indexes slow every INSERT/UPDATE/DELETE for zero read benefit.
-- Ref: Supabase Postgres Best Practices §1.3 (composite index left-prefix rule)
-- ============================================================================

-- 1. idx_stories_companion_id is a left-prefix of idx_stories_companion_created
DROP INDEX IF EXISTS idx_stories_companion_id;

-- 2. idx_story_reactions_user_id is a left-prefix of UNIQUE(user_id, media_id)
DROP INDEX IF EXISTS idx_story_reactions_user_id;

-- 3. idx_users_email_hash is redundant with the UNIQUE btree on email
DROP INDEX IF EXISTS idx_users_email_hash;

-- 4. idx_messages_companion_id — no query filters by companion_id alone
DROP INDEX IF EXISTS idx_messages_companion_id;

-- 5. idx_memories_companion_id — no query filters by companion_id alone
DROP INDEX IF EXISTS idx_memories_companion_id;

-- 6. idx_relationship_states_companion_id — no query filters by companion_id alone
DROP INDEX IF EXISTS idx_relationship_states_companion_id;

-- ============================================================================
-- New index for GetActiveStories pagination:
--   WHERE expires_at > NOW() ORDER BY created_at DESC LIMIT N
-- A btree on (created_at DESC) lets Postgres walk the index in order and stop
-- after LIMIT rows, filtering expires_at on the fly.
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_stories_created_at ON stories (created_at DESC);

-- ============================================================================
-- Cleanup function: delete expired stories (and their media via CASCADE).
-- Call periodically via pg_cron, an edge function, or application cron.
--
--   SELECT cleanup_expired_stories();
--
-- With pg_cron (if enabled):
--   SELECT cron.schedule('cleanup-stories', '0 * * * *',
--     $$ SELECT cleanup_expired_stories() $$);
-- ============================================================================
CREATE OR REPLACE FUNCTION cleanup_expired_stories()
RETURNS integer
LANGUAGE sql
AS $$
    WITH deleted AS (
        DELETE FROM stories
        WHERE expires_at < NOW() - INTERVAL '1 hour'
        RETURNING id
    )
    SELECT count(*)::integer FROM deleted;
$$;
