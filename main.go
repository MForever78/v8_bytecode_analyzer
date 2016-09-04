package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "sort"
    "encoding/csv"
    "strconv"
)

type Result struct {
    Bytecodes string
    Count int
}

type ByCount []Result

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }

func (r Result) String() string {
    return fmt.Sprintf("%s: %d\n", r.Bytecodes, r.Count)
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: ./main /path/to/bytecode /path/to/output")
        return
    }

    file, err := os.Open(os.Args[1]);
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    dictionary := make(map[string]int)
    sequence := make([]string, 128)
    sequenceIndex := 0

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        if !strings.Contains(line, "@") {
            continue
        }
        line = line[40:]

        words := strings.Split(line, " ")
        var bytecode string
        for i := 0; i < len(words); i++ {
            bytecode = words[i]
            if (len(bytecode) > 2) {
                break
            }
        }

        sequence[sequenceIndex] = bytecode
        sequenceIndex++;
        sequenceString := strings.Join(sequence[:sequenceIndex], ",")
        if sequenceIndex == 128 {
            sequenceIndex = 0
        }

        count, ok := dictionary[sequenceString]
        if !ok {
            // add new table item
            dictionary[sequenceString] = 1
            sequenceIndex = 0;
        } else {
            dictionary[sequenceString] = count + 1
        }
    }

    if err:= scanner.Err(); err != nil {
        log.Fatal(err)
        return
    }

    result := make([]Result, len(dictionary))

    i := 0;
    for bytecode, count := range dictionary {
        result[i].Bytecodes = bytecode
        result[i].Count = count
        i += 1
    }

    sort.Sort(ByCount(result))

    records := make([][]string, len(result));
    for index, r := range result {
        records[index] = make([]string, 2)
        records[index][0] = r.Bytecodes
        records[index][1] = strconv.Itoa(r.Count)
    }

    output, err := os.Create(os.Args[2])
    if err != nil {
        log.Fatal(err)
        return
    }
    defer output.Close()

    w := csv.NewWriter(output)
    w.WriteAll(records)
}

