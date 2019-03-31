# TODO

## Next

1. Package command

   <https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-cli-package.html>

1. Nested stacks

    <https://theithollow.com/2018/01/29/using-change-sets-nested-cloudformation-stacks/>

## Backlog

1. Drift detection

   - DriftInformation

1. Change set descriptions

   ```yaml
   stacks:
     - changeSetDescription: ${{env:COMMIT_MESSAGE}}
   ```

1. Override stack policy

   ```yaml
   defaults:
     # how to make this one-time use?
     stackPolicyDuringChangeSet:
       StackPolicyBody: null
       StackPolicyURL: null
   ```

1. Cost threshold

   ```yaml
   defaults:
     limitCost:
       maximum: 1000
   ```

1. Exit codes

   ```yaml
   defaults:
     exitCodes:
       # enforce limitCost
       onCost: 2

       # do not proceed blindly if drift is detected
       onDrift: 1
   ```
