%undefine _debugsource_packages

Name:           tintfetch
Version:        0.1.0
Release:        0.1.dev%{?dist}
Summary:        System information with colorful distribution artwork
License:        MIT AND BSD-3-Clause
URL:            https://github.com/ffgan/gf
Source0:        %{url}/archive/refs/tags/v%{version}/%{name}-%{version}.tar.gz
BuildRequires:  golang >= 1.24
BuildRequires:  make
BuildRequires:  gcc
BuildRequires:  glibc-devel
Suggests:       bash
Suggests:       hyfetch
Suggests:       pciutils

%description
Tintfetch displays operating system, hardware and desktop information using
a native Go engine. Optional installed upstream engines provide additional
Neofetch and HyFetch compatibility.

%prep
%autosetup

%build
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
export GOFLAGS="-buildvcs=false -buildmode=pie"
export CGO_ENABLED=1
%make_build VERSION=%{version}-dev GO_LDFLAGS="-B gobuildid -linkmode external -extldflags -Wl,-z,relro,-z,now"

%install
export GOFLAGS="-buildvcs=false -buildmode=pie"
export CGO_ENABLED=1
%make_install PREFIX=%{_prefix} BINDIR=%{_bindir} MANDIR=%{_mandir} DOCDIR=%{_docdir}/%{name} LICENSEDIR=%{_licensedir}/%{name} VERSION=%{version}-dev GO_LDFLAGS="-B gobuildid -linkmode external -extldflags -Wl,-z,relro,-z,now"

%check
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
%make_build check

%files
%license %{_licensedir}/%{name}/LICENSE
%license %{_licensedir}/%{name}/THIRD_PARTY_NOTICES
%doc %{_docdir}/%{name}
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1*

%changelog
* Sun Nov 15 2026 ffgan <ffgan@users.noreply.github.com> - 0.1.0-0.1.dev
- Initial development packaging, not submitted for Fedora review
