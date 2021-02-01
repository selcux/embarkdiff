## EmbarkDiff

A command line sync reporting tool for files and folders.

The tool displays a list of _operations_ that _would_ be needed to sync a target directory to a reference directory.

### Usage

`$ embarkdiff [command]`

Available Commands:

- _add_: Adds the defined resource folder.

  - _--source_: Source directory path
  - _--target_: Target directory path

- _diff_: Compares given directories
- _help_: Help about any command
- _list_: Lists the given folders.

### Output

After `embarkdiff diff` command is executed, the output should be similar to the example.

```shell
delete `JustDifferent.txt`
create `new`
create `new/folder`
copy `new/folder/NewFile.txt`
delete `a/file/in/other/folder/LongChangedAtEnd.exe`
delete `a/file/in/other/folder`
delete `a/file/in/other`
copy `folder/ContentSameButLonger.txt`
```
