Name:           otto
Version:        %{_version}
Release:        1
Summary:        The Otto Host client software
License:        Apache-2.0
Source0:        %{name}-%{version}.tar.gz
BuildRequires:  systemd-rpm-macros
Provides:       %{name} = %{version}
Prefix:         /opt

%description
Otto is an automation toolkit for Unix-like computers. This package provides the Otto client software for Otto hosts.

%global debug_package %{nil}

%prep
%autosetup

%build
cd otto
CGO_ENABLED=0 go get
cd cmd/client
CGO_ENABLED=0 go build -buildmode=exe -trimpath -ldflags="-s -w -X 'main.Version=%{version}' -X 'main.BuildDate=%{_date}'" -v -o %{name}
./%{name} -v

%install
mkdir -p %{buildroot}/opt/%{name}
install -Dpm 0755 %{name}/cmd/client/%{name} %{buildroot}/opt/%{name}/%{name}
install -Dpm 644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service

%check
cd otto
CGO_ENABLED=0 go get
CGO_ENABLED=0 go test -v ./...

%post
%systemd_post %{name}.service

%posttrans
if test pidof otto = 1; then
    systemctl restart %{name}.service
fi

%preun
%systemd_preun %{name}.service

%files
/opt/%{name}/%{name}
%{_unitdir}/%{name}.service
