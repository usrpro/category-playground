-- name: bc-json

SELECT json_agg(crumbs) as result
FROM (
    WITH RECURSIVE
        q AS
        (
        SELECT  h.*, 1 AS level
        FROM    categories h
        WHERE   id = $1
        UNION ALL
        SELECT  hp.*, level + 1
        FROM    q
        JOIN    categories hp
        ON      hp.id = q.parent
        )
    SELECT id, parent, data
    FROM    q
    ORDER BY
            level DESC
) AS crumbs;

-- name: cat-tree

WITH    RECURSIVE
        q AS
        (
            SELECT  id, parent, lft, rgt, data, ARRAY[id] AS level
            FROM    categories hc
            WHERE   id = 5
            UNION ALL
            SELECT  hc.id, hc.parent, hc.lft, hc.rgt, hc.data, q.level || hc.id
            FROM    q
            JOIN    categories hc
            ON      hc.parent = q.id
            WHERE   array_upper(level, 1) < 3
        )
SELECT  id
FROM    q
ORDER BY
        level;

-- name: cat-tree-v2

with recursive categories_from_parents as
(
      select id, data, parent, '{}'::int[] as parents, 0 as level
        from categories
       where parent = $1  -- Offset

   union all

      select c.id, c.data, c.parent, parents || c.parent, level+1
        from      categories_from_parents p
             join categories c
               on c.parent = p.id
        where not c.id = any(parents) -- Loop protector
        and level < $2 -- Depth limit
)
select id, data, parent 
  from categories_from_parents;

-- name: recursive-cat

with recursive categories_from_parents as
(
      select id, data, parent, 0 as level
        from categories
       where parent = 0  -- Offset

   union all

      select c.id, c.data, c.parent, level+1
        from      categories_from_parents p
             join categories c
               on c.parent = p.id
       and level < 3 -- Depth limit
),
    categories_from_children as
(
        select c.parent,
            jsonb_build_object('id', c.id, 'data', c.data) as js
        from categories_from_parents tree
            join categories c using(id)
        where level > 0
        group by c.parent, c.id

    union all

        select c.parent,
            jsonb_build_object('id', c.id, 'data', c.data)
        ||  jsonb_build_object('children', jsonb_build_array(js)) as js
        from categories_from_children tree
            join categories c on c.id = tree.parent
)
select jsonb_pretty(js)
from categories_from_children
where parent = 0;