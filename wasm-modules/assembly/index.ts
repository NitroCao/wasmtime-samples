// The entry file of your WebAssembly module.

@external("", "")
declare function doDouble(n: i32): i32

@external("foo", "bar")
declare function plusOne(n: i32): i32

export function add(a: i32, b: i32): i32 {
  a = doDouble(a);
  b = plusOne(b);
  return a + b;
}
