---
title: useTimeDate
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useTimeDate`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useTimeDate:
    # no configuration options
```

## Details

Reports bad usage of `time.Date`. This rule detects several classes of problems:

- **Nil timezone**: Passing `nil` as the timezone argument causes a runtime panic.
- **Invalid dates**: Zero month or day arguments, out-of-bounds values for month (1-12), day (1-31), hour (0-23), minute (0-59), second (0-60), and nanosecond (0-999999999).
- **Calendar-aware validation**: Detects invalid dates like June 31st, February 29th in non-leap years, etc.
- **Swapped arguments**: Detects when the month and day arguments appear to be swapped.
- **Non-decimal notation**: Detects octal, hexadecimal, binary, exponential, and other non-decimal notations used as arguments, which are likely unintentional.
- **Sign issues**: Detects negative values (uncommon) and useless plus signs on arguments.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Nil timezone causes runtime panic
_ = time.Date(2023, 1, 2, 3, 4, 5, 0, nil)

// Leading zeros are octal notation and may not produce the intended value
_ = time.Date(2023, 01, 02, 03, 04, 05, 000000006, time.UTC)

// June has only 30 days
_ = time.Date(2023, 6, 31, 3, 4, 5, 0, time.UTC)

// February 29th in a non-leap year
_ = time.Date(2023, 2, 29, 3, 4, 5, 0, time.UTC)

// Month and day appear to be swapped
_ = time.Date(2023, 31, 6, 3, 4, 5, 0, time.UTC)

// Hexadecimal notation for year
_ = time.Date(0x7e7, 1, 2, 3, 4, 5, 6, time.UTC)
```

### Valid

```golang
// Standard usage with decimal literals
_ = time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)

// Using time.Month constants
_ = time.Date(2023, time.January, 2, 3, 4, 5, 1234567, time.UTC)

// Midnight is valid
_ = time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

// Leap year February 29th
_ = time.Date(2024, 2, 29, 3, 4, 5, 0, time.UTC)

// Negative year (BC) is valid
_ = time.Date(-500, 1, 2, 3, 4, 5, 6, time.UTC)
```
