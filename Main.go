package main

import "flag"

type CMDFlags struct {
    binary bool
    check bool
    text bool
}


func main() {
    parsedFlags := CMDFlags{}
    flag.BoolVar(&parsedFlags.binary,"binary",false,"read in binary mode")
    flag.BoolVar(&parsedFlags.binary,"b",false,"read in binary mode")

    flag.BoolVar(&parsedFlags.check,"check",false,"read checksums from the FILEs and check them")
    flag.BoolVar(&parsedFlags.check,"c",false,"read checksums from the FILEs and check them")

    flag.BoolVar(&parsedFlags.text,"text",false,"read in binary mode")
    flag.BoolVar(&parsedFlags.text,"t",false,"read in binary mode")

    flag.Parse()



}
