# GitHub Action

You can use `yq` in your GitHub action, for instance:

```yaml
  - uses: actions/checkout@v2
  - name: Get SDK Version from config
    id: lookupSdkVersion
    uses: mikefarah/yq@master
    with:
      cmd: yq '.renutil.version' 'config.yml'
  - name: Restore Cache
    id: restore-cache
    uses: actions/cache@v2
    with:
      path: ../renpy
      key:  ${{ runner.os }}-sdk-${{ steps.lookupSdkVersion.outputs.result }}
      restore-keys: |
        ${{ runner.os }}-sdk
  # ... more
```

The `yq` action sets a `result` variable in its output, making it available to subsequent steps. In this case it's available as `steps.lookupSdkVersion.outputs.result`.

Details of how the GitHub action itself is configured can be found [here](https://github.com/mikefarah/yq/issues/844#issuecomment-856700574)

If you [enable step debug logging](https://docs.github.com/en/actions/managing-workflow-runs/enabling-debug-logging#enabling-step-debug-logging), you can see additional information about the exact command sent as well as the response returned within the GitHub Action logs.

Thanks @[**devorbitus**](https://github.com/devorbitus)**!**
