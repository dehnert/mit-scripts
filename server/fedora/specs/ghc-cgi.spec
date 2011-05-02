# For Haskell Packaging Guidelines see:
# - https://fedoraproject.org/wiki/Packaging:Haskell
# - https://fedoraproject.org/wiki/PackagingDrafts/Haskell

%global pkg_name cgi

# common part of summary for all the subpackages
%global common_summary Haskell %{pkg_name} library

# main description used for all the subpackages
%global common_description A %{pkg_name} library for Haskell.

# Haskell library dependencies (used for buildrequires and devel/prof subpkg requires)
%global ghc_pkg_deps ghc-network-devel, ghc-parsec-devel, ghc-mtl-devel, ghc-MonadCatchIO-mtl-devel, ghc-xhtml-devel

# foreign library dependencies (used for buildrequires and devel subpkg requires)
#%%global ghc_pkg_c_deps @CDEP1@-devel

Name:           ghc-%{pkg_name}
Version:        3001.1.8.2
Release:        0.%{scriptsversion}%{?dist}
Summary:        %{common_summary}

Group:          System Environment/Libraries
License:        BSD
URL:            http://hackage.haskell.org/package/%{pkg_name}
Source0:        http://hackage.haskell.org/packages/archive/%{pkg_name}/%{version}/%{pkg_name}-%{version}.tar.gz
# fedora ghc archs:
ExclusiveArch:  %{ix86} x86_64 ppc alpha sparcv9
BuildRequires:  ghc, ghc-doc, ghc-prof
# macros for building haskell packages
BuildRequires:  ghc-rpm-macros >= 0.7.3
BuildRequires:  hscolour
%{?ghc_pkg_deps:BuildRequires:  %{ghc_pkg_deps}, %(echo %{ghc_pkg_deps} | sed -e "s/\(ghc-[^, ]\+\)-devel/\1-doc,\1-prof/g")}
%{?ghc_pkg_c_deps:BuildRequires:  %{ghc_pkg_c_deps}}

%description
%{common_description}


%prep
%setup -q -n %{pkg_name}-%{version}


%build
%ghc_lib_build


%install
%ghc_lib_install


# define the devel and prof subpkgs, devel post[un] scripts, and filelists:
# ghc-%pkg_name{,devel,prof}.files
%ghc_lib_package


%changelog
* Mon May  2 2011 Alexander Chernyakhovsky <achernya@mit.edu> - 3001.1.8.2-0
- regenerated packaging with cabal2spec-0.22.5

* Thu Sep  9 2010 Anders Kaseorg <andersk@mit.edu> - 3001.1.8.1-0
- initial packaging for Fedora automatically generated by cabal2spec-0.22.1
