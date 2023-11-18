Name:           otto-agent
Version:        %{_version}
Release:        1
Summary:        The Otto agent
License:        Apache-2.0
Source0:        %{name}-%{version}.tar.gz
BuildRequires:  systemd-rpm-macros
Provides:       %{name} = %{version}
Prefix:         /opt
Obsoletes:      otto <= 0.11.6

%description
Otto is an automation toolkit for Unix-like computers. This package provides the Otto agent software for Otto hosts.

%global debug_package %{nil}

%prep
%autosetup

%build
cd otto/cmd/agent
CGO_ENABLED=0 GOAMD64=v2 go build -buildmode=exe -trimpath -ldflags="-s -w -X 'main.Version=%{version}' -X 'main.BuildDate=%{_date}' -X 'main.BuildRevision=%{_revision}'" -v -o agent
./agent -v

%install
mkdir -p %{buildroot}/opt/%{name}
install -Dpm 0755 otto/cmd/agent/agent %{buildroot}/opt/%{name}/agent
install -Dpm 644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service

%check
cd otto
CGO_ENABLED=0 GOAMD64=v2 go build -v ./...
CGO_ENABLED=0 GOAMD64=v2 go test -v ./...

%post
%systemd_post %{name}.service

%posttrans
if test $(readlink /proc/*/exe | grep /opt/%{name}/agent | wc -l) = 1; then
    systemctl restart %{name}.service
fi

%preun
%systemd_preun %{name}.service

%files
/opt/%{name}/agent
%{_unitdir}/%{name}.service
