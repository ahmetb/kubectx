This repository provides both `kubectx` and `kubens` tools. Purpose of this
project is to provide an utility and facilitate discussion about how `kubectl`
can manage contexts better.

# kubectx(1)

kubectx is an utility to manage and switch between kubectl(1) contexts.

```
USAGE:
  kubectx                   : list the contexts
  kubectx <NAME>            : switch to context
  kubectx -                 : switch to the previous context
  kubectx <NEW_NAME>=<NAME> : create alias for context
  kubectx -h,--help         : show this message
```

Purpose of this project is to provide an utility and facilitate discussion
about how `kubectl` can manage contexts better.

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

`kubectx` also supports <kbd>Tab</kbd> completion, which helps with long context
names.

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

`kubectx` also supports <kbd>Tab</kbd> completion, which helps with long context
names.

-----

## Installation

For macOS:

> Use [Homebrew](https://brew.sh/) package manager:
>
>      brew tap ahmetb/kubectx https://github.com/ahmetb/kubectx.git
>      brew install kubectx
> this will also set up bash/zsh completion scripts automatically.

Other platforms:

> Download the `kubectx` script, make it executable and add it to your PATH. You
> can also install bash/zsh [completion scripts](completion/) manually.

-----

Disclaimer: This is not an official Google product.
