package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type CMDFlags struct {
    binary bool
    check bool
    text bool
}

func calcMD5sum (path string) string {

    f,err := os.ReadFile(path)

    if err != nil {
        fmt.Printf("could not open file %s", path)
        panic(err)
    }

    content := string(f)
    return MD5Sum(content)
}

func readMDFile(path string) {
    f,err := os.Open(path)

    if err != nil {
        fmt.Printf("could not open file %s", path)
        panic(err)
    }
    defer func() {
        if err := f.Close(); err != nil {
            fmt.Println("could not close file")
            panic(err)
        }
    }()

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        line := scanner.Text()
        words := strings.Fields(line)

        if (len(words) != 2) {
            fmt.Printf("md5sum: %s: no properly formatted checksum lines found\n", path)
            return
        }
        if (len(words[0]) != 32) {
            fmt.Printf("md5sum: %s: no properly formatted checksum lines found\n", path)
            return
        }

        hash := calcMD5sum(words[1])

        fmt.Print(words[1])
        fmt.Print("\t")
        if (hash == words[0]) {
            fmt.Print("OK")
        } else {
            fmt.Print("Failed")
        }
        fmt.Println()
    }


}

func handleCheckOption(paths []string ) {
    for _, path := range paths {
        readMDFile(path)
    }
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

    if (parsedFlags.check) {
        handleCheckOption(flag.Args())
    } else {
        
        for _, path := range flag.Args() {
            hash := calcMD5sum(path)
            fmt.Printf("%s %s\n", hash,path)
        }
        
    }
    


}
