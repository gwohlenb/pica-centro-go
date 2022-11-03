# pica-centro-go

Pica Centro game

Objective:
* To guess a secret number of a specified length in the fewest attempts
* Each digit is in the range of 0-9; there can be duplicates

Rules:
* To make an guess, enter a number of the specified length
* In response, the program prints a string with the following digit clues:
  * P = (pica) the digit is part of the answer but not in the right position
  * C = (centro) the digit is part of the answer and is in the correct position
  * X = the digit is not part of the answer
* Continue guessing until you figure out the secret number or run out of guesses (default 20)

To Run:\
**go run pica-centro.go <secret_number_length>**
* where secret_number_length is optional; default is 4
* hit ESC or ^c to give up
