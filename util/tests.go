//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package util

var Tests = map[int]string{
	1: `print("Success!")`,
	2: `console.log("Success!");`,
	3: `echo "Success!"`,
	4: `print "Success!";`,
	5: `print("Success!")`,
	6: `puts "Success!"`,
	7: `
#include <stdio.h>
int main() {
   printf("Success!");
   return 0;
}`,
	8: `
#include <iostream>
int main() {
   std::cout << "Success!";
   return 0;
}`,
	9: `fn main() { println!("Success!"); }`,
	10: `
program hello
 print *, "Success!"
end program hello`,
	11: `main = putStrLn "Success!"`,
	12: `
public class HelloWorld {
   public static void main(String[] args) {
       System.out.println("Success!");
   }
}`,
	13: `
package main
import "fmt"
func main() { fmt.Println("Success!") }`,
	14: `
let message: string = "Success!";
console.log(message);`,
	15: `(write-line "Success!")`,
	16: `
#lang racket
"Success!"`,
	17: `puts "Success!"`,
	18: `(println "Success!")`,
	19: `
section .text
global _start
_start:
   mov edx, len
   mov ecx, msg
   mov ebx, 1
   mov eax, 4
   int 0x80
   mov eax, 1
   int 0x80
section .data
msg db "Success!", 0xa
len equ $ -msg`,
	20: `
const std = @import("std");
pub fn main() !void {
   std.io.getStdOut().writeAll("Success!") catch unreachable;
}`,
	21: `echo "Success!"`,
	22: `
import std.stdio;
void main() { writeln("Success!"); }`,
	23: `Console.WriteLine("Success!");`,
	24: `print("Success!")`,
	25: `void main() { print("Success!"); }`,
	26: `
Module Program
    Sub Main()
         Console.WriteLine("Success!")
    End Sub
End Module`,
	27: `printfn "Success!"`,
	28: `<?php echo "Success!";`,
}
