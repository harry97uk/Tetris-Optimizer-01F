# Tetris-Optimizer-01F

Develop a program that receives only one argument, a path to a text file which will contain a list of tetrominoes and assemble them in order to create the smallest square possible.

The program must :

-Compile successfully   
-Assemble all of the tetrominoes in order to create the smallest square possible    
-Identify each tetromino in the solution by printing them with uppercase latin letters (A for the first one, B for the second, etc)   
-Expect at least one tetromino in the text file   
-In case of bad format on the tetrominoes or bad file format it should print ERROR    
-The project must be written in Go.   
-The code must respect the good practices.    
-It is recommended to have test files for unit testing.   
-Only the standard go packages are allowed    

Example of a text File

    #...
    #...
    #...
    #...

    ....
    ....
    ..##
    ..##

If it isn't possible to form a complete square, the program should leave spaces between the tetrominoes. For example:

    ABB.
    ABB.
    A...
    A...
