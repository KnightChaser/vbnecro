# vbnecro
> **This project is under development and not stable. It's currently under at idea realization phase.** 

`vbnecro` is a simple tool for controlling VirtualBox machines, especially Linux-based guests. It allows you to perform operations such as managing snapshots, starting, executing shell commands inside the guest, and shutting down the VMs.

## Constraints

- **Manual Setup:** Users must manually configure the machines on VirtualBox.
- **Guest Additions:** The guest OS must have VirtualBox Guest Additions installed.
- **Auto-Login:** For guest control commands to work reliably, auto-login must be enabled on the guest.

## Example Configuration File

Below is an example `config.yaml` file:

```yaml
vms:
  - alias: "vm/vbnecro_ubuntu2204"
    vm_name: "vbnecro_ubuntu2204"
    username: "vbnecro"
    password: "pass12##"

jobs:
  - vm_alias: "vm/vbnecro_ubuntu2204"
    ensure_off: true
    operations:
      - type: "RestoreSnapshot"
        params:
          snapshot: "Setup003"
      - type: "StartVM"
      - type: "ExecuteShellCommand"
        params:
          command: "ls"
          args:
            - "-l"
            - "/home/vbnecro"
      - type: "ShutdownVM"

```

## Usage

1.  **Build the project:**
    
    ```bash
    go build .
    
    ```
    
2.  **Run vbnecro with your configuration:**
    
    ```bash
    ./vbnecro --config-path=./config.yaml
    ```
   

