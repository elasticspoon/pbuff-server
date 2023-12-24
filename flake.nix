{
  description = "Ruby on Rails development environment";

  inputs = {
    nixpkgs.url = "nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    nix-filter.url = "github:numtide/nix-filter";
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , nix-filter
    ,
    }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      overlays = [
        (self: super: {
          ruby = pkgs.ruby_3_2;
        })
      ];
      pkgs = import nixpkgs { inherit overlays system; };

      rubyEnv = pkgs.bundlerEnv {
        # The full app environment with dependencies
        name = "rails-env";
        inherit (pkgs) ruby;
        gemdir = ./.; # Points to Gemfile.lock and gemset.nix
      };
    in
    {
      apps.default = {
        type = "app";
        program = "${rubyEnv}/bin/rails";
      };

      devShells = rec {
        default = run;

        run = pkgs.mkShell {
          buildInputs = with pkgs; [ nodejs rubyEnv rubyEnv.wrappedRuby ];

          shellHook = ''
            export NIX_SHELL="true"
            ${rubyEnv}/bin/rails --version
          '';
        };
      };
    });
}
