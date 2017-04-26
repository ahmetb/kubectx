class Kubectx < Formula
  desc "Tool that can switch between kubectl contexts easily and create aliases"
  homepage "https://github.com/ahmetb/kubectx"
  url "https://github.com/ahmetb/kubectx/archive/v0.2.0.tar.gz"
  sha256 "9fb6557416e4be3ef7e9701527f89fa73a1a0545a3cafbe6bed7527061c6cfb7"

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
