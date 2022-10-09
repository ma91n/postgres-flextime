WITH update_cnt AS (
    UPDATE notification
        SET read_status_typ = '2'
            , updated_at = $1
            , revision = revision + 1
        WHERE user_id = $2
            AND read_status_typ = '0' -- 未読
        RETURNING 1)
SELECT count(*) as cnt
FROM update_cnt
;
