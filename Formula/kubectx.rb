class Kubectx < Formula
    desc "Tool that can switch between kubectl contexts easily and create aliases"
    homepage "https://github.com/ahmetb/kubectx"
    url "https://github.com/ahmetb/kubectx/archive/v0.1.zip"
    sha256 "3d014027e38c476164638b2138f190c43fd65a22ec50035c36926555233247c0"

    bottle :unneeded

    def install
        bin.install "kubectx"
        bash_completion.install "completion/kubectx.bash" => "kubectx"
    end

    test do
        system "kubectx", "--help"
    end

    def caveats; <<-EOS.undent
      To install zsh completion, add this to your .zshrc:

        [ -f /usr/local/etc/bash_completion.d/kubectx ] && source /usr/local/etc/bash_completion.d/kubectx
      EOS
    end
end
