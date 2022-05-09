DO $$
DECLARE
    column_exist boolean := false;
BEGIN
SELECT count(*) != 0 INTO column_exist
    FROM information_schema.columns
    WHERE table_name = 'oauthapps'
    AND column_name = 'matterfossappid';
IF column_exist THEN
    UPDATE OAuthApps SET MatterfossAppID = '' WHERE MatterfossAppID IS NULL;
    ALTER TABLE OAuthApps ALTER COLUMN MatterfossAppID SET DEFAULT '';
    ALTER TABLE OAuthApps ALTER COLUMN MatterfossAppID SET NOT NULL;
END IF;
END $$;
