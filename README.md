# Notes

Take notes right from your terminal.

## Installation

There are binaries in the releases, but you can also install it using Go :

```sh
go install github.com/Natouche68/notes@latest
```

To use it, just type `notes` in your terminal, and a TUI will pop up with all options.

> The first time you launch the application, it may take a while to load because of the database initialization.

## Saving

`notes` saves all data on Charm cloud, so you can sync it between multiple machines by installing the [Charm](https://github.com/charmbracelet/charm) client and running `charm link`.
