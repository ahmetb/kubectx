class Kubectx < Formula
  desc "Tool that can switch between kubectl contexts easily and create aliases"
  homepage "https://github.com/ahmetb/kubectx"
  url "https://github.com/ahmetb/kubectx/archive/v0.2.0.tar.gz"
  sha256 "28069aff84aaba1aa38f42d3b27e64e460a5c0651fb56b1748f44fd832d912e3"
  head "https://github.com/ahmetb/kubectx.git", :branch => "short-names"
  bottle :unneeded

  option "with-short-names", "link as \"kctx\" and \"kns\" instead"

  def install
    bin.install "kubectx" => build.with?("short-names") ? "kctx" : "kubectx"
    bin.install "kubens" => build.with?("short-names") ? "kns" : "kubens"
    include.install "utils.bash"

    bash_completion.install "completion/kubectx.bash" => "kubectx"
    bash_completion.install "completion/kubens.bash" => "kubens"
    zsh_completion.install "completion/kubectx.zsh" => "_kubectx"
    zsh_completion.install "completion/kubens.zsh" => "_kubens"
  end

  test do
    system "which", build.with?("short-names") ? "kctx" : "kubectx"
    system "which", build.with?("short-names") ? "kns" : "kubens"
  end
end
