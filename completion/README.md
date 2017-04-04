kubectx provides shell completion scripts to complete context names, making it
even faster to switch between contexts easily.

## Bash setup

Copy the `kubectx.bash` file to your HOME directory:

    cp kubectx.bash ~/.kubectx.bash

And source it in your `~/.bashrc` file by adding the line:

    [ -f ~/.kubectx.bash ] && source ~/.kubectx.bash

Start a new shell, type `kubectx`, then hit <kbd>Tab</kbd> to see the existing
contexts.

You can Add `TAB: menu-complete` to your `~/.inputrc` to cycle through the
options with <kbd>Tab</kbd>.

