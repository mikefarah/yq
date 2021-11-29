#!/bin/bash
# ./yq ea '.[]'  examples/data*.yaml


./yq ea '
  ((.[] | {.name: .}) as $item ireduce ({}; . * $item )) as $uniqueMap
  | ( $uniqueMap  | to_entries | .[]) as $item ireduce([]; . + $item.value)
' examples/data*.yaml