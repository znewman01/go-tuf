{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:tweag/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    nix-filter.url = "github:numtide/nix-filter";
  };
  outputs = { nixpkgs, flake-utils, gomod2nix, nix-filter, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlay ];
        };
        python = pkgs.python39;
        securesystemslib = python.pkgs.buildPythonPackage rec {
          pname = "securesystemslib";
          version = "0.11.2";
          src = python.pkgs.fetchPypi {
            inherit pname version;
            sha256 = "Q1VDcf7u9QGWWHqgZs/9a5zv9rSE+nsSfhOfr7XA4j4=";
          };
          propagatedBuildInputs = with python.pkgs; [
            six
            cryptography
            pynacl
            colorama
          ];
          doCheck = false;
          prePatch = ''
            substituteInPlace ./setup.py --replace 'six>=0.11.1' 'six'
            substituteInPlace ./setup.py --replace 'cryptography>=2.2.2' 'cryptography'
            substituteInPlace ./setup.py --replace 'pynacl>=1.2.0' 'pynacl'
            substituteInPlace ./setup.py --replace 'colorama>=0.3.9' 'colorama'
          '';
        };
        python-tuf = python.pkgs.buildPythonPackage rec {
          pname = "tuf";
          version = "0.11.1";
          src = python.pkgs.fetchPypi {
            inherit pname version;
            sha256 = "ZdX4ekGDBJS/WF+KUIJhirJgFdFWpn8j43VSQZ5CfPE=";
          };
          doCheck = false;
          propagatedBuildInputs =
            [ python.pkgs.six python.pkgs.iso8601 securesystemslib ];
          prePatch = ''
            substituteInPlace ./setup.py --replace 'iso8601>=0.1.12' 'iso8601'
            substituteInPlace ./setup.py --replace 'six>=0.11.1' 'six'
            substituteInPlace ./setup.py --replace 'securesystemslib>0.11.2' 'securesystemslib'
          '';
        };
      in rec {
        packages.default = pkgs.buildGoApplication {
          pname = "go-tuf";
          version = "v0.0.0-deadbeef";
          src = nix-filter.lib.filter {
            root = ./.;
            exclude = [ (nix-filter.lib.matchExt "nix") ];
          };
          modules = ./gomod2nix.toml;
          checkInputs = [ python python-tuf ];
        };
        apps = {
          tuf = flake-utils.lib.mkApp {
            name = "tuf";
            drv = packages.default;
          };
          tuf-client = flake-utils.lib.mkApp {
            name = "tuf-client";
            drv = packages.default;
          };
        };
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs;
            [
              nixfmt
              gocode
              gore
              gomodifytags
              gopls
              go-symbols
              gopkgs
              go-outline
              gotests
              gotools
              golangci-lint
              gomod2nix
              jq
            ] ++ packages.default.nativeBuildInputs;
        };
      });
}
