# Integration Test Suite

I haven't found a very good way to test things end-to-end in CI so I've opted to run these on my local machine.

## Setup

Create "Sandbox" Vault to isolate the script from your normal items. The script will maniplate items within that vault,
but you can override it with the `VAULT` environment variable.

## Running

```
; integration/test.bash
```
