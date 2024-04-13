# Proofs for Date Functions

## Proof of Add(date, num)

**Precondition**: 
```
IsADate(date)
```

**Prove Postcondition**:
```
err != nil or
convertToAbsoluteDate(resultDate) = convertToAbsoluteDate(date) + num
```
```
case 1  1 <= convertToAbsoluteDate(date)  + num <= MaxAbsoluteDate
|   | convertToAbsoluteDate(date) + num                         |
| = |   convertToDate postcondition                             |
|   | convertToAbsoluteDate(convertToDate(absolulteDate + num)) |
| = |   function code and declaration of resultDate             |
|   | convertToAbsoluteDate(resultDate)                         |

Case 2  convertToAbsoluteDate(date) + num < 1
| => |   code                    |
|    | err != nil                |
| => |  definition of or         |
|    | postcondition is true     |

case 3 convertToAbsoluteDate(date) + num > MaxAbsoluteDate
| => |     code                  |
|    | err != nil                |
| =  |     definition of or      |
|    | postcondition is true     |

```