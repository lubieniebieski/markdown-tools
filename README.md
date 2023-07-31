# markdown-tools: A collection of markdown-related tools packed in a CLI

The markdown-tools repository is a collection of markdown-related tools packed in a CLI. It provides a script that can be used to cleanup `*.md` files and use reference-style links instead of inline urls.

## Example

### Input

```markdown
# My markdown file

I have an inline link to [Google](https://www.google.com) and [Facebook](https://www.facebook.com).

Last line of the file.
```

### Output

```markdown
# My markdown file

I have an inline link to [Google][1] and [Facebook][2].

Last line of the file.

[1]: https://www.google.com
[2]: https://www.facebook.com
```

## Usage

### Simple version

You can install it via brew -- repo contents available under [lubieniebieski/homebrew-tools](https://github.com/lubieniebieski/homebrew-tools)

```bash
brew install lubieniebieski/tools/markdown-tools
markdown-tools links_as_references <PATH>
```

## Known issues and potential improvements

- not everything is right if you run the script multiple times on the same content
- path handling should be better - if you're passing a directory, it will parse all files in the directory recursively
- it would be nice to have a flag to specify the output directory
- it wouldn't hurt to make sure that someone is really expecting changes in the existing files

## Feedback

If you have any feedback, please reach out to me on Twitter or email - you can find the links on my GitHub profile [@lubieniebieski](https://https://github.com/lubieniebieski/). Probably a lot of things could be done better, so I'm open to suggestions.