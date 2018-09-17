This repository provides both `kubectx` and `kubens` tools.


**`kubectx`** helps you switch between clusters back and forth:
![kubectx demo GIF](img/kubectx-demo.gif)

**`kubens`** helps you switch between Kubernetes namespaces smoothly:
![kubens demo GIF](img/kubens-demo.gif)

# kubectx(1)

kubectx is an utility to manage and switch between kubectl(1) contexts.

```
USAGE:
  kubectx                   : list the contexts
  kubectx <NAME>            : switch to context <NAME>
  kubectx -                 : switch to the previous context
  kubectx <NEW_NAME>=<NAME> : rename context <NAME> to <NEW_NAME>
  kubectx <NEW_NAME>=.      : rename current-context to <NEW_NAME>
  kubectx -d <NAME>         : delete context <NAME> ('.' for current-context)
                              (this command won't delete the user/cluster entry
                              that is used by the context)
```

### Usage

```sh
$ kubectx minikube
Switched to context "minikube".

$ kubectx -
Switched to context "oregon".

$ kubectx -
Switched to context "minikube".

$ kubectx dublin=gke_ahmetb_europe-west1-b_dublin
Context "dublin" set.
Aliased "gke_ahmetb_europe-west1-b_dublin" as "dublin".
```

`kubectx` supports <kbd>Tab</kbd> completion on bash/zsh/fish shells to help with
long context names. You don't have to remember full context names anymore.

-----

# kubens(1)

kubens is an utility to switch between Kubernetes namespaces.

```
USAGE:
  kubens                    : list the namespaces
  kubens <NAME>             : change the active namespace
  kubens -                  : switch to the previous namespace
```


### Usage

```sh
$ kubens kube-system
Context "test" set.
Active namespace is "kube-system".

$ kubens -
Context "test" set.
Active namespace is "default".
```

`kubens` also supports <kbd>Tab</kbd> completion on bash/zsh/fish shells.

-----

## Installation

### macOS

:confetti_ball: Use the [Homebrew](https://brew.sh/) package manager:

    brew install kubectx

This command will set up bash/zsh/fish completion scripts automatically.


- Running `brew install` with `--with-short-names` will install tools with names
`kctx` and `kns` to prevent prefix collision with `kubectl` name.

- If you like to add context/namespace info to your shell prompt (`$PS1`),
  I recommend trying out [kube-ps1](https://github.com/jonmosco/kube-ps1).

### Linux

Since `kubectx`/`kubens` are written in Bash, you should be able to instal
them to any POSIX environment that has Bash installed.

- Download the `kubectx`, and `kubens` scripts.
- Either:
  - save them all to somewhere in your `PATH`,
  - or save them to a directory, then create symlinks to `kubectx`/`kubens` from
    somewhere in your `PATH`, like `/usr/local/bin`
- Make `kubectx` and `kubens` executable (`chmod +x ...`)
- Install bash/zsh/fish [completion scripts](completion/).  
  - For zsh:  
    The completion scripts have to be in a path that belongs to `$fpath`. Either link or copy them to an existing folder.  
    If using oh-my-zsh you can do as follows:
    ```bash
    mkdir -p ~/.oh-my-zsh/completions
    chmod -R 755 ~/.oh-my-zsh/completions
    ln -s /opt/kubectx/completion/kubectx.zsh ~/.oh-my-zsh/completions/_kubectx.zsh
    ln -s /opt/kubectx/completion/kubens.zsh ~/.oh-my-zsh/completions/_kubens.zsh
    ```  
    Note that the leading underscore seems to be a convention.  
    If not using oh-my-zsh, you could link to `/usr/share/zsh/functions/Completion` (might require sudo), depending on the `$fpath` of your zsh installation.  
    In case of error, calling `compaudit` might help.
  - For bash:
    ```bash
    git clone https://github.com/ahmetb/kubectx.git ~/.kubectx
    COMPDIR=$(pkg-config --variable=completionsdir bash-completion)
    ln -sf ~/.kubectx/completion/kubens.bash $COMPDIR/kubens
    ln -sf ~/.kubectx/completion/kubectx.bash $COMPDIR/kubectx
    cat << FOE >> ~/.bashrc
    
    
    #kubectx and kubens
    export PATH=~/.kubectx:\$PATH
    FOE
    ```
  - For fish: Figure out how to install completion scripts and please document here
  
Example installation steps:

``` bash
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```

#### Arch Linux

An unofficial [AUR package](https://aur.archlinux.org/packages/kubectx) `kubectx`
is available. Install instructions can be found on the [Arch 
wiki](https://wiki.archlinux.org/index.php/Arch_User_Repository#Installing_packages).

#### Debian/Ubuntu

Available as a Debian package for [Debian Buster (testing)](https://packages.debian.org/buster/kubectx), [Sid (unstable)](https://packages.debian.org/sid/kubectx) (_note: if you are unfamiliar with Debian release process and how to enable testing/unstable repos, check the [Debian Wiki](https://wiki.debian.org/DebianReleases)_):

``` bash
sudo apt install kubectx
```

-----

### Interactive mode

If you want `kubectx` and `kubens` commands to present you an interactive menu
with fuzzy searching, you just need to [install
`fzf`](https://github.com/junegunn/fzf) in your PATH.

![kubectx interactive search with fzf](img/kubectx-interactive.gif)

If you have `fzf` installed, but want to opt out of using this feature, set the environment variable `KUBECTX_IGNORE_FZF=1`.


-----

### Customizing colors

If you like to customize the colors indicating the current namespace or context, set the environment variables `KUBECTX_CURRENT_FGCOLOR` and `KUBECTX_CURRENT_BGCOLOR` (refer color codes [here](https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/)):

```
export KUBECTX_CURRENT_FGCOLOR=$(tput setaf 6) # blue text
export KUBECTX_CURRENT_BGCOLOR=$(tput setaf 7) # white background
```

Colors in the output can be disabled by setting the
[`NO_COLOR`](http://no-color.org/) environment variable.

-----

####  Users

| What are others saying about kubectx? |
| ---- |
| _“Thank you for kubectx & kubens - I use them all the time & have them in my k8s toolset to maintain happiness :) ”_ – [@pbouwer](https://twitter.com/pbouwer/status/925896377929949184) |
| _“I can't imagine working without kubectx and especially kubens anymore. It's pure gold.”_ – [@timoreimann](https://twitter.com/timoreimann/status/925801946757419008) |
| _“I'm liking kubectx from @ahmetb, makes it super-easy to switch #Kubernetes contexts [...]”_ &mdash; [@lizrice](https://twitter.com/lizrice/status/928556415517589505) |
| _“Also using it on a daily basis. This and my zsh config that shows me the current k8s context 😉”_ – [@puja108](https://twitter.com/puja108/status/928742521139810305) |
| _“Lately I've found myself using the kubens command more than kubectx. Both very useful though :-)”_ – [@stuartleeks](https://twitter.com/stuartleeks/status/928562850464907264) |
| _“yeah kubens rocks!”_ – [@embano1](https://twitter.com/embano1/status/928698440732815360) |
| _“Special thanks to Ahmet Alp Balkan for creating kubectx, kubens, and kubectl aliases, as these tools made my life better.”_ – [@strebeld](https://medium.com/@strebeld/5-ways-to-enhance-kubectl-ux-97c8893227a)

> If you liked `kubectx`, you may like my [`kubectl-aliases`](https://github.com/ahmetb/kubectl-aliases) project, too.

-----

Disclaimer: This is not an official Google product.


#### Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/ahmetb/kubectx.svg)](https://starcharts.herokuapp.com/ahmetb/kubectx)

