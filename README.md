# Errorutils

Errorutils is a reusable optional functions error framework that extends the logrus package. It provides a simple and flexible way to instantiate custom error types with additional information, sets up a custom logger for better formatting, and provides descriptive one-line checks for errors.

## Installation

To use Errorutils in your Go project, you can install it using `go get`:

```
> go get github.com/pydpll/errorutils
```

## Usage

To use Errorutils in your project, you first need to import it. It's init funtion will make the custom formatter accessible through normal logrus entries.

```go
package myPackage
import "github.com/pydpll/errorutils"
```

### When to use errorutils

Errorutils enables error checking and handling at the location where the error is generated or captured. In general, our recommendation is to use `ExitOnFail()` in places where `logrus.Error()` could be used, `WarnOnFail()` provides notifications to users in cases where the program could continue operating. These functions are nil-error checks.

### When not to use errorutils

The functionality of this package is constrained to only handling failure errors. The following are some examples where the use of alternatives is encouraged:

- In places where it makes more sense to return the error, use the common construct `if err != nil`.
- When defered calls are necessary for cleanup or other tasks, use `panic()`. This library uses os.Exit().
- In control flow scenarios where there should be conditional execution of keywords such as `continue` and `break` based on the error.
- In code tests, do not replace any of the error and logging functionality of a testing object.
- When Sharing information such as EOF should be not be handled with this custom type, some workarounds by wrapping errors `WithInner()` might work but it is not a guarantee (see next section example).
- Recoverable panics or terminations that expect the defer stack to be executed should rely on built-in `panic()` instead.
- The `Details` error type is not meant to be compared.

### Creating a new error with line references

To add details to an error, use the `New` function. This function takes an error and functions of type `Option` to add information, and wrap other error types. This is a nil check. If both errors are provided (by argument and inner error option) and are nil the returned value is a nil *Details.

```go
//err is nil, otherErr is not
detailed := errorutils.New(err, errorutils.WithExitCode(3), errorutils.WithLineRef("OKP8PK1CosD"), errorutils.WithInner(otherErr))
// detailed is now showing the inner error otherErr.Error() as the message. Type information has been lost.
```

Errorutils provides a way to add line references to error values that are only printed when `logrus.DebugLevel` is enabled. Line references indicate the location in the code where the error occurred. Ideally, unique identifiers such as random strings are better to avoid outdating the reference. Alternatively assign the information to present the offending input or other useful information.

### Error handling abstraction

As mentioned before, Errorutils offers error handling that replaces simple cases of error checking. For more complex error handling use `errorutils.HandleFailure` as described below.

```go
err := myFunc()
//handle with one of these
errorutils.LogFailures(err)
errorutils.LogFailuresf(err, "other error info: %%s")
errorutils.WarnOnFail(err)
errorutils.WarnOnFailf(err, "additional info: %%s")
errorutils.ExitonFail(err)
```

### Handler functions and safe closer

Errorutils accepts handler functions that deal with errors consistently. Additionally, the package offers safe closing functions that visibilize closing errors for Closers.

The following example stashes the result if writing fails.

```go
func x() {
    //... this function did some valuable work that needs to be saved.
    seqFile, err := os.Create(filepath.Join(libLoc, fmt.Sprintf("%s.fasta", name)))
        if err != nil {
            goto handleError
        }

    {
        defer errorutils.NotifyClose(seqFile)
        //bufio writer
        seqWriter := bufio.NewWriter(seqFile)
        _, err = seqWriter.Write(formattedBytes)
        indexWriter.Flush()
    }
handleError:
    // if file cannot be created or there is a writting error, stash the sequences
    stashingErr := errorutils.HandleFailure(
                    err,
                    errorutils.Handler(func() *Details {
                        r, err2 := writeTemp(sha1)
                        if err2 != nil {
                        return errorutils.New(err2)
                        }
                        _, err2 = r.Write(formattedBytes)
                        errorutils.SafeClose(r, &err2)
                        return errorutils.New(err2)
                        }),
                    errorutils.WithMsg(fmt.Sprintf("sequences file could not be created for %s at %s, a stash was ATTEMPTED as temporaryfile accessible with hash name %s", name, libLoc, sha1)),
                    errorutils.WithLineRef("uDIKN3XCREp"))
    // in this case format string must have a scaped string verb '%%s' to ensure WarnOnFailf will have a place to print error value.
    errorutils.WarnOnFailf(stashingErr, fmt.Sprintf("Sequences for %s cound not be saved: %%s\nSkipping...", name), errorutils.WithLineRef("XqZsHJI8ABs"))
}
```

## License

Errorutils is released under MIT Licensing. see [LICENSE](LICENSE) for details.
