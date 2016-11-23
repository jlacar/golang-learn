# Conway's Game of Life - field definition file format

## Purpose

The field definition file is an easy way to configure a two-dimensional
board for Conway's Game of Life.

Each line of the field definition file is parsed and interpreted based on
notations at the start of the line. Notations include comments, column offset,
relative row, and absolute row. Notations are separated from the cell configuration
by a colon.

### General line format

A line in the field format file can have two general forms:

    # comment
    XX: ... cell configuration ...

The first line is a comment line and is described in the next section.

In the second line, XX is a notation that can be one of ">>", "++", or a number.
These notations are described in detail below.

The cell configuration part of the line is whatever comes after the ":" separator.
Any non-space characters are considered to be live cells. Any spaces are dead
cells. The column positions of the cells are determined by their relative
position from the colon separator. The first position after the colon is column
0 default. This can be changed by using the column offset notation.

### Comments

A line that starts with "#" is treated as a comment line

Example:

    # These are parsed as comment lines
    #
    # Gosper's Glider Gun
    #
    # column guide
    # 0...4....9....4....9....4....9....4....9

### Absolute row

A line that starts with a number indicates the 0-based row number the cells on
that line will be placed.

    # Puts live cells on row 15
    15: @@@  @@@@   #####

### Relative row

A line that starts with a "++" indicates that the cells on that line are to
be placed on the row immediately after that of the last line parsed.

    15: @@@  @@@@   #####
    #
    # These will go on row 16
    ++: @@    @@@@     @@@@
    # These will go on row 17
    ++: @     @        @
    # These will go on row 1
    01:   @ @
    # These will go on row 2
    ++:   @ @
    # These will go on row 3
    ++:   @@@

Note that the above example shows that parsing of lines goes from top to bottom
and that absolute row numbers don't have to be sequential. Your configuration
can jump around, although too much of that can make the configuration confusing.

### Column offset

A line that starts with ">>:NN" controls the offset of the columns of any
subsequent lines parsed. The offset remains in effect until another offset
setting is parsed.

    # Sets the column offset to 20
    \>>:20
    # This cell will be located at row 3, column 21
    #  0123456789.....
    03: @
    \>>:35
    # These cells will be located at (row:4, col:36) and (row:4, col:38)
    ++: @ @
