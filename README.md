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

To create a new error with additional information, you can use the `New` function. This function takes an error and functions that match the type `Option` that takes in the informational values.

```go
err := errors.New("something went wrong")
details := errorutils.New(err, errorutils.WithExitCode(1), errorutils.WithLineRef("OKP8PK1CosD"))
```

Errorutils provides a way to add line references to error values. Line references are useful for debugging and can be used to indicate the location in the code where the error occurred. Ideally, unique identifiers such as random strings are better to avoid outdating the reference.

### Creating a new error report

Errorutils also provides a succinct syntax for creating new error reports without having to use a premade error values.

```go
details := errorutils.NewReport("something went wrong", "OKP8PK1CosD")
```

### Handler functions and safe closer

Errorutils also provides handler functions that can be used to deal with errors in a consistent way. Additionally, the package offers a safe close function for \*os.file that visibilizes closing errors.

The following example based on TGenNorth/kmare/database, stashes the result if writing fails.

```go
func x() {
    //... this function did some valuable work that needs to be saved.
    seqFile, err := os.Create(filepath.Join(libLoc, fmt.Sprintf("%d.fasta", name)))
        if err != nil {
            goto handleError
        }

    {
        defer errorutils.SafeClose(seqFile, &err)
        //bufio writer
        seqWriter := bufio.NewWriter(seqFile)
        _, err = seqWriter.Write(formattedBytes)
        indexWriter.Flush()
    }
handleError:
    // if file cannot be created or there is a writting error, stash the sequences
    if err != nil {
    stashingErr := errorutils.HandleFailure(
                     err,
                     errorutils.Handler(func() error {
                         r, err2 := writeTemp(sha1)
                         if err2 != nil {
                         return err2
                         }
                         _, err2 = r.Write(formattedBytes)
                         errorutils.SafeClose(r, &err2)
                         return err2
                         }),
                     errorutils.WithMsg(fmt.Sprintf("sequences file could not be created for %s at %s, a stash was ATTEMPTED as temporaryfile accessible with hash name %s", name, libLoc, sha1)),
                     errorutils.WithLineRef("uDIKN3XCREp"))
    if stashingErr != nil {
         errorutils.LogFailures(stashingErr, errorutils.WithLineRef("XqZsHJI8ABs"))
    }
    }
    return err
}
```


## License

Errorutils is released under our custom Academic and Research license. See [LICENSE](LICENSE.rst) for details.
