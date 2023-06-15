FROM public.ecr.aws/i1i0w6p5/mikefarah/yq:4

COPY entrypoint.sh /entrypoint.sh

# this seems to be the default user in github actions
USER 1001

ENTRYPOINT ["/entrypoint.sh"]
