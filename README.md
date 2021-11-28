# `kubectx` + `kubens`: Power tools for kubectl

![Latest GitHub release](https://img.shields.io/github/release/ahmetb/kubectx.svg)
![GitHub stars](https://img.shields.io/github/stars/ahmetb/kubectx.svg?label=github%20stars)
![Homebrew downloads](https://img.shields.io/homebrew/installs/dy/kubectx?label=macOS%20installs)
[![Go implementation (CI)](https://github.com/ahmetb/kubectx/workflows/Go%20implementation%20(CI)/badge.svg)](https://github.com/ahmetb/kubectx/actions?query=workflow%3A"Go+implementation+(CI)")
![Proudly written in Bash](https://img.shields.io/badge/written%20in-bash-ff69b4.svg)

This repository provides both `kubectx` and `kubens` tools.
[Install &rarr;](#installation)

## What are `kubectx` and `kubens`?

**kubectx** is a tool to switch between contexts (clusters) on kubectl
faster.<br/>
**kubens** is a tool to switch between Kubernetes namespaces (and
configure them for kubectl) easily.

Here's a **`kubectx`** demo:
![kubectx demo GIF](img/kubectx-demo.gif)

...and here's a **`kubens`** demo:
![kubens demo GIF](img/kubens-demo.gif)

### Examples

```sh
# switch to anoter cluster that's in kubeconfig
$ kubectx minikube
Switched to context "minikube".

# switch back to previous cluster
$ kubectx -
Switched to context "oregon".

# create an alias for the context
$ kubectx dublin=gke_ahmetb_europe-west1-b_dublin
Context "dublin" set.
Aliased "gke_ahmetb_europe-west1-b_dublin" as "dublin".

# change the active namespace on kubectl
$ kubens kube-system
Context "test" set.
Active namespace is "kube-system".

# go back to the previous namespace
$ kubens -
Context "test" set.
Active namespace is "default".
```

If you have [`fzf`](https://github.com/junegunn/fzf) installed, you can also
**interactively** select a context or cluster, or fuzzy-search by typing a few
characters. To learn more, read [interactive mode &rarr;](#interactive-mode)

Both `kubectx` and `kubens` support <kbd>Tab</kbd> completion on bash/zsh/fish
shells to help with long context names. You don't have to remember full context
names anymore.

-----

## Installation

Stable versions of `kubectx` and `kubens` are small bash scripts that you
can find in this repository.

Starting with v0.9.0, `kubectx` and `kubens` **are now rewritten in Go**.  They
should work the same way (and we'll keep the bash-based implementations around)
but the new features will be added to the new Go programs.  Please help us test
this new Go implementation by downloading the binaries from the [**Releases page
&rarr;**](https://github.com/ahmetb/kubectx/releases)

**Installation options:**

- [as kubectl plugins (macOS & Linux)](#kubectl-plugins-macos-and-linux)
- [with Homebrew (macOS & Linux)](#homebrew-macos-and-linux)
- [with MacPorts (macOS)](#macports-macos)
- [with apt (Debian)](#apt-debian)
- [with pacman (Arch Linux)](#pacman-arch-linux)
- [manually (macOS & Linux)](#manual-installation-macos-and-linux)

If you like to add context/namespace information to your shell prompt (`$PS1`),
you can try out [kube-ps1].

[kube-ps1]: https://github.com/jonmosco/kube-ps1

### Kubectl Plugins (macOS and Linux)

You can install and use the [Krew](https://github.com/kubernetes-sigs/krew/) kubectl
plugin manager to get `kubectx` and `kubens`.

**Note:** This will not install the shell completion scripts. If you want them,
*choose another installation method
or install the scripts [manually](#manual-installation-macos-and-linux).

```sh
kubectl krew install ctx
kubectl krew install ns
```

After installing, the tools will be available as `kubectl ctx` and `kubectl ns`.

### Homebrew (macOS and Linux)

If you use [Homebrew](https://brew.sh/) you can install like this:

```sh
brew install kubectx
```

This command will set up bash/zsh/fish completion scripts automatically.


### MacPorts (macOS)

If you use [MacPorts](https://www.macports.org) you can install like this:

```sh
sudo port install kubectx
```

### apt (Debian)

``` bash
sudo apt install kubectx
```
Newer versions might be available on repos like
[Debian Buster (testing)](https://packages.debian.org/buster/kubectx),
[Sid (unstable)](https://packages.debian.org/sid/kubectx)
(_if you are unfamiliar with the Debian release process and how to enable
testing/unstable repos, check out the
[Debian Wiki](https://wiki.debian.org/DebianReleases)_):


### pacman (Arch Linux)

Available as official Arch Linux package. Install it via:

```bash
sudo pacman -S kubectx
```


### Manual Installation (macOS and Linux)

Since `kubectx` and `kubens` are written in Bash, you should be able to install
them to any POSIX environment that has Bash installed.

- Download the `kubectx`, and `kubens` scripts.
- Either:
  - save them all to somewhere in your `PATH`,
  - or save them to a directory, then create symlinks to `kubectx`/`kubens` from
    somewhere in your `PATH`, like `/usr/local/bin`
- Make `kubectx` and `kubens` executable (`chmod +x ...`)

Example installation steps:

``` bash
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```

If you also want to have shell completions, pick an installation method for the
[completion scripts](completion/) that fits your system best: [`zsh` with
`antibody`](#completion-scripts-for-zsh-with-antibody), [plain
`zsh`](#completion-scripts-for-plain-zsh),
[`bash`](#completion-scripts-for-bash) or
[`fish`](#completion-scripts-for-fish).

#### Completion scripts for `zsh` with [antibody](https://getantibody.github.io)

Add this line to your [Plugins File](https://getantibody.github.io/usage/) (e.g.
`~/.zsh_plugins.txt`):

```
ahmetb/kubectx path:completion kind:fpath
```

Depending on your setup, you might or might not need to call `compinit` or
`autoload -U compinit && compinit` in your `~/.zshrc` after you load the Plugins
file. If you use [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh), load the
completions before you load `oh-my-zsh` because `oh-my-zsh` will call
`compinit`. 

#### Completion scripts for plain `zsh`

The completion scripts have to be in a path that belongs to `$fpath`. Either
link or copy them to an existing folder.

Example with [`oh-my-zsh`](https://github.com/ohmyzsh/ohmyzsh):

```bash
mkdir -p ~/.oh-my-zsh/completions
chmod -R 755 ~/.oh-my-zsh/completions
ln -s /opt/kubectx/completion/_kubectx.zsh ~/.oh-my-zsh/completions/_kubectx.zsh
ln -s /opt/kubectx/completion/_kubens.zsh ~/.oh-my-zsh/completions/_kubens.zsh
```

If completion doesn't work, add `autoload -U compinit && compinit` to your
`.zshrc` (similar to
[`zsh-completions`](https://github.com/zsh-users/zsh-completions/blob/master/README.md#oh-my-zsh)).

If you are not using [`oh-my-zsh`](https://github.com/ohmyzsh/ohmyzsh), you
could link to `/usr/share/zsh/functions/Completion` (might require sudo),
depending on the `$fpath` of your zsh installation.

In case of errors, calling `compaudit` might help.

#### Completion scripts for `bash`

```bash
git clone https://github.com/ahmetb/kubectx.git ~/.kubectx
COMPDIR=$(pkg-config --variable=completionsdir bash-completion)
ln -sf ~/.kubectx/completion/kubens.bash $COMPDIR/kubens
ln -sf ~/.kubectx/completion/kubectx.bash $COMPDIR/kubectx
cat << EOF >> ~/.bashrc


#kubectx and kubens
export PATH=~/.kubectx:\$PATH
EOF
```

#### Completion scripts for `fish`

```fish
mkdir -p ~/.config/fish/completions
ln -s /opt/kubectx/completion/kubectx.fish ~/.config/fish/completions/
ln -s /opt/kubectx/completion/kubens.fish ~/.config/fish/completions/
```

-----

### Interactive mode

If you want `kubectx` and `kubens` commands to present you an interactive menu
with fuzzy searching, you just need to [install
`fzf`](https://github.com/junegunn/fzf) in your `$PATH`.

![kubectx interactive search with fzf](img/kubectx-interactive.gif)

If you have `fzf` installed, but want to opt out of using this feature, set the
environment variable `KUBECTX_IGNORE_FZF=1`.

If you want to keep `fzf` interactive mode but need the default behavior of the
command, you can do it by piping the output to another command (e.g. `kubectx |
cat `).

-----

### Customizing colors

If you like to customize the colors indicating the current namespace or context,
set the environment variables `KUBECTX_CURRENT_FGCOLOR` and
`KUBECTX_CURRENT_BGCOLOR` (refer color codes
[here](https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/)):

```sh
export KUBECTX_CURRENT_FGCOLOR=$(tput setaf 6) # blue text
export KUBECTX_CURRENT_BGCOLOR=$(tput setab 7) # white background
```

Colors in the output can be disabled by setting the
[`NO_COLOR`](https://no-color.org/) environment variable.

-----

If you liked `kubectx`, you may like my
[`kubectl-aliases`](https://github.com/ahmetb/kubectl-aliases) project, too. I
recommend pairing kubectx and kubens with [fzf](#interactive-mode) and
[kube-ps1].

#### Stargazers over time

[![Stargazers over time](https://starchart.cc/ahmetb/kubectx.svg)](https://starchart.cc/ahmetb/kubectx)
![Google Analytics](https://ga-beacon.appspot.com/UA-2609286-17/kubectx/README?pixel) <!-- TODO broken since Aug 2021 as igrigorik left Google -->
