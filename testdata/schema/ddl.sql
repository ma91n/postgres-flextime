DROP TABLE IF EXISTS notification;
CREATE TABLE notification
(
    notification_id VARCHAR(32)              NOT NULL,
    user_id         VARCHAR(20)              NOT NULL,
    registered_at   TIMESTAMPTZ              NOT NULL,
    title           VARCHAR(50)              NOT NULL,
    content         VARCHAR(1000),
    read_status_typ VARCHAR(1), -- 0:未読 1:既読 2:一括既読
    opened_at       TIMESTAMPTZ,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    revision        BIGINT                   NOT NULL,
    CONSTRAINT notification_pkc PRIMARY KEY (notification_id)
);

CREATE OR REPLACE FUNCTION flex_timestamp()
    RETURNS TIMESTAMPTZ AS
$$
BEGIN
    DECLARE
        flex_ts TIMESTAMPTZ := (SELECT fix_time
                                FROM flex_time
                                WHERE seq_num = (SELECT max(seq_num) FROM flex_time));
    BEGIN
        IF flex_ts IS NOT NULL
        THEN
            RETURN flex_ts;
        ELSE
            RETURN current_timestamp;
        END IF;
    END;
END;
$$ LANGUAGE PLPGSQL;
;
