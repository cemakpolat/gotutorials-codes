Okay, let's implement the second example from Chapter 9, the "Word Frequency Counter," as a complete project with tests. We'll follow the same structure, setup instructions, and testing approach as before.

**Project: Word Frequency Counter with Tests**

This project will create a command-line tool that reads a text file, counts the frequency of each word, and prints the top N most frequent words. It will include unit tests to ensure its functionality.

**Project Structure:**

```
wordcounter/
├── go.mod
├── main.go
└── counter/
    ├── counter.go
    └── counter_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir wordcounter
    cd wordcounter
    ```

2.  **Initialize the Go Module:**
    Inside the `wordcounter` directory, run:

    ```bash
    go mod init wordcounter
    ```

3.  **Create the `counter` Package Directory:**

    ```bash
    mkdir counter
    ```

4.  **Create the `counter.go` file (inside the `counter` directory):**

    ```bash
    touch counter/counter.go
    ```

5.  **Create the `counter_test.go` file (inside the `counter` directory):**

    ```bash
    touch counter/counter_test.go
    ```

6.  **Create the `main.go` file (inside the `wordcounter` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
wordcounter/
├── go.mod
├── main.go
└── counter/
    ├── counter.go
    └── counter_test.go
```

**Now, copy the following code into their corresponding files:**

**1. `go.mod` File:**

```
module wordcounter

go 1.21
```

**2. `counter/counter.go` (Word Counter Logic):**

```go
// counter/counter.go
package counter

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func CountWords(filename string) (map[string]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    wordCounts := make(map[string]int)
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)
    fmt.Println("wordCounts");
    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        
        wordCounts[word]++
    }

	if err = scanner.Err(); err != nil {
		return nil, err;
	}
    return wordCounts, nil
}

func PrintTopWords(wordCounts map[string]int, topN int) {
    type kv struct {
        key   string
        value int
    }

    var ss []kv
    for k, v := range wordCounts {
        ss = append(ss, kv{k, v})
    }

    sort.Slice(ss, func(i, j int) bool {
        return ss[i].value > ss[j].value
    })
	
	
    for i := 0; i < topN && i < len(ss); i++ {
		fmt.Printf("%s: %d\n", ss[i].key, ss[i].value);
    }
}
```
* Contains functions for counting and printing the most frequent words.
*  Uses `bufio.Scanner` to read the file word by word, ignoring any punctuation.
* The functions are exported, so they can be used in the tests and the `main.go` file.

**3. `counter/counter_test.go` (Word Counter Logic Tests):**

```go
// counter/counter_test.go
package counter

import (
	"os"
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	testContent := "This is a test file. This file has some text. This text is for testing purposes."

	filename := "testfile.txt";
	err := os.WriteFile(filename, []byte(testContent), 0644);
	if err != nil {
		t.Fatalf("Error creating the file: %v", err);
	}
	defer os.Remove(filename)

	wordCounts, err := CountWords(filename)
	if err != nil {
		t.Fatalf("CountWords failed: %v", err);
	}

	expectedWordCounts := map[string]int{
		"this": 3,
		"is": 2,
		"a": 1,
		"test": 1,
		"file": 2,
		"has": 1,
		"some": 1,
		"text": 2,
		"for": 1,
		"testing": 1,
		"purposes": 1,
	}

	if !reflect.DeepEqual(wordCounts, expectedWordCounts) {
		t.Errorf("CountWords failed: expected: %v, got: %v", expectedWordCounts, wordCounts);
	}
}
```

*   Tests the logic for the `CountWords` function.
* Creates a test file and tests different types of words, while asserting that the values are correct.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"wordcounter/counter"
)

func main() {
    if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <filename> <topN>")
		os.Exit(1);
    }
    filename := os.Args[1];
	topN, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid topN value!");
		os.Exit(1);
	}

    wordCounts, err := counter.CountWords(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error counting words: %v\n", err)
        os.Exit(1)
    }

    counter.PrintTopWords(wordCounts, topN)
}
```
*   Imports the `counter` package and uses the exported functions.
*   Uses command-line arguments to read the filename and top N.
*   Outputs to standard error if an issue happens while parsing the arguments.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `wordcounter` directory, and run:

    ```bash
    go run . mytext.txt 5
    ```

	Create a file named `mytext.txt` with some text content before running this command.
    Replace `mytext.txt` with the name of your input file, and 5 with the desired `topN` value.

2.  **Run the Tests:**
    Open a terminal, navigate to the `wordcounter` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      wordcounter/counter        0.004s
```

**Output (If Tests Fail):**
If any tests fail, the output will provide details about the error, including the expected output, the actual output and the line of the code where it failed.

```
--- FAIL: TestCountWords (0.00s)
    counter_test.go:37: CountWords failed: expected: map[file:1 is:1 this:1], got: map[file:2 is:2 this:2]
FAIL
exit status 1
FAIL	wordcounter/counter	0.003s
```

**Key Features of This Project:**

*   **Complete CLI Tool:** Implements the word frequency counter functionality.
*   **Modularity:** Uses a `counter` package to encapsulate the core logic.
*   **Testing:** Includes unit tests to verify the correctness of the counting and sorting logic.
*  **Arguments**: Reads arguments passed from the command line, parsing them and verifying their values.

This project provides another practical example of how to build a command-line application in Go and demonstrates best practices for modularity, testing, and command line arguments.
