#!/usr/bin/env bash
PACKAGE=api
echo "Starting build for ${PACKAGE}"
set -e
if [[ -z ${PACKAGE} ]]; then
    echo "need to supply which bot to build"
    exit 1
fi

ctr=$(buildah from fedora)
mnt=$(buildah mount $ctr)

cp build/${PACKAGE}/${PACKAGE} $mnt/${PACKAGE}
chmod +x $mnt/${PACKAGE}

buildah config --entrypoint /${PACKAGE} $ctr
buildah umount $ctr
buildah commit $ctr almagest-${PACKAGE}
echo "Commiting to <almagest-${PACKAGE}>"
