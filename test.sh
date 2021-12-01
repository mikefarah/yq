#!/bin/bash
# load array into a bash array
# need to output each entry as a single line
# readarray identityMappings < <(./yq e -o=j -I=0 '.identitymappings[]' test.yml )

# for identityMapping in "${identityMappings[@]}"; do
#     # identity mapping is a yaml snippet representing a single entry
#     roleArn=$(echo "$identityMapping" | yq e '.arn' -)
#     echo "roleArn: $roleArn"
# done




while IFS=$'\t' read -r roleArn group user _; do
  echo "Role:  $roleArn"
  echo "Group: $group"
  echo "User:  $user"
done < <(yq -j read test.yaml \
         | jq -r '.identitymappings[] | [.arn, .group, .user] | @tsv')