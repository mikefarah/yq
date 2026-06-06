FROM mikefarah/yq:4@sha256:11a1f0b604b13dbbdc662260d8db6f644b22d8553122a25c1b5b2e8713ca6977

COPY entrypoint.sh /entrypoint.sh

# github action recommendation is to run as root.
# https://docs.github.com/en/actions/creating-actions/dockerfile-support-for-github-actions#user
USER root

ENTRYPOINT ["/entrypoint.sh"]
