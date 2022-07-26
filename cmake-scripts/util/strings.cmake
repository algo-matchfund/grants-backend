# Converts a string to list, optionally escaping the semicolons with the
# $<SEMICOLON> generator.
# This positional arguments are:
#  variable            The name of the variable to set the result in
#  input               The input variable
# The named arguments are:
#  SEMICOLON_ESCAPE    Option to specify whether or not to escape semicolons
function(string_to_list variable input)
    cmake_parse_arguments(S "SEMICOLON_ESCAPE" "" "" ${ARGN})

    separate_arguments(list_var NATIVE_COMMAND "${input}")

    # if requested, convert the semicolons to the semicolon generator
    if (S_SEMICOLON_ESCAPE)
        string(REPLACE ";" "$<SEMICOLON>" list_var "${list_var}")
    endif()

    set(${variable} "${list_var}" PARENT_SCOPE)
endfunction()
