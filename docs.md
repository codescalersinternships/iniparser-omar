# Documentation

#### `func (*INIParser) LoadFromString(str string)`

This method parses the provided INI string `str` and loads its contents into the INIParser object. It allows you to read and manipulate the INI data directly from a string.

#### `func (*INIParser) LoadFromFile(filePath string)`

This method reads the contents of the specified INI file `filePath` from the file system and loads them into the INIParser object. It enables you to work with the data stored in an INI file.

#### `func (*INIParser) GetSectionNames() []string`

This method retrieves a list of all the section names present in the INIParser object. It returns a slice of strings containing the names of all sections in the INI data.

#### `func (*INIParser) GetSections() map[string]map[string]string`

This method serializes the INIParser object into a dictionary/map structure. Each section name is a key in the outer map, and the corresponding value is an inner map representing the key-value pairs within that section.

#### `func (*INIParser) Get(sectionName, key string) string`

This method retrieves the value of the specified `key` within the specified `sectionName`. It returns the value of the key as a string.

#### `func (*INIParser) Set(sectionName, key, value string)`

This method sets the value of the specified `key` within the specified `sectionName`.

#### `func (*INIParser) ToString() string`

This method converts the INIParser object back into a string representation, allowing you to retrieve the modified or parsed INI data as a string.

#### `func (*INIParser) SaveToFile(filePath string)`

This method saves the INI data from the INIParser object to the specified `filePath`.
