{
  description = "A Nix-flake-based development environment for dependencytrack-operator";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forEachSupportedSystem =
        f:
        nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            pkgs = import nixpkgs {
              inherit system;
            };
          }
        );
    in
    {
      devShells = forEachSupportedSystem (
        { pkgs }: {
          default = pkgs.mkShell {
            packages = with pkgs; [
              coreutils
              curl
              go
              golangci-lint
              helmify
              kubectl
              kubernetes-controller-tools
              kubernetes-helm
              kustomize
              kind
              operator-sdk
              openapi-generator-cli
              podman
              setup-envtest
            ];
          };
        }
      );
    };
}
