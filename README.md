# Errorutils

Errorutils is a reusable optional functions error framework that extends the logrus package. It provides a simple and flexible way to instantiate custom error types with additional information.

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

### Creating a new error with Line references

To create a new error with additional information, you can use the `New` function. This function takes an error and functions of type `Option` that take in the informational values.

```go
err := errors.New("something went wrong")
details := errorutils.New(err, errorutils.WithExitCode(1), errorutils.WithLineRef("OKP8PK1CosD"))
```

Errorutils provides a way to add line references to error values that are only printed when logrus.DebugLevel is enabled. Line references indicate the location in the code where the error occurred. Ideally, unique identifiers such as random strings are better to avoid outdating the reference.

### Creating a new error report

Errorutils has a succinct syntax for creating new `Detail` errors without having to use pre-made error values.

```go
details := errorutils.NewReport("something went wrong", "OKP8PK1CosD")
```

### Handler functions and safe closer

Errorutils accepts handler functions that deal with errors consistently. Additionally, the package offers safe closing functions that visibilize closing errors for io.Closer instances.

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
### Error handling abstraction

Errorutils offers error handling that supplants simple cases of traditional error checking with if statements. for more complex error handling use `errorutils.HandleFailure` as described above

```go
err := fmt.Error("this is a simple error")
errorutils.LogFailures(err)
errorutils.WarnOnFail(err)
errorutils.WarnOnFailf(err, "additional info: %%s")
errorutils.PaniconFail(err)
```

## License

Errorutils is released under our custom Academic and Research license. See [LICENSE](LICENSE.rst) for details.