SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'OAuthApps'
        AND table_schema = DATABASE()
        AND column_name = 'MatterfossAppID'
    ) > 0,
    'UPDATE OAuthApps SET MatterfossAppID = "" WHERE MatterfossAppID IS NULL;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'OAuthApps'
        AND table_schema = DATABASE()
        AND column_name = 'MatterfossAppID'
    ) > 0,
    'ALTER TABLE OAuthApps MODIFY MatterfossAppID varchar(32) NOT NULL DEFAULT "";',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
