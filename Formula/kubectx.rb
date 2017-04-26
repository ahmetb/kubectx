class Kubectx < Formula
  desc "Tool that can switch between kubectl contexts easily and create aliases"
  homepage "https://github.com/ahmetb/kubectx"
  url "https://github.com/ahmetb/kubectx/archive/v0.2.0.tar.gz"
  sha256 "28069aff84aaba1aa38f42d3b27e64e460a5c0651fb56b1748f44fd832d912e3"

  bottle :unneeded

  def install
    bin.install "kubectx"
    bash_completion.install "completion/kubectx.bash" => "kubectx"
    zsh_completion.install "completion/kubectx.zsh" => "_kubectx"
  end

  test do
    system "which", "kubectx"
  end
end
