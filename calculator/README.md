# Go Calculator

by Aryan

```
Expr -> Term | Expr + Term | Expr - Term
Term -> Factor | Term \* Factor | Term / Factor
Factor -> Number | ( Expr )

a \* (b + c)

Lexer -> get all Tokens -> numbers, operands, parentheses
Parser -> Generate Abstract Symbol Tree from Tokens
Calculation
```

## My scribbles:

- Create lexer []

2 \* (1 + 3)

{\*, {2, (1 + 3)}}

{\*, 2, {+, {1, 3}}}

{operand, left, right}
{operand, expression, expression}

```
2 * 1 + 3

{ +, (2 * 1), 3}
{ +, {*, {2, 1}}, 3}

```

```
[ 2, *, (, (, 1, +, 3, ), *, 4, ) - 2 ] 1 + 1 - 1 - 1

{ -, 2 * ((1 + 3) * 4), 2 } addsub
{ -, { *, 2, ((1 + 3) * 4) }, 2 } muldiv

```

```
[ (, (, 1, +, 3, ), *, 4, ) - 2 ] 1 + 1 - 1 - 1

{ -, ((1 + 3) * 4), 2 } addsub
{ -, ((1 + 3) * 4), 2 } muldiv
{ -, (1 + 3) * 4, 2 } paren
{ -, (1 + 3) * 4, 2 } addsub + muldiv

```

```
2 * 2 + 3 + 4

{ +, (2 * 2), (3 + 4)}
{ +, {*, 2, 2}, {+, 3, 4}}
```

```
2 + 3 * 2 * 2

{ +, {2, 3 * 2 * 2}}

{ +, {2, {*, 3, 2 * 2}}}

```

```
i = 0: "2"; i++
i = 1 (found *):
  create: Node(*, "2", x)

x -> i++:
i = 2: "1"; i++
i = 3 (found +):
  create: Node(+, "1", x)

x -> i++:
i = 4: "3"
```

```
2 * (1 + 3)

{2 * }
```
