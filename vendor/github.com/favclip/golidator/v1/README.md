# golidator

The programable validator.

## Description

Existing validator is not fully flexible.
We want to customize error object generation.
We want to modify value instead of raise error when validation failed.
We want it.

## Samples

see [usage](https://github.com/favclip/golidator/blob/master/v1/usage_test.go)

### Basic usage

```
parent := golidator.NewValidator()
err := parent.Validate(obj)
```

### Use Custom Validator

```
parent := golidator.NewValidator()
parent.SetValidationFunc("req", func(t *validator.Target, param string) error {
    val := t.FieldValue
    if str := val.String(); str == "" {
        return fmt.Errorf("unexpected value: %s", str)
    }

    return nil
})
```

### Use Customized Error

```
parent := &golidator.Validator{}
parent.SetTag("validate")
parent.SetValidationFunc("req", validator.ReqFactory(&validator.ReqErrorOption{
    ReqError: func(f reflect.StructField, actual interface{}) error {
        return fmt.Errorf("%s IS REQUIRED", f.Name)
    },
}))
```
