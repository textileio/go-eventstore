# API Rationale

A `Model` has the concept of explicit and implicit transactions, which can be readonly 
or allow writes.

### Model IDs
The idea is that every `struct` that will be registered as a model, should have a 
property named `ID` with type `eventstore.EntityID`. When registering the model 
this will be enforced, and will be automatically generated if has an empty value. 

## Examples
Here're some examples to understand a bit more.

### Implicit transactions
Under the hood creates a transaction to execute the action:
```go
store := NewStore(...)
model, _ := store.Register("Person", &Person{})

p := &Person{Name: "Alice", ...}
_ = model.Create(p)

p.Name = "Bob"
_ = model.Save(p) 

shouldBeTrue := model.Has(p.ID)

_ = model.Delete(p.ID)

```
Notes:
* Pros: Easy to Use
* Cons: Not much guarantees what happens between operations.

### Explicit transactions
Are separated between `ReadTxn` and `WriteTxn` ones. `ReadTxn` transactions only 
allow operations that don't mutate data. `WriteTxn` transactions allow read and 
write operations.

#### ReadTxn
```go
store := NewStore(...)
model, _ := store.Register("Person", &Person{})

p := &Person{Name: "Alice", ...}
_ = model.Create(p)

_ = model.ReadTxn(func(txn *Txn) error {
    shouldBeTrue := txn.Has(p.ID) // Good, only reads.
    
    p2 := &Person{}
    txn.FindByID(p.ID, p2) // Good, only reads

    p2.Name = "Bob"
    _ = txn.Save(p2) // FAIL
    _ = txn.Delete(p2.ID) // FAIL

    return nil
})
```

#### WriteTxn
Same example as above:
```go
store := NewStore(...)
model, _ := store.Register("Person", &Person{})

p := &Person{Name: "Alice", ...}
_ = model.Create(p)

_ = model.WriteTxn(func(txn *Txn) error {
    shouldBeTrue := txn.Has(p.ID) // Good, only reads.
    
    p2 := &Person{}
    txn.FindByID(p.ID, p2) // Good, only reads

    p2.Name = "Bob"
    _ = txn.Save(p2) // Good
    _ = txn.Delete(p2.ID) // Good

    return nil
})
```

As can be seen, explicit txns may feel more bloated when used but should 
make sense in most logic that wants isolation guarantees during 
business-logic operations. The separation between `Read` and `Update` 
txn types is more of a safeguard for the developer than anything else.
