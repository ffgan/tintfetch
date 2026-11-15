%undefine _debugsource_packages

Name:           tintfetch
Version:        0.3.7
Release:        0.1.dev%{?dist}
Summary:        System information with colorful distribution artwork
License:        MIT AND BSD-3-Clause
URL:            https://github.com/ffgan/tintfetch
Source0:        %{url}/archive/refs/tags/v%{version}/%{name}-%{version}.tar.gz
BuildRequires:  golang >= 1.26
BuildRequires:  python3
BuildRequires:  make
BuildRequires:  gcc
BuildRequires:  glibc-devel
Suggests:       bash
Suggests:       hyfetch
Suggests:       pciutils

%description
Tintfetch displays operating system, hardware and desktop information using
a native Go engine with HyFetch color profiles and fastfetch-style modules.
Optional installed external applications provide additional rendering backends.

%prep
%autosetup

%build
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
export GOFLAGS="-buildvcs=false -buildmode=pie"
export CGO_ENABLED=1
export CGO_CFLAGS="%{build_cflags}"
export CGO_LDFLAGS="%{build_ldflags}"
%make_build VERSION=%{version}-dev GO_LDFLAGS="-B gobuildid -linkmode external -extldflags -Wl,-z,relro,-z,now"

%install
export GOFLAGS="-buildvcs=false -buildmode=pie"
export CGO_ENABLED=1
export CGO_CFLAGS="%{build_cflags}"
export CGO_LDFLAGS="%{build_ldflags}"
%make_install PREFIX=%{_prefix} BINDIR=%{_bindir} MANDIR=%{_mandir} DOCDIR=%{_docdir}/%{name} LICENSEDIR=%{_licensedir}/%{name} VERSION=%{version}-dev GO_LDFLAGS="-B gobuildid -linkmode external -extldflags -Wl,-z,relro,-z,now"

%check
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
%make_build check
sh scripts/smoke.sh

%files
%license %{_licensedir}/%{name}/LICENSE
%license %{_licensedir}/%{name}/THIRD_PARTY_NOTICES
%doc %{_docdir}/%{name}
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1*

%changelog
* Sun Nov 15 2026 ffgan <ffgandev@gmail.com> - 0.3.7-0.1.dev
- Initial development packaging, not submitted for Fedora review
