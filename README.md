# Smarterzettelkasten

Litterally "*smart splip-box*" as in the slip-box method (Zettelkasten) invented by Niklas Luhman (cf. [Wikipedia](https://en.wikipedia.org/wiki/Zettelkasten)).

This slip-box note-taking method allows to organise information in a gigantic tree of small notes. Initially designed for a physical interface (a real slip-box), it has been generalized to digital entries with softwares like Obsidian.

Taking at our advantage the fact that Obsidian (and other similar software) rely on the portable Markdown format, we chose to build a tool written in Golang to manage such digital slip-box.

This allows to:

- rename a tag in all notes
- change a prefix for all notes and directory containing notes

## Usage

```
zlk change-tag <dir> <old name> <new name>
zlk change-prefix <dir> <old prefix> <new prefix>
```