kubectx provides shell completion scripts to complete context names, making it
even faster to switch between contexts easily.

## Bash setup

Copy the `kubectx.bash` file to your HOME directory:

```sh
cp kubectx.bash ~/.kubectx.bash
```

And source it in your `~/.bashrc` file by adding the line:

```sh
[ -f ~/.kubectx.bash ] && source ~/.kubectx.bash
```

Start a new shell, type `kubectx`, then hit <kbd>Tab</kbd> to see the existing
contexts.

You can Add `TAB: menu-complete` to your `~/.inputrc` to cycle through the
options with <kbd>Tab</kbd>.

## Zsh setup

`zsh` can leverage the `bash` completion scripts. Copy the `kubectx.bash` file
to your HOME directory:

```sh
cp kubectx.bash ~/.kubectx.bash
```

And add the following to your `.zshrc`:

```sh
[ -f ~/.kubectx.bash ] && source ~/.kubectx.bash
```

Start a new shell, type `kubectx`, then hit <kbd>Tab</kbd> to see the existing
contexts. If it does not work, modify the line above to:

```sh
[ -f ~/.kubectx.bash ] && autoload bashcompinit && bashcompinit && \
      source ~/.kubectx.bash
```
