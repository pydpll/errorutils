# Errorutils

Errorutils is a reusable optional functions error framework that extends the zerolog package. It provides a simple and flexible way to instantiate custom error types with additional information.

## Installation

To use Errorutils in your Go project, you can install it using `go get`:

```
> go get github.com/TgenNorth/errorutils
```

## Usage

To use Errorutils in your project, you first need to import it:

```go
package myPackage
import "github.com/TGenNorth/errorutils"
```

### When to use errorutils

Errorutils enables error checking and handling at the location where the error is generated or captured. In general, our recommendation is to use `ExitOnFail()` in places where `panic()` or `log.fatal()` could be used, `WarnOnFail()` provides notifications to users in cases where the program could continue operating. These functions do not need to be wrapped in if conditions to check for nil.

### When not to use errorutils

The functionality of this package is constrained to only handling failure errors. The following are some examples where the use of alternatives is encouraged:

- In places where it makes more sense to return the error, use the common construct `if err != nil`.
- In code tests, do not replace any of the error and logging functionality of a testing object.
- Sharing information such as EOF should not be handled with this framework.
- Recoverable panics or terminations that expect the defer stack to be executed should rely on built-in `panic()` instead.
- The `Details` error type is not meant to be compared.

### Creating a new error with line references

To add details to an error, use the `New` function. This function takes an error and functions of type `Option` that take in the informational values.

```go
err := myfunc()
if err !=nil {
    detailed := errorutils.New(err, errorutils.WithExitCode(3), errorutils.WithLineRef("OKP8PK1CosD"))
    return detailed
}
```

Errorutils provides a way to add line references to error values that are only printed when `zerolog.SetGlobal.DebugLevel` is enabled. Line references indicate the location in the code where the error occurred. Ideally, unique identifiers such as random strings are better to avoid outdating the reference.

### Error handling abstraction

Errorutils offers error handling that replaces simple cases of error checking with if statements. For more complex error handling use `errorutils.HandleFailure` as described below.

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

The following example based on TGenNorth/kmare/database, stashes the result if writing fails.

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

Errorutils is released under our custom Academic and Research license. See [LICENSE](LICENSE.rst) for details.