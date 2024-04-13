# Proofs of Date Range Functions

## Proof of Overlap(dateRange1, dateRange2)

**Preconditions**
```
IsDateRange(dateRange1) and IsDateRange(dateRange2)
```

**Prove Postcondition**
```
If two ranges overlap, then there is a date that is in both ranges.  

overlap(range1, range2) => 
    exists date & inRange(range1, date) and inRange(range2, date)

if two ranges do not overlap, then no date is in both ranges.

not overlap(range1, range2) => 
    forall date & not inRange(range1, date) or  not inRange(range2, date)

Theorem 1: (overlap(range1, range2) = true) => 
    exists date & inRange(range1, date) and inRange(range2, date)

|    | assume overlap(range1, range2) = true                                          |
| => |   <Overlap code>                                                               |

|    | Case 1: assume range2.first >= range1.first and date = range2.first            |
| => |   <introduce exists; range2.first is in range2 by definitoon>                  |
|    | (exists date & range1.first <= date <= range1.last) and inRange(range2, date)  |
| => |   <definition of inRange>                                                      |
|    | exists date & inRange(range1, date) and inRange(range2, date)                  |

|    | Case 2: assume range2.first < range1.first and date = range1.first             |
| => |   <instroduce exists; range1.rirst is in range1 by defintion                   |
|    | (exists date & range2.first <= date <= range2.last ) and inRange(range1, date) |
| => |   exists date & inRange(range1, date) and inRange(range2, date)                |

| => |   <from Case 1 and Case 2                                                      |
|    | exists date & inRange(range1, date) and inRange(range2, date)                  |

Theorem 2: not overlap(range1, range2) => 
    forall date & not inRange(range1, date) or  not inRange(range2, date)

|    | assume overlap(range1, range2) = false                                         |
| => |    <Overlap code>                                                              |
|    | range1.last < range2.first or range1.first > range2.last                       |

|    | Case 1: range1.last < range2.first and InRange(range1,date)                    |                                           | 
| => |    <Overlap code, substitution>                                                |
|    | date <= range1.last < range2.first                                             |
| => |    <Definition of InRange>                                                     |
|    | date not inRange(range2)                                                       |

|    | Case 2: range1.first > range2.last and InRange(range1, date)                    |
| => |    <Overlap code, substitution>                                                |
|    | date >= range1.first > range2.last                                             |
| => |    <Definition of InRange                                                      |
|    | date not inRange)(range2)                                                      |

| => |    <from Case 1 and Case 2                                                     |
|    | forall date & not inRange(range1, date) or not InRange(range2, date)           |
