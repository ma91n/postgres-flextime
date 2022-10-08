DROP TABLE IF EXISTS company
;
CREATE TABLE company
(
    company_cd   varchar(5)               NOT NULL,
    company_name varchar(256)             NOT NULL,
    founded_year integer                  NOT NULL,
    status       varchar(1)               NOT NULL default 0,
    created_at   timestamp with time zone NOT NULL,
    updated_at   timestamp with time zone NOT NULL,
    revision     integer                  NOT NULL,
    CONSTRAINT company_pkc PRIMARY KEY (company_cd)
)
;

CREATE TABLE flex_time
(
    seq_num  BIGSERIAL,
    fix_time TIMESTAMPTZ
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
