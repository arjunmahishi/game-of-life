# Game of Life

A simple(ish) program that simulates [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) on the terminal. 

## Running it

### Build the binary

```bash
$ make build
```

### Simulate an existing pattern

There are a few patterns saved in the [./patterns](./patterns) directory. These patterns are simple text files with a 2D matrix of `1`s and `0`s representing the initial state of the simulation

```bash
$ ./bin/life patterns/butterfly-hatch
```

### Create your own pattern

Create a pattern file similar to the ones in [./patterns](./patterns). Just create a 2D matric with `0`s representing dead cells and `1`s representing alive cells.

**Example**

```text
$  cat canvas.txt

00000000000000000
00000000000000000
00011111111111000
00000000000000000
00000000000000000
```
Dont worry about the size of the matrix. During the simulation the size will be adjusted to center the pattern in your terminal window.

**Run it**

```bash
$ ./bin/life canvas.txt
```

**Output**

![](https://i.imgur.com/aGp5dHN.gif)

### More options

```bash
$ ./bin/life --help

Usage of ./bin/life:
  -canvas-only string
        Empty canvas
  -empty string
        Empty cell character (default " ")
  -filled string
        Filled cell character (default "■")
  -freq int
        Refresh frequency (default 100)
  -line-space string
        Line space character
  -no-color
        Disable colors
```
