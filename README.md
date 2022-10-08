# postgres-flextime


```sql
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
```

```sql
postgres=# SELECT flex_timestamp();
flex_timestamp        
------------------------------
 2022-10-08 22:50:28.52979+09
(1 row)

postgres=# SELECT flex_timestamp();
flex_timestamp         
-------------------------------
 2022-10-08 22:52:33.674613+09
(1 row)

-- Fix time
postgres=# INSERT INTO flex_time(fix_time) VALUES (TO_TIMESTAMP('2022-04-01 15:30:00', 'YYYY-MM-DD HH24:MI:SS'));
INSERT 0 1

    postgres=# SELECT flex_timestamp();
flex_timestamp     
------------------------
 2022-04-01 15:30:00+09
(1 row)

postgres=# SELECT flex_timestamp();
flex_timestamp     
------------------------
 2022-04-01 15:30:00+09
(1 row)
```


## Try

```
$ docker-compose up -d
$ docker exec -ti postgres-flextime_postgres_1 bash
#  psql -h localhost -p 5432 -U sample  -d postgres
postgres=# SELECT flex_timestamp();
```
