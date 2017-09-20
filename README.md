hal-ops
=======

[WIP] hal wrapper CLI command

## Concepts

- Send P-R
- Review the changes
- Login the hal instance with ssh
- Run `git pull`
- Checkout the branch
    - git checkout -> `hal-ops check` (validation only)?
    - `hal-ops check "branch-name"`   (checkout && validation)?
- Merge the P-R if there are no problems
- `hal-ops deploy`
    - Almost equals to `hal deploy apply`
    - Send an event log to datadog
    - Notify to our Slack that the deployment is done
