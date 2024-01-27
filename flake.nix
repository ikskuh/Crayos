{
  description = "Crayos Global Game Jam 2024 Game";
  inputs.nixpkgs.url = "nixpkgs/nixos-23.11";

  outputs = {
    self,
    nixpkgs,
  }: let
    supportedSystems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});
  in {
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in {
      crayos-backend = pkgs.buildGoModule {
        pname = "crayos-backend";
        src = ./backend;

        # This hash locks the dependencies of this package. It is
        # necessary because of how Go requires network access to resolve
        # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
        # details. Normally one can build with a fake sha256 and rely on native Go
        # mechanisms to tell you what the hash should be or determine what
        # it should be "out-of-band" with other tooling (eg. gomod2nix).
        # To begin with it is recommended to set this, but one must
        # remember to bump this hash when your dependencies change.
        #vendorSha256 = pkgs.lib.fakeSha256;

        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };
    });

    # Add dependencies that are only needed for development
    devShells = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};

      crayos-python-packages = python-packages:
        with python-packages; [
          (
            buildPythonPackage rec {
              pname = "case-converter";
              version = "1.1.0";
              src = fetchPypi {
                inherit pname version;
                sha256 = "sha256-LtP8bj/6jWAfmjH/y8j70Z6utIZxp5qO8WOUZygkUQ4=";
              };
              doCheck = false;
              propagatedBuildInputs = [
              ];
            }
          )
        ];
      crayos-python = pkgs.python311.withPackages crayos-python-packages;
    in {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [go gopls gotools go-tools python311 crayos-python];
      };
    });

    defaultPackage = forAllSystems (system: self.packages.${system}.crayos-backend);
  };
}
