%global git_org mikefarah

# https://bugzilla.redhat.com/show_bug.cgi?id=995136#c12
%global _dwz_low_mem_die_limit 0

Name:           yq
Version:        2.3.0
Release:        1%{?dist}
Summary:        Process YAML documents from the CLI

License:        MIT
URL:            https://github.com/%{git_org}/%{name}
Source0:        https://github.com/%{git_org}/%{name}/archive/%{version}.tar.gz

BuildRequires:  golang rsync

%description
yq is a lightweight and portable command-line YAML processor.

The aim of the project is to be the jq or sed of yaml files.



%prep
%autosetup -n %{name}-%{version}

%build

mkdir -p gopath/src/github.com/%{git_org}/%{name}
export GOPATH=${PWD}/gopath
export PATH=${GOPATH}:${PATH}
rsync -az --exclude=gopath/ ./ gopath/src/github.com/%{git_org}/%{name}
cd gopath/src/github.com/%{git_org}/%{name}

# Get go dependencies
go get -u github.com/kardianos/govendor
go install github.com/kardianos/govendor

# Build the yq binary
${GOPATH}/bin/govendor sync
go build -o %{name} -ldflags=-linkmode=external

%install
rm -rf %{buildroot}

# Install binaries & scripts
install -d %{buildroot}%{_bindir}
install -p -m 755 gopath/src/github.com/%{git_org}/%{name}/%{name} %{buildroot}%{_bindir}

%files
%doc README.md LICENSE

%{_bindir}/%{name}

%changelog
* Mon Mar 25 2019 Bradford Dabbs <brad@perched.io>
 - Initial creation of spec file
