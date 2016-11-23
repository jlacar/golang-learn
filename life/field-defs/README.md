# Conway's Game of Life - field definition file format

## Purpose

The field definition file is an easy way to configure a two-dimensional
board for Conway's Game of Life.

Each line of the field definition file is parsed and interpreted based on
notations at the start of the line. Notations include comments, column offset,
relative row, and absolute row. Notations are separated from the cell configuration
by a colon.

### General line format

A line in the field definition file can have the following forms:

    # comment

    NN: ... cell configuration ...

    ++: ... cell configuration ...

    >>:NN

The first form is a comment line.

The second form is a cell configuration line with an absolute row.

The third form is a cell configuration line with a relative row.

The fourth form is a column offset setting line.

Cell configurations are determined by whatever comes after the ":" separator
in the second and third forms. Any non-space characters can be used to denote
live cells. Spaces are used to denote dead cells and need only be included to
offset live cells from the leftmost column.

The column offsets of the cells are determined by their relative position from
the colon separator. The first position after the colon is column 0 by default.
The column offset can be changed by using a column offset setting line.

Here is a full example of a configuration for Gosper's Glider Gun:

    # Gosper's Glider Gun configuration
    #
    #  0...4....9....4....9....4....9....4....
    01:@
    >>:60
    01:@
    >>:05
    05:                         @
    ++:                       @ @
    ++:             @@      @@            @@
    ++:            @   @    @@            @@
    ++: @@        @     @   @@
    ++: @@        @   @ @@    @ @
    ++:           @     @       @
    ++:            @   @
    ++:             @@

Note the use of the column offset setting and the repeated absolute row.
Those lines set the horizontal boundaries of the field.

If a field definition requires the field boundaries to be larger than what is
specified on the command line with the -x and -y options, those options will
be superseded by the minimum width and height of a field that can accommodate
the population defined in the field definition file.

### Comments

A line that starts with "#" is treated as a comment line. Use comment lines to
document the field definition or temporarily disable a configuration line.

Example:

    # These are parsed as comment lines
    #
    # Gosper's Glider Gun
    #
    # column guide
    # 0...4....9....4....9....4....9....4....9
    #
    # following configuration lines are commented out
    #>>:50
    #10: @  @@@@@

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
setting is parsed. Don't put anything after the column number as this will
cause a parsing error and the column offset setting on that line will be ignored.

    # Sets the column offset to 20
    >>:20
    # This cell will be located at row 3, column 21
    #  0123456789.....
    03: @
    >>:35
    # These cells will be located at (row:4, col:36) and (row:4, col:38)
    ++: @ @
