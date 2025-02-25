# git-auto-commit

![GoDoc](https://pkg.go.dev/badge/github.com/ivy/git-auto-commit)
![Build](https://github.com/ivy/git-auto-commit/actions/workflows/build.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ivy/git-auto-commit)
![License](https://img.shields.io/github/license/ivy/git-auto-commit)

> 🚧 **Work in Progress**
> 
> `git-auto-commit` is actively being developed. Expect frequent updates and improvements.

`git-auto-commit` ✨ is an open-source tool written in Go that leverages a large-language model (LLM) to automatically draft commit and pull request messages. This allows developers to stay focused on the *why* and *how* of software development while maintaining high-quality documentation in commit histories. By automating commit message generation, `git-auto-commit` helps save time, streamline development workflows, and make changes easier to review.

## 🚀 Installation

To install `git-auto-commit`, run the following command:

```sh
curl -fsSL https://get.git-auto-commit.org | sh
```

## 📖 Usage

### ✍️ git auto-commit  

`git auto-commit` analyzes your staged changes and generates a clear, contextual commit message using an LLM.  

#### Options:  
- **`-v, --verbose`** _(default)_ – Opens your `$EDITOR` (or falls back to `nano` or `vi`) with a suggested commit message. Edit and save to finalize the commit.  
- **`-y, --yes`** – Commits your changes with the suggested message without prompting.  
- **`-m MSG, --message MSG`** – Adds extra context to the LLM, useful for explaining _why_ the change was made.  
- **`-M MODEL, --model MODEL`** – Overrides the default model used for message generation.  
- **`-p PROVIDER, --provider PROVIDER`** – Overrides the default LLM provider.  

Additional arguments can be passed to `git commit`:

```sh
git auto-commit -m "My message" -- --amend
```

### 🔀 git auto-pr

`git auto-pr` automates PR descriptions using AI, reducing manual effort and ensuring well-structured messages. Requires the [GitHub CLI (`gh`)](https://cli.github.com/).  

#### Usage
All options from `git auto-commit` apply, with the ability to pass additional arguments to `gh`.  

```sh
git auto-pr -m "My message" -- --open
```  

This example adds a custom message and opens the PR immediately after creation.

## 📌 Roadmap

`git-auto-commit` is under active development, and several features are planned for future releases:

- **Customization & Configuration**
  - Users can configure tool behavior through `git-config` (#20)
  - Specify LLM model with `--model`/`-M` flag (#4)
  - Supply context manually via `--message`/`-m` option (#2)

- **Commit Message Generation**
  - Automatically draft commit messages (#1)
  - Derive commit style from:
    - Built-in templates (#18)
    - User preferences (#16)
    - Past commit messages (#10)
  - Split large commits into smaller, focused commits (#13)

- **Pull Request Generation**
  - Draft pull request messages from `git log -p` (#3)
  - Derive pull request style from:
    - Built-in templates (#19)
    - User preferences (#17)
    - Repository templates (#15)
    - Past merged pull requests (#11)

- **Enhanced Functionality**
  - Stream LLM output to the terminal before opening `$EDITOR` (#12)
  - Offer to update the changelog if one is present (#9)
  - Users can provide context via voice (speech-to-text) (#14)

- **Packaging & Distribution**
  - Publish Homebrew formula (#7)
  - Publish Debian/Ubuntu packages (#6)
  - Support easy install through `curl | bash` (#5)
  - Include Go/Devcontainer boilerplate (#8)

## 🛠 Development

To develop locally, clone the repository and use the provided scripts:

```sh
bin/setup    # Set up the development environment
bin/build    # Build the project
bin/test     # Run the test suite
```

Additionally, the repository includes a [Devcontainer](https://code.visualstudio.com/docs/devcontainers/containers) configuration compatible with IDEs such as VS Code and Cursor, enabling a seamless development experience.

## 🤝 Contributing

Contributions are welcome! 🎉 If you’d like to contribute, please check the open issues, fork the repository, and submit a pull request. Discussions, feedback, and suggestions are encouraged in the issue tracker.

## 📜 License

This project is open-source and licensed under the [ISC License](LICENSE.md).
