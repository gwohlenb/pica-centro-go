// Pica Centro game
// Objective: To guess a secret number of a specified length in the fewest attempts
//            Each digit is in the range of 0-9; there can be duplicates
// Rules: To make an guess, enter a number of the specified length
//        In response, the program prints a string with the following character clues:
//          P = (pica) the digit is part of the answer but not in the right position
//          C = (centro) the digit is part of the answer and is in the correct position
//          X = the digit is not part of the answer
//        Continue guessing until you figure out the secret number or run out of guesses (20)
//
// To run:
// go run pica-centro.go <secret_number_length>
// where secret_number_length is optional; default is 4
// hit ESC or ^c to give up
//
// Author: gwohlenb@yahoo.com
// GitHub: https://github.com/gwohlenb
// Date:   November 2022 (v1.0)
//////////////////////////////////////

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"golang.org/x/term" // term is for collecting user input without them having to press return (stdin raw mode)
			    // must be installed with "go get -u golang.org/x/term"
			    // requires go version >=  1.17
)

var secretNumberLength int = 4 // default length if the user doesn't specify otherwise on the command line
const maxSecretNumberLength int = 10
const maxGuessCount int = 20
const asciiEsc byte = 27
const asciiCtrlC byte = 0x03
const asciiBackspace byte = 0x08
const asciiDelete byte = 0x7F // 127
const asciiP byte = 80 // pica
const asciiC byte = 67 // centro
const asciiX byte = 88 // not present
const asciiDigit0 byte = 48
const asciiDigit9 byte = 57

func main() {

  fmt.Println("Welcome to Pica Centro")
  fmt.Println("Press ESC to give up")

  if len(os.Args) > 1 {
    var err error
    secretNumberLength, err = strconv.Atoi(os.Args[1])
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    if secretNumberLength < 1 || secretNumberLength > maxSecretNumberLength {
      fmt.Println("Sorry, the secret number can only be from 1 to", maxSecretNumberLength, "digits")
      os.Exit(1)
    }
  }

  fmt.Println("The secret number is", secretNumberLength, "digits long")

  var secretNumber []int // store the secretNumber in a slice of integers
  secretNumber = generateSecretNumber(secretNumber)

  // These lines can be useful when testing
  //secretNumber = []int{0,4,4,3} // if needed for testing
  //fmt.Println("Secret number is", intSliceToString(secretNumber))

  // Now start the guessing and evaluation process
  var currentGuess []int
  var guessCount int
  var analysisString []byte
  var solved bool

  for guessCount = 1; guessCount <= maxGuessCount && !solved; guessCount++ {
    currentGuess, validGuess := collectGuess(currentGuess, guessCount)
    if (!validGuess) {
      fmt.Println("Failure! The secret number was", intSliceToString(secretNumber))
      os.Exit(1)
    }
    analysisString, solved = analyzeGuess(currentGuess, secretNumber)
    if solved {
      fmt.Println("Success! You got the secret number in", guessCount, "guess(es)!")
    } else {
      fmt.Printf("%s\n", analysisString)
    }
  }
}

//////////////////////////////////////

func generateSecretNumber(secretNumber []int) []int {

  rand.Seed(time.Now().UnixNano()) // seed the random number generator

  for i := 0; i < secretNumberLength; i++ {
    secretNumber = append(secretNumber, rand.Intn(9))
  }

  return secretNumber
}

//////////////////////////////////////

func collectGuess(currentGuess []int, guessCount int) (guess []int, validGuess bool) {

  // put stdin into raw mode (no user CR required)
  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState) // leave stdin raw mode when this function returns

  fmt.Printf("Enter %d-digit guess #%d: ", secretNumberLength, guessCount)
  for i := 0; i < secretNumberLength; i++ {
    fmt.Printf("_\b")

    digit := make([]byte, 1)
    _, err = os.Stdin.Read(digit)
    if err != nil {
      fmt.Printf("\r\n%v\r\n", err) // display the error in raw mode
      return []int {0}, false
    }

    // Check if user entered an exit character
    if (digit[0] == asciiEsc || digit[0] == asciiCtrlC) {
      fmt.Printf("\r\n")
      return []int {0}, false
    }

    // Check if user entered a backspace or delete character
    if (digit[0] == asciiBackspace || digit[0] == asciiDelete) {
      i -= 1 // don't count this as a guess digit even if the user hasn't entered anything yet
      if len(currentGuess) > 0 {

        // Remove the guess digit from the screen and the currentGuess slice
        fmt.Printf("\b")
        currentGuess = currentGuess[:len(currentGuess)-1]
        i -= 1
      }
    } else {

      // Process the guess digit entered; anything other than 1-9 is treated as a zero/0
      digitInt, _ := strconv.Atoi(string(digit))
      fmt.Printf("%d", digitInt)
      currentGuess = append(currentGuess, digitInt)
    }
  }
  fmt.Printf("\r\n")

  return currentGuess, true
}

//////////////////////////////////////

func analyzeGuess(currentGuess []int, secretNumber []int) (analysisString []byte, solved bool) {
  var isPresent bool = false
  var position int
  solved = false

  for i, value := range(currentGuess) {
    isPresent, position = analyzeCharacter(value, i, secretNumber)

    if (isPresent) {
      if i == position {
        analysisString = append(analysisString, asciiC)
      } else {
        analysisString = append(analysisString, asciiP)
      }
    } else {
      analysisString = append(analysisString, asciiX)
    }
  }

  // Check if the guess is 100% correct
  solved = true
  for i := range(analysisString) {
    if analysisString[i] != asciiC {
      solved = false
    }
  }

  return analysisString, solved
}

//////////////////////////////////////

func analyzeCharacter(currentCharacter int, currentPosition int, secretNumber []int) (isPresent bool, positionFound int) {
  isPresent = false   // true if the digit is present in the secret number
  positionFound = -1  // position digit found in the secret number; zero-indexed; invalid if isPresent is false

  for i, value := range(secretNumber) {
    if value == currentCharacter {
      isPresent = true
      positionFound = i;
      // Once the digit is found, we need to check if this is the current search position in the guess. If so, we need to stop searching.
      // This is to handle the case of secret numbers containing duplicate digits (like "5144" or "4143"))
      if currentPosition == positionFound {
        return isPresent, positionFound
      }
    }
  }
  return isPresent, positionFound
}

//////////////////////////////////////
// Utility function to convert an int slice to a string

func intSliceToString(intSlice []int) (intString string) {
  intString = ""
  for _, value := range(intSlice) {
    intString += strconv.Itoa(value)
  }
  return intString
}
