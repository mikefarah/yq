#!/bin/bash

setUp() {
    rm test*.yml || true
}

## Convenient bash shortcut to read records of NUL separated values
## from stdin the safe way. See example usage in the next tests.
read-0() {
    local eof="" IFS=''
    while [ "$1" ]; do
        ## - The `-r` avoids bad surprise with '\n' and other interpreted
        ##   sequences that can be read.
        ## - The `-d ''` is the (strange?) way to refer to NUL delimiter.
        ## - The `--` is how to avoid unpleasant surprises if your
        ##   "$1" starts with "-" (minus) sign. This protection also
        ##   will produce a readable error if you want to try to start
        ##   your variable names with a "-".
        read -r -d '' -- "$1" || eof=1
        shift
    done
    [ -z "$eof" ] ## fail on EOF
}

## Convenient bash shortcut to be used with the next function `p-err`
## to read NUL separated values the safe way AND catch any errors from
## the process creating the stream of NUL separated data.  See example
## usage in the tests.
read-0-err() {
    local ret="$1" eof="" idx=0 last=
    read -r -- "${ret?}" <<<"0"
    shift
    while [ "$1" ]; do
        last=$idx
        read -r -d '' -- "$1" || {
            ## Put this last value in ${!ret}
            eof="$1"
            read -r -- "$ret" <<<"${!eof}"
            break
        }
        ((idx++))
        shift
    done
    [ -z "$eof" ] || {
        if [ "$last" != 0 ]; then
            ## Uhoh, we have no idea if the errorlevel of the internal
            ## command was properly delimited with a NUL char, and
            ## anyway something went really wrong at least about the
            ## number of fields separated by NUL char and the one
            ## expected.
            echo "Error: read-0-err couldn't fill all value $ret = '${!ret}', '$eof', '${!eof}'" >&2
            read -r -- "$ret" <<<"not-enough-values"
        else
            if ! [[ "${!ret}" =~ ^[0-9]+$ && "${!ret}" -ge 0 && "${!ret}" -le 127 ]]; then
                ## This could happen if you don't use `p-err` wrapper,
                ## or used stdout in unexpected ways in your inner
                ## command.
                echo "Error: last value is not a number, did you finish with an errorlevel ?" >&2
                read -r -- "$ret" <<<"last-value-not-a-number"
            fi
        fi
        false
    }
}

## Simply runs command given as argument and adds errorlevel in the
## standard output. Is expected to be used in tandem with
## `read-0-err`.
p-err() {
    local exp="$1"
    "$@"
    printf "%s" "$?"
}

wyq-r() {
    local exp="$1"
    ./yq e -0 -r=false "$1"
    printf "%s" "$?"
}

testBasicUsageRaw() {
  cat >test.yml <<EOL
a: foo
b: bar
EOL

  printf "foo\0bar\0" > expected.out

  ## We need to compare binary content here. We have to filter the compared
  ## content through a representation that gets rid of NUL chars but accurately
  ## transcribe the content.
  ## Also as it would be nice to have a pretty output in case the test fails,
  ## we use here 'hd': a widely available shortcut to 'hexdump' that will
  ## pretty-print any binary to it's hexadecimal representation.
  ##
  ## Note that the standard `assertEquals` compare its arguments
  ## value, but they can't hold NUL characters (this comes from the
  ## limitation of the C API of `exec*(..)` functions that requires
  ## `const char *arv[]`). And these are NUL terminated strings.  As a
  ## consequence, the NUL characters gets removed in bash arguments.
  assertEquals "$(hd expected.out)" \
               "$(./yq e -0 '.a, .b' test.yml | hd)"

  rm expected.out
}

testBasicUsage() {
  local a b
  cat >test.yml <<EOL
a: foo
b: bar
EOL

  ## We provide 2 values, and ask to fill 2 variables.
  read-0 a b < <(./yq e -0 '.a, .b' test.yml)
  assertEquals "$?" "0"      ## Everything is fine
  assertEquals "foo" "$a"    ## Values are correctly parsed
  assertEquals "bar" "$b"

  a=YYY ; b=XXX
  ## Not enough values provided to fill `a` and `b`.
  read-0 a b < <(./yq e -0 '.a' test.yml)
  assertEquals "$?" "1"      ## An error was emitted
  assertEquals "foo" "$a"    ## First value was correctly parsed
  assertEquals "" "$b"       ## Second was still reset

  ## Error from inner command are not catchable !. Use
  ## `read-0-err`/`p-err` for that.
  read-0 a < <(printf "\0"; ./yq e -0 'xxx' test.yml; )
  assertEquals "$?" "0"

}

