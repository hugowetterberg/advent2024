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

https://adventofsql.com/

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

### Day 5

``` sql
WITH
comparison AS (
     SELECT
        production_date,
        toys_produced,
        LAG(toys_produced, 1) OVER(
                           ORDER BY production_date
        ) previous_day_production
     FROM toy_production
     ORDER BY production_date
)
SELECT production_date, toys_produced, previous_day_production,
       toys_produced-previous_day_production AS production_change,
       toys_produced::float/previous_day_production AS production_change_percentage
FROM comparison
ORDER BY toys_produced::float/previous_day_production DESC NULLS LAST;
```

### Day 6

``` sql
WITH average AS (
    SELECT AVG(price) AS price FROM gifts
)
SELECT c.name AS child_name, g.name AS gift_name, g.price AS gift_price
FROM children AS c
    INNER JOIN gifts AS g ON g.child_id = c.child_id
    INNER JOIN average ON g.price > average.price
ORDER BY g.price
LIMIT 10;
```

...or to get the exact answer (name of the kid with most expensive gift that costs more than average):

``` sql
WITH average AS (
    SELECT AVG(price) AS price FROM gifts
)
SELECT c.name
FROM children AS c
    INNER JOIN gifts AS g ON g.child_id = c.child_id
    INNER JOIN average ON g.price > average.price
ORDER BY g.price
LIMIT 1;
```

### Day 7

I don't know why the results of this query wasn't accepted by the website, I've checked the results with manual queries and it seems to be fine. Seeing a lot of noise on reddit about correct answes not being accepted, so I don't really know if I'm doing something wrong here :shrug:

``` sql
SELECT DISTINCT ON (primary_skill)
    first_value(elf_id) OVER (
        PARTITION BY primary_skill ORDER BY years_experience DESC, elf_id
    ) AS elf_1_id,
    first_value(elf_id) OVER (
        PARTITION BY primary_skill ORDER BY years_experience, elf_id DESC
    ) AS elf_2_id,
    primary_skill AS shared_skill
FROM workshop_elves
ORDER BY primary_skill;
```
