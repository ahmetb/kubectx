class Kubectx < Formula
  desc "Tool that can switch between kubectl contexts easily and create aliases"
  homepage "https://github.com/ahmetb/kubectx"
  url "https://github.com/ahmetb/kubectx/archive/v0.1.tar.gz"
  sha256 "841817f928af25061b1b6794400394c3e6e807e8a1c48c179f1fd8bdd553ca79"

  bottle :unneeded

  def install
    bin.install "kubectx"
    bash_completion.install "completion/kubectx.bash" => "kubectx"
  end

  def caveats; <<-EOS.undent
    To install zsh completion, add this to your .zshrc:

      [ -f /usr/local/etc/bash_completion.d/kubectx ] && source /usr/local/etc/bash_completion.d/kubectx
    EOS
  end

  test do
    system "kubectx", "--help"
  end
end
