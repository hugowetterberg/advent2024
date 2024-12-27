# Advent of code 2024

Usage:

``` shellsession
$ go run ./cmd/advent2024 -day 1 -solution 2 -use-sample
Similarity: 31

$ go run ./cmd/advent2024 -h
Usage of advent2024:
  -day int
    	day to run (default 1)
  -sample string
    	name of specific sample input to use
  -solution int
    	solution to run (default 1)
  -use-sample
    	use the sample input
```

## Use specific samples

Beginning at day 12 a more robust sample input selection process was needed. Specify `-use-sample` to use "sample-[solution number].txt" or "sample.txt" as input, or specify `-sample=name` to use "sample-[name].txt" as input.

For day 12 the sample collection looked like this:

``` shell
sample-abba.txt   ## -sample=abba
sample-eshape.txt ## -sample=eshape
sample-small.txt  ## -sample=small
sample.txt        ## -use-sample
sample-xoxo.txt   ## -sample=xoxo
```

## Solution flags

Sometimes solutions have solution specific flags that can be set, these flags are separated from the main flags by `--`. To see the documentation for solution flags pass `-- -help` after the standard flags:

``` shellsession
$ go run ./cmd/advent2024 -day 12 -solution 1 -- -help
Usage of day12:
  -debug-region string
    	output debug image for region
  -verbose
    	verbose region output
```

For day 12 I needed some help visualising the region borders for the second solution, like so:

``` shellsession
$ go run ./cmd/advent2024 -day 12 -solution 1 -sample=eshape -- -debug-region E
```

![Graphical representation of the E-shape](d12/debug-eshape.png)

# Advent of SQL

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

~~I don't know why the results of this query wasn't accepted by the website, I've checked the results with manual queries and it seems to be fine. Seeing a lot of noise on reddit about correct answes not being accepted, so I don't really know if I'm doing something wrong here :shrug:~~

The elf_id was supposed to be used as an ASC tie-breaker for both most and least experienced, I previously assumed DESC for least experience.

``` sql
SELECT DISTINCT ON (primary_skill)
    first_value(elf_id) OVER (
        PARTITION BY primary_skill ORDER BY years_experience DESC, elf_id
    ) AS elf_1_id,
    first_value(elf_id) OVER (
        PARTITION BY primary_skill ORDER BY years_experience, elf_id
    ) AS elf_2_id,
    primary_skill AS shared_skill
FROM workshop_elves
ORDER BY primary_skill;
```

### Day 8

Expensive query, there might be a smarter way to do this.

``` sql
WITH RECURSIVE employees AS (
    SELECT staff_id, staff_name, manager_id,
        CASE WHEN manager_id IS NULL THEN array[]::int[] ELSE array[manager_id] END AS employees
    FROM staff
    UNION
    SELECT s.staff_id, s.staff_name, s.manager_id, mr.employees||s.manager_id
    FROM staff AS s
    INNER JOIN employees AS mr ON mr.staff_id = s.manager_id
),
deduped AS (
    SELECT DISTINCT ON (staff_id)
        staff_id, staff_name,
        COALESCE(array_length(employees, 1), 0)+1 AS level,
        employees AS path
    FROM employees
    ORDER BY staff_id, array_length(employees, 1) DESC NULLS LAST
)
SELECT * FROM deduped ORDER BY level DESC;
```

### Day 9

``` sql
WITH avg_ex AS (
    SELECT r.reindeer_name, t.exercise_name, AVG(t.speed_record) avg_speed
    FROM Reindeers AS r
    INNER JOIN Training_Sessions AS t ON t.reindeer_id = r.reindeer_id
    WHERE r.reindeer_name != 'Rudolf'
    GROUP BY r.reindeer_name, t.exercise_name
),
top_result AS (
    SELECT reindeer_name, MAX(avg_speed) best_average
    FROM avg_ex
    GROUP BY reindeer_name
)
SELECT reindeer_name, to_char(best_average, '9999D00') AS best_time FROM top_result
ORDER BY best_average DESC
LIMIT 3;
```

### Day 10

~~This is probably not the best solution, I'll have to do some googling of pivot queries.~~

Nah, this was fine.

``` sql
WITH reckoning AS (
    SELECT date,
        SUM(case when drink_name='Eggnog' then quantity end) AS eggnog,
        SUM(case when drink_name='Hot Cocoa' then quantity end) AS hot_cococoa,
        SUM(case when drink_name='Peppermint Schnapps' then quantity end) AS peppermint_schnapps
    FROM Drinks GROUP BY date
)
SELECT * FROM reckoning
WHERE eggnog=198
    AND hot_cococoa=38
    AND peppermint_schnapps=298;
```
