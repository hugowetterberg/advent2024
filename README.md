# Advent of code 2024

Usage:

``` shellsession
$ go run ./cmd/advent2024 -day 1 -solution 2 -use-sample
Similarity: 31

$ go run ./cmd/advent2024 -h
Usage of /tmp/go-build1045924321/b001/exe/advent2024:
  -day int
    	day to run (default 1)
  -solution int
    	solution to run (default 1)
  -use-sample
    	use the sample input
```

## Advent of SQL

Nerd sniped myself (with help from colleagues) into doing this as well. We'll see how long I'll stick to doing both ...or any :)

### Day 1

``` sql
WITH
complexity AS (
  SELECT
        unnest('{1,2,3}'::int[]) AS id,
        unnest('{"Simple Gift", "Moderate Gift", "Complex Gift"}'::text[]) AS description
),
category_workshop AS (
  SELECT
        unnest('{"outdoor", "educational"}'::text[]) AS category,
        unnest('{"Outside Workshop", "Learning Workshop"}'::text[]) AS name
)
SELECT
        c.name,
        wl.wishes->>'first_choice' AS primary_wish,
        wl.wishes->>'second_choice' AS backup_wish,
        wl.wishes->'colors'->>0 AS favorite_color,
        json_array_length(wl.wishes->'colors') AS color_count,
        cx.description AS gift_complexity,
        COALESCE(cw.name, 'General Workshop') AS workshop_assignment
FROM children AS c
        INNER JOIN wish_lists AS wl ON wl.child_id = c.child_id
        INNER JOIN toy_catalogue AS toy ON toy.toy_name = wl.wishes->>'first_choice'
        INNER JOIN complexity AS cx ON
              cx.id = CASE WHEN toy.difficulty_to_make <= 3 THEN toy.difficulty_to_make ELSE 3 END
        LEFT OUTER JOIN category_workshop AS cw ON cw.category = toy.category
ORDER BY c.name
LIMIT 5;
```

### Day 2

``` sql
WITH merged AS (
     SELECT * FROM (
            SELECT * FROM letters_a
            UNION
            SELECT * FROM letters_b
     ) ORDER BY id
)
SELECT string_agg(CHR(value), '') FROM merged
WHERE
        (value >= ASCII('a') AND value <= ASCII('z'))
        OR (value >= ASCII('A') AND value <= ASCII('Z'))
        OR (value >= ASCII('0') AND value <= ASCII('9'))
        OR (position(CHR(value) in ' !"''(),-.:;?')>0);
```

### Day 3

XML, nah, can't be bothered.

### Day 4

``` sql
WITH tag_cmp AS (
     SELECT toy_id,
            ARRAY(SELECT unnest(new_tags) EXCEPT SELECT unnest(previous_tags)) AS added_tags,
            ARRAY(SELECT unnest(previous_tags) EXCEPT SELECT unnest(new_tags)) AS removed_tags,
            ARRAY(SELECT unnest(previous_tags) INTERSECT SELECT unnest(new_tags)) AS unchanged_tags
     FROM toy_production
)
SELECT toy_id,
       COALESCE(array_length(added_tags, 1), 0) AS added_len,
       COALESCE(array_length(unchanged_tags, 1), 0) AS unchanged_len,
       COALESCE(array_length(removed_tags, 1), 0) AS removed_len
FROM tag_cmp
ORDER BY array_length(added_tags, 1) DESC NULLS LAST
LIMIT 1;
```
