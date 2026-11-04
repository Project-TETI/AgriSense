BEGIN;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM refresh_tokens
        GROUP BY token_hash HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'Duplicate refresh token hashes exist; no data was removed. Resolve conflicts before migration.';
    END IF;
END
$$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash
    ON refresh_tokens (token_hash);

COMMIT;
