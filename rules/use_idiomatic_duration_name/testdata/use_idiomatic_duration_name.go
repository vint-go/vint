package fixtures

import "time"

// Invalid: name implies seconds, but Duration has its own unit.
var timeoutSecs time.Duration = 5 * time.Second // MATCH /var timeoutSecs is of type time.Duration; don't use unit-specific suffix "Secs"/

// Invalid: name implies milliseconds.
var delayMs time.Duration = 100 * time.Millisecond // MATCH /var delayMs is of type time.Duration; don't use unit-specific suffix "Ms"/

// Invalid: name implies milliseconds with Milli suffix.
var retryMilli time.Duration = 200 * time.Millisecond // MATCH /var retryMilli is of type time.Duration; don't use unit-specific suffix "Milli"/

// Invalid: name implies minutes.
var waitMin time.Duration = 2 * time.Minute // MATCH /var waitMin is of type time.Duration; don't use unit-specific suffix "Min"/

// Invalid: name implies hours.
var pollHour time.Duration = time.Hour // MATCH /var pollHour is of type time.Duration; don't use unit-specific suffix "Hour"/

// Invalid: name implies microseconds.
var intervalUsec time.Duration = 500 * time.Microsecond // MATCH /var intervalUsec is of type time.Duration; don't use unit-specific suffix "Usec"/

// Valid: no unit-specific suffix.
var timeout time.Duration = 5 * time.Second

// Valid: no unit-specific suffix.
var delay time.Duration = 100 * time.Millisecond

// Valid: not a time.Duration type.
var nameSec string = "hello"

// Valid: no unit-specific suffix.
var interval time.Duration = time.Second
