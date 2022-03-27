# wordle-guesser
A tool to guess possible Wordle words.
It filters out words from a given dictionary with possible constraints.

## How to Build
```
go build
```

## Usage
```
Usage of ./wordle-guesser:
  -correct-spot string
        Characters in the correct spots.
        Format : <position1>:<characters>;<position2>:<characters>,...
        Example: 1:e;2:p;3:o
  -dictionary string
        Path to the dictionary file. (default "sgb-words.txt")
  -invalid string
        Invalid characters.
        Format : <chars>
        Example: t,a,s,d
  -wrong-spot string
        Characters in the wrong spots.
        Format : <position1>:<characters>;<position2>:<characters>,...
        Example: 2:e;3:p,e;4:o
```
