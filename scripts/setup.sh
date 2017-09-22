#!/bin/bash

set -eu

find_mgr() {
    if hash minishift 2>/dev/null; then
        echo "minishift"
    else
        if hash docker-machine 2>/dev/null; then
            echo "docker-machine"
        fi
    fi
}

get_vm_name() {
    case "$1" in
        minishift)
            echo "minishift"
            ;;
        docker-machine)
            echo "${DOCKER_MACHINE_NAME}"
            ;;
        *)
            ;;
    esac
}

is_vm_running() {
    local vm=$1
    declare -a running=($(VBoxManage list runningvms | awk '{ print $1 }'))
    local result='false'

    for rvm in "${running[@]}"; do
        if [[ "${rvm}" == *"${vm}"* ]]; then
            result='true'
        fi
    done
    echo "$result"
}

if hash cygpath 2>/dev/null; then
    PROJECT_DIR=$(cygpath -w -a "$(pwd)")
else
    PROJECT_DIR=$(pwd)
fi

VM_MGR=$(find_mgr)
if [[ -z $VM_MGR ]]; then
    echo "ERROR: No VM Manager found; expected one of ['minishift', 'docker-machine']"
    exit 1
fi

VM_NAME=$(get_vm_name "$VM_MGR")
if [[ -z $VM_NAME ]]; then
    echo "ERROR: No VM found; try running 'eval $(docker-machine env)'"
    exit 1
fi

if ! hash VBoxManage 2>/dev/null; then
    echo "VirtualBox executable 'VBoxManage' not found in path"
    exit 1
fi

avail=$(is_vm_running "$VM_NAME")
if [[ "$avail" == *"true"* ]]; then
    res=$(VBoxManage sharedfolder add "${VM_NAME}" --name "${PROJECT}" --hostpath "${PROJECT_DIR}" --transient 2>&1)
    if [[ -z $res || $res == *"already exists"* ]]; then
        # no need to show that it already exists
        :
    else
        echo "$res"
        exit 1
    fi
    echo "VM: [${VM_NAME}] -- Added Sharedfolder [${PROJECT}] @Path [${PROJECT_DIR}]"
else
    echo "$VM_NAME is not currently running; please start your VM and try again."
    exit 1
fi

SSH_CMD="sudo mkdir -p /${PROJECT} ; sudo mount -t vboxsf ${PROJECT} /${PROJECT}"
case "${VM_MGR}" in
    minishift)
        minishift ssh "${SSH_CMD}"
        echo "VM: [${VM_NAME}] -- Mounted Sharedfolder [${PROJECT}] @VM Path [/${PROJECT}]"
        ;;
    docker-machine)
        docker-machine ssh "${VM_NAME}" "${SSH_CMD}"
        echo "VM: [${VM_NAME}] -- Mounted Sharedfolder [${PROJECT}] @VM Path [/${PROJECT}]"
        ;;
    *)
        ;;
esac
