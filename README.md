# INIParser Library

The INIParser library is a tool developed in Golang for parsing and manipulating INI (Initialization) files.

## Features

The INIParser library offers the following key features:

1. **INI File Reading**: It allows you to read the contents of an INI file from the file system while preserving the relative order of the data. The library maintains the original order of sections, key-value pairs, and comments within the INI file.

2. **INI String Reading**: Similarly, when parsing an INI string, the library ensures that the relative order of sections, key-value pairs, and comments is maintained.

3. **get & set**: The library provides the ability to modify an INI file while preserving the order of sections and key-value pairs. Also provides the ability to retrieve data from INI file

4. **Section Retrieval**: You can retrieve a list of section names present in the INI file.

5. **Serialization to Map**: INI files can be serialized into a map structure, where sections are represented as keys, and the corresponding key-value pairs are stored as nested maps.

6. **Saving to File**: Once you have made changes to the INI data, you can save it to a file.

7. **Conversion to String**: The library allows you to convert the parsed or modified INI data back into a string format.

## Assumptions

To ensure proper usage and understanding of the INIParser library, the following assumptions have been made:

1. **No Global Keys**: The library assumes that all keys must be part of a section. There are no global keys allowed in the INI files or strings. Each key-value pair should be associated with a specific section.

2. **Key-Value Separator**: The key-value pairs in the INI files or strings are separated by the equals sign (=) character. This separator indicates the assignment of a value to a key.

3. **Trimmed Keys, Values, and Section Headers**: The library trims any leading or trailing spaces from keys, values, and section headers. Additionally, empty keys, values, and section headers are not allowed. They must contain at least one non-space character.

4. **Comments**: Comments in the INI files or strings are only valid when placed at the beginning of a line. They are denoted by a semicolon (;) character. Comments can be used to provide explanatory or informational text but do not affect the parsing or functionality of the library.

Please consider these assumptions while using the INIParser library in your project.