testBasicUsageJson() {
  cat >test.yml <<EOL
a:
  x: foo
b: bar
EOL

  read-0 a b < <(./yq e -0 -o=json '.a, .b' test.yml)

  assertEquals '{
  "x": "foo"
}' "$a"
  assertEquals '"bar"' "$b"

}

testFailWithValueContainingNUL() {
  local a b c
  ## Note that value of field 'a' actually contains a NUL char !
  cat >test.yml <<EOL
a: "foo\u0000bar"
b: 1
c: |
  wiz
  boom
EOL

  ## We are looking for trouble with asking to separated fields with NUL
  ## char and requested value `.a` actually contains itself a NUL char !
  read-0 a b c < <(./yq e -0 '.a, .b, .c' test.yml)
  assertNotEquals "0" "$?"   ## read-0 failed to fill all values

  ## But here, we can request for one value, even if `./yq` fails
  read-0 b < <(./yq e -0 '.b, .a' test.yml)
  assertEquals "0" "$?"   ## read-0 succeeds at feeding the first value
  ## Note: to catch the failure of `yq`, see in the next tests the usage
  ## of `read-0-err`.

  ## using -r=false solves any NUL containing value issues, but keeps
  ## all in YAML representation:
  read-0 a b c < <(./yq e -0 -r=false '.a, .b, .c' test.yml)
  assertEquals "0" "$?"    ## All goes well despite asking for `a` value

  assertEquals '"foo\0bar"' "$a"   ## This is a YAML string representation
  assertEquals '1' "$b"
  assertEquals '|
  wiz
  boom' "$c"
}

testStandardLoop() {
    local E a b res

    ## Here everything is normal: 4 values, that will be paired
    ## in key/values.
    cat >test.yml <<EOL
- yay
- wiz
- hop
- pow
EOL

    res=""
    while read-0-err E a b; do
        res+="$a: $b;"
    done < <(p-err ./yq -0 '.[]' test.yml)

    assertEquals "0" "$E"                     ## errorlevel of internal command
    assertEquals "yay: wiz;hop: pow;" "$res"  ## expected result
}

testStandardLoopWithoutEnoughValues() {
    local E a b res

    ## Here 5 values, there will be a missing value when reading
    ## pairs of value.
    cat >test.yml <<EOL
- yay
- wiz
- hop
- pow
- kwak
EOL

    res=""
    ## The loop will succeed 2 times then fail
    while read-0-err E a b; do
        res+="$a: $b;"
    done < <(p-err ./yq -0 '.[]' test.yml)

    assertEquals "not-enough-values" "$E"     ## Not enough value error
    assertEquals "yay: wiz;hop: pow;" "$res"  ## the 2 full key/value pairs

}

testStandardLoopWithInternalCmdError() {
    local E a b res

    ## Note the third value contains a NUL char !
    cat >test.yml <<EOL
- yay
- wiz
- "foo\0bar"
- hop
- pow
EOL

    res=""
    ## It should be only upon the second pass in the loop that
    ## read-0-err will catch the fact that there is an error !
    while read-0-err E a b; do
        res+="$a: $b;"
    done < <(p-err ./yq -0 '.[]' test.yml)
    assertEquals "1" "$E"            ## Internal command errorlevel (from `./yq`)
    assertEquals "yay: wiz;" "$res"  ## first 2 values were ok at least

}

testStandardLoopNotEnoughErrorEatsCmdError() {
    local E a b res

    ## Because of possible edge cases where the internal errorlevel
    ## reported by `p-err` in the standard output might be mangled
    ## with the unfinished record, `read-0-err E ...` will NOT report
    ## the internal command error in the variable E and instead will
    ## store the value 'not-enough-values'. In real world, anyway, you
    ## will want to react the same if the internal command failed
    ## and/or you didn't get as much values as expected while
    ## reading. Keep in mind also that standard error is not
    ## swallowed, so you can read reports from the inner command AND
    ## from `read-0-err`.

    ## Here, note that the fourth value contains a NUL char !
    cat >test.yml <<EOL
- yay
- wiz
- hop
- "foo\0bar"
- pow
EOL

    res=""
    ## It should be only upon the second loop that read-0-err will catch
    ## the fact that there are not enough data to fill the requested variables
    while read-0-err E a b; do
        res+="$a: $b;"
    done < <(p-err ./yq -0 '.[]' test.yml)
    assertEquals "not-enough-values" "$E"          ## Not enough values error eats internal error !
    assertEquals "yay: wiz;" "$res"  ## first 2 values were ok at least
}


source ./scripts/shunit2