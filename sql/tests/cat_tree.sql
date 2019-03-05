INSERT
INTO    categories
WITH RECURSIVE
        ini AS
        (
        SELECT  5 AS level, 10 AS children
        ),
        range AS
        (
        SELECT  level, children
        FROM    ini
        ),
        q AS
        (
        SELECT  s AS id, 0 AS parent, level, children
        FROM    (
                SELECT  r.*, generate_series(1, children) s
                FROM    range r
                ) q2
        UNION ALL
        SELECT  id * children + position, id, level - 1, children
        FROM    (
                SELECT  generate_series(1, children) AS position, q.*
                FROM    q
                ) q2
        WHERE   level > 0
        )
SELECT  id, parent, 'Value' || id
FROM    q;
 
ANALYZE categories;