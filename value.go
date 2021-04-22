package anyiter

// SafeValue is an interface wrapper for the methods of reflect.Value, but with the methods that can panic
// modified to return errors instead in the cases where they would panic.
type SafeValue interface {
}
