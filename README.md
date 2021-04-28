# Smarterzettelkasten

Litterally "*smart splip-box*" as in the slip-box method (Zettelkasten) invented by Niklas Luhman (cf. [Wikipedia](https://en.wikipedia.org/wiki/Zettelkasten)).

This slip-box note-taking method allows to organise information in a gigantic tree of small notes. Initially designed for a physical interface (a real slip-box), it has been generalized to digital entries with softwares like Obsidian.

Taking at our advantage the fact that Obsidian (and other similar software) rely on the portable Markdown format, we chose to build a tool written in Golang to manage such digital slip-box.

This allows to:

- rename a note and change all the links to this notes in all the other notes accordingly
- rename a tag in all notes
- change a prefix for all notes and directory containing notes
- detect conflicts in note ids
- detect discrepancies between files prefict and parent directory prefix

## Rules

1. Ids must contain no space
2. There must be a space between a note id and its title in the file name
3. The id of a directory should be a prefix for the ids of all its subdirectories and files

## TODO

- [x] Add function to collect all links between notes and all tags
- [ ] Add operation to rename a note and all the mentions of the note
- [ ] Add operation to rename all mentions of a tag
- [ ] Add checking for rule 3
- [ ] Add operation to change prefix
- [ ] Add history for atomic operation (renaming file, directory or tag) and add possibility to cancel said operation
