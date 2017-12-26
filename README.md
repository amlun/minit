# minit
a mini tantan app server

## Build

TODO

## Environment Require

TODO

## Design

TODO

## Important!

Relationship

> When A like B, we add two records into db.

> That means A like B, and B is liked by A.

> They are in one transaction. // pg tx

```
id | owner_id | user_id |  state   | re_state
---+----------+---------+----------+---------
 1 |        1 |       5 | liked    |
 2 |        5 |       1 |          | liked
```