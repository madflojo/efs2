# Efs2file Instructions

An `Efs2file` is much like a `Dockerfile` in that it uses a simple set of instructions to execute actions against target servers. The number of instructions within Efs2 is small but mighty. This page serves as a reference for available instructions.

| Instruction | Example | Description |
| ----------- | ------- | ----------- |
| `RUN` | `RUN apt-get install something` | Execute command against remote target system |
| `PUT` | `PUT local.file /path/to/remote.file 0644` | Upload file to the target system using the destination path and permissions mode |

## Legacy Instructions

The following instructions are legacy however still officially supported within Efs2.

| Instruction | Example | Description |
| ----------- | ------- | ----------- |
| `RUN CMD` | `RUN CMD apt-get install something` | Execute command against remote target systems |
| `RUN SCRIPT` | `RUN SCRIPT localfile.sh` | Upload and Execute the provided script on the remote target systems |
