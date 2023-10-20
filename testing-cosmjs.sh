#!/bin/bash

currentPath=$(dirname "$0")

# Pretend that all future 4 it tests have failed
exit $(expr 4 + 10)

# To get the number of failures, do:
# $ echo $? 
# if it is 0 -> no errors
# if it is 11 or more -> subtract 10 and you have the number of errors.
# any other number and there was a problem with testing proper.