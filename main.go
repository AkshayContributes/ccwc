package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		readEverythingFromStdIn()
		return
	}

	switcher(args[1:])
	return

}

func switcher(args []string) {
	filePath := ""
	command := ""

	pattern := `^-(c|w|l|m)$`

	patternExclusive := `^-[a-zA-Z]$`

	re := regexp.MustCompile(pattern)
	reExclusive := regexp.MustCompile(patternExclusive)

	fmt.Println(args)
	if re.MatchString(args[0]) {
		command = strings.TrimSpace(args[0])
	} else if reExclusive.MatchString(args[0]) {
		fmt.Printf("cwwc: invalid option -- '%s'\n", args[0])
		os.Exit(1)
	} else {
		filePath = strings.TrimSpace(args[0])
	}

	switch command {
	case "-c":
		countBytes(filePath)
	case "-w":
		countWords(filePath)
	case "-l":
		countLines(filePath)
	case "-m":
		countChars(filePath)
	default:
		countEverything(filePath)
	}
}

func readFile(filePath string) (*os.File, error) {

	fmt.Println("Reading file:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error Opening File:", err)
		return nil, err
	}
	return file, nil
}

func countBytes(filePath string) (int64, error) {

	if filePath == "" {
		return countBytesFromStdIn()
	}

	return countBytesFromFile(filePath)

}

func countBytesFromStdIn() (int64, error) {
	fmt.Println("Counting Bytes from stdin")
	_, bytes := readInputFromStdIn()
	fmt.Printf("%d", bytes)
	return bytes, nil
}

func countBytesFromFile(filePath string) (int64, error) {

	bytes, _, _, _, _, err := readEverythingFromFile(filePath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0, err
	}

	fmt.Printf("%d %s", bytes, filePath)

	return bytes, nil

}

func countWords(filePath string) (int, error) {

	if filePath == "" {
		return countWordsFromStdin()
	}

	return countWordsFromFile(filePath)

}

func countWordsFromStdin() (int, error) {
	_, _, wordCount, _, _, _ := readEverythingFromStdIn()
	fmt.Printf("%d", wordCount)
	return wordCount, nil
}

func countWordsFromFile(filePath string) (int, error) {

	_, wordCount, _, _, _, err := readEverythingFromFile(filePath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0, err
	}

	fmt.Printf("%d %s", wordCount, filePath)
	return wordCount, nil
}

func countLines(filePath string) (int, error) {
	if filePath == "" {
		return countLinesFromStdin()
	}

	return countLinesFromFile(filePath)

}

func countLinesFromStdin() (int, error) {
	_, _, _, lineCount, _, _ := readEverythingFromStdIn()
	fmt.Printf("%d", lineCount)
	return lineCount, nil
}

func countLinesFromFile(filePath string) (int, error) {

	_, _, lineCount, _, _, err := readEverythingFromFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0, err
	}

	return lineCount, nil

}

func countChars(filePath string) (int, error) {
	if filePath == "" {
		return countCharsFromStdin()
	}

	return countCharsFromFile(filePath)
}

func countCharsFromStdin() (int, error) {
	_, charCount, _, _, _, _ := readEverythingFromStdIn()
	fmt.Printf("%d", charCount)
	return charCount, nil
}

func countCharsFromFile(filePath string) (int, error) {

	_, _, _, charCount, _, err := readEverythingFromFile(filePath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0, err
	}

	fmt.Printf("%d %s", charCount, filePath)

	return charCount, nil
}

func countEverything(filePath string) {
	bytes, lineCount, wordCount, charCount, fileName, err := int64(0), 0, 0, 0, "", error(nil)
	if filePath != "" {
		bytes, lineCount, wordCount, charCount, fileName, err = readEverythingFromFile(filePath)

	} else {
		bytes, lineCount, wordCount, charCount, fileName, err = readEverythingFromStdIn()
	}

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("%d %d %d %d %s", bytes, lineCount, wordCount, charCount, fileName)

}

func readEverythingFromFile(filePath string) (int64, int, int, int, string, error) {

	file, err := readFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0, 0, "", err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return 0, 0, 0, 0, "", err
	}

	bytes := fileInfo.Size()

	scanner := bufio.NewScanner(file)
	wordCount := 0
	lineCount := 0
	charCount := 0

	for scanner.Scan() {
		wordCountTemp, charCountTemp := countWordsInLine(scanner.Text())
		wordCount += wordCountTemp
		charCount += charCountTemp
		lineCount++

	}

	return bytes, lineCount, wordCount, charCount, fileInfo.Name(), nil
}

func countWordsInLine(line string) (int, int) {
	wordCount := 0
	charCount := 0
	inWord := false
	for _, char := range line {
		charCount++
		if unicode.IsSpace(char) {
			inWord = false
		} else {
			if !inWord {
				wordCount++
				inWord = true
			}
		}
	}
	return wordCount, charCount
}

func readEverythingFromStdIn() (int64, int, int, int, string, error) {
	fmt.Println("Reading from stdin")
	byteCount, charCount, wordCount, lineCount, err := int64(0), 0, 0, 0, error(nil)

	inputLines, byteCount := readInputFromStdIn()

	for _, line := range inputLines {
		wordCountTemp, charCountTemp := countWordsInLine(line)
		wordCount += wordCountTemp
		charCount += charCountTemp
		lineCount++
	}

	return byteCount, charCount, wordCount, lineCount, "", err
}

func readInputFromStdIn() ([]string, int64) {
	var inputLines []string
	bytes := int64(0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		bytes += int64(len(line))
		inputLines = append(inputLines, line)
	}
	return inputLines, bytes
}
