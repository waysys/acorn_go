# Proofs of Week Functions

## Proof of WeekDay(date)

**Precondition**:
```
IsADate(date)
```

**Prove Postcondition**:
```
SUNDAY <= WeekDay(date) <= SATURDAY
and
WeekDay(date.Increment()) = daysOfWeek[(WeekDay(date) + 1) mod 7]
```
```
Theorem 1: SUNDAY <= WeekDay(date) <= SATURDAY

|    | WeekDay(date)                                |
| =  |    WeekDay function code                     |
|    | convertToAbsolute(date) mod 7                |
| => |    definition of mod                         |
|    | 0 <= WeekDay(date) <= 6                      |
| =  |    SUNDAY = 0 and SATURDAY = 6, substitution |
|    | SUNDAY <= WeekDay(date) <= SATURDAY          |  

Theorem 2: WeekDay(date.Increment()) = daysOfWeek[(WeekDay(date) + 1) mod 7]

|    | WeekDay(date.Increment())                    |
| =  |   convertToAbsolute(date.Increment() )       |
|    | (convertToAbsolute(date) + 1) mod 7          |
| =  |   definition of WeekDay() and daysOfWeek     |
|    | daysOfWeek[(WeekDay(date) + 1) mod 7]        |
  
Proof: Theoreem 1 and Theorem 2 => Postcondition
```