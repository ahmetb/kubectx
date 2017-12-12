This repository provides both `kubectx` and `kubens` tools.


**`kubectx`** help you switch between clusters back and forth:
![kubectx demo GIF](img/kubectx-demo.gif)

**`kubens`** help you switch between Kubernetes namespaces smoothly:
![kubens demo GIF](img/kubens-demo.gif)

# kubectx(1)

kubectx is an utility to manage and switch between kubectl(1) contexts.

```
USAGE:
  kubectx                   : list the contexts
  kubectx <NAME>            : switch to context <NAME>
  kubectx -                 : switch to the previous context
  kubectx <NEW_NAME>=<NAME> : rename context <NAME> to <NEW_NAME>
  kubectx -h,--help         : show this message
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

`kubectx` supports <kbd>Tab</kbd> completion on bash/zsh shells to help with 
long context names. You don't have to remember full context names anymore.

-----

# kubens(1)

kubens is an utility to switch between Kubernetes namespaces.

```
USAGE:
  kubens                    : list the namespaces
  kubens <NAME>             : change the active namespace
  kubens -                  : switch to the previous namespace
  kubens -h,--help          : show this message
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

`kubens` also supports <kbd>Tab</kbd> completion on bash/zsh shells.

-----

## Installation

**For macOS:**

:tada: kubectx is now in Homebrew! :confetti_ball:

> Use the [Homebrew](https://brew.sh/) package manager:
>
>     brew install kubectx
>
> this will also set up bash/zsh completion scripts automatically.

Running `brew install` with `--with-short-names` will install tools with names
`kctx` and `kns` to prevent prefix collision with `kubectl` name.

> Note: If you installed kubectx before it was accepted to Homebrew core
> repository, reinstall with:
> `brew untap ahmetb/kubectx && brew uninstall --force kubectx && brew update && brew install kubectx`

**Other platforms:**

Since `kubectx`/`kubens` are written in Bash, they can run in shells that support POSIX standards.

- Download the `kubectx`, `kubens` and `utils.bash` scripts
- Either:
  - save them all to soemwhere in your `PATH`,
  - or save them to a directory, then create symlinks to `kubectx`/`kubens` from somewhere in your `PATH`, like `/usr/local/bin`
- Make `kubectx` and `kubens` executable (`chmod +x ...`)
- You’re on your own to install bash/zsh [completion scripts](completion/) manually.

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
| _“Special thanks to Ahmet Alp Balkan for creating kubectx, kubens, and kubectl aliases, as these tools made my life better.” – [@strebeld](https://medium.com/@strebeld/5-ways-to-enhance-kubectl-ux-97c8893227a)
-----

Disclaimer: This is not an official Google product.


#### Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/ahmetb/kubectx.svg)](https://starcharts.herokuapp.com/ahmetb/kubectx)

