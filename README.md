# Go Crossword

Package to generate a valid crossword puzzle from a list of words/clues.

Loosely based on: https://stackoverflow.com/a/22256214

The package does not guarantee all words can be placed in the grid but in general will 
usually get there if it's possible to do so (given enough attempts). 

Example: 

```bash
  $ go run cmd/generate/main.go example/words.json  
```

This will run 50 attempts and render the result with the most placed words.

If the crossword is being solved interactively you would need to store the
generated `Crossword` (e.g. json encode it to a file). This can easily 
be decoded and rendered without altering the layout.