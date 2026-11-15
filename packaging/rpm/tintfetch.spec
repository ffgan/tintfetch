# Development snapshot for initial Fedora review; replace with a release tarball
# when 0.3.7 is published. The source is an immutable, public upstream commit.
%global commit 92430d863e1bc931a4dba7d3767a67e6ce89590d
%global shortcommit %(c=%{commit}; echo ${c:0:7})
%global snapdate 20261115
%undefine _debugsource_packages

Name:           tintfetch
Version:        0.3.7
Release:        0.2.%{snapdate}git%{shortcommit}%{?dist}
Summary:        System information with colorful distribution artwork
License:        GPL-2.0-or-later AND MIT AND BSD-3-Clause
URL:            https://github.com/ffgan/tintfetch
Source0:        %{url}/archive/%{commit}/%{name}-%{commit}.tar.gz
Source1:        https://raw.githubusercontent.com/ffgan/tintfetch/%{commit}/packaging/rpm/go-vendor-tools.toml
BuildRequires:  golang >= 1.26
BuildRequires:  go-rpm-macros
BuildRequires:  go-vendor-tools
BuildRequires:  python3
BuildRequires:  make
BuildRequires:  gcc
BuildRequires:  glibc-devel
ExclusiveArch:  %{golang_arches}
Suggests:       bash
Suggests:       hyfetch
Suggests:       pciutils

%description
Tintfetch displays operating system, hardware and desktop information using
a native Go engine with HyFetch color profiles and configurable information
modules. It supports plain text and JSON output. Optional installed external
applications provide additional rendering backends.

%prep
%autosetup -n %{name}-%{commit}

%generate_buildrequires
%go_vendor_license_buildrequires -c %{SOURCE1}

%build
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GO111MODULE=on
export CGO_ENABLED=1
export CGO_CFLAGS="%{build_cflags}"
export CGO_LDFLAGS="%{build_ldflags}"
export GO_LDFLAGS="-X github.com/ffgan/tintfetch/cmd.Version=%{version}-dev"
mkdir -p build
%gobuild -trimpath -o build/%{name} .

%install
install -Dpm0755 build/%{name} %{buildroot}%{_bindir}/%{name}
install -Dpm0644 docs/%{name}.1 %{buildroot}%{_mandir}/man1/%{name}.1
install -Dpm0644 docs/config.conf %{buildroot}%{_docdir}/%{name}/examples/config.conf

%check
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
%make_build check
sh scripts/smoke.sh
%go_vendor_license_check -c %{SOURCE1}

%files
%license LICENSE COPYING THIRD_PARTY_NOTICES
%doc README.md docs/compatibility.md docs/provenance.md
%doc docs/provenance.json docs/assets.sha256
%doc %{_docdir}/%{name}/examples
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1*

%changelog
* Sun Nov 15 2026 ffgan <ffgandev@gmail.com> - 0.3.7-0.2.20261115git92430d8
- Pin the review source to a public upstream snapshot
- Build with Fedora Go macros and run offline provenance and smoke tests
