-- ============================================================================
-- Mood history: daily mood snapshots for trend visualization.
-- One row per user+companion per calendar day (upsert on interaction).
-- ============================================================================

CREATE TABLE IF NOT EXISTS mood_history (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    companion_id  uuid NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
    recorded_date date NOT NULL DEFAULT CURRENT_DATE,
    mood_score    numeric(5,2) NOT NULL,
    UNIQUE(user_id, companion_id, recorded_date)
);

-- Query pattern: WHERE user_id = $1 AND companion_id = $2 AND recorded_date >= $3
-- ORDER BY recorded_date ASC
CREATE INDEX IF NOT EXISTS idx_mood_history_lookup
    ON mood_history (user_id, companion_id, recorded_date);

ALTER TABLE mood_history ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'mood_history' AND policyname = 'mood_history_own_access') THEN
        CREATE POLICY mood_history_own_access ON mood_history FOR ALL
            USING (user_id = (select current_setting('app.current_user_id', true))::uuid);
    END IF;
END $$;
