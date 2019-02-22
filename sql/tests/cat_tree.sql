INSERT
INTO    categories
WITH RECURSIVE
        ini AS
        (
        SELECT  8 AS level, 5 AS children
        ),
        range AS
        (
        SELECT  level, children,
                (
                SELECT  SUM(POW(children, n)::INTEGER * ((n < level)::INTEGER + 1))
                FROM    generate_series(level, 0, -1) n
                ) width
        FROM    ini
        ),
        q AS
        (
        SELECT  s AS id, 0 AS parent, level, children,
                1 + width * (s - 1) AS lft,
                1 + width * s - 1 AS rgt,
                width / children AS width
        FROM    (
                SELECT  r.*, generate_series(1, children) s
                FROM    range r
                ) q2
        UNION ALL
        SELECT  id * children + position, id, level - 1, children,
                1 + lft + width * (position - 1),
                1 + lft + width * position - 1,
                width / children
        FROM    (
                SELECT  generate_series(1, children) AS position, q.*
                FROM    q
                ) q2
        WHERE   level > 0
        )
SELECT  id, parent, lft, rgt, 'Value ' || id, RPAD('', 100, '*')
FROM    q;
 
ANALYZE categories;