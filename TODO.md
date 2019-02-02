# TODO

## Drift detection

- DriftInformation

## Status detection

- Stream stack events to console
- StackStatus, StackStatusReason?

## Package command

<https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-cli-package.html>

## Change set descriptions

```yaml
stacks:
  - changeSetDescription: ${{env:COMMIT_MESSAGE}}
```

## Override stack policy

```yaml
defaults:
  # how to make this one-time use?
  stackPolicyDuringChangeSet:
    StackPolicyBody: null
    StackPolicyURL: null
```

## Cost threshold

```yaml
defaults:
  limitCost:
    maximum: 1000
```

## Exit codes

```yaml
defaults:
  exitCodes:
    # enforce limitCost
    onCost: 2

    # do not proceed blindly if drift is detected
    onDrift: 1
```
