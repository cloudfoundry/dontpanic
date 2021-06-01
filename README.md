# DONTPANIC!!

`dontpanic` is a tool for debugging issues with [Garden containers](https://github.com/cloudfoundry/garden-runc-release) and their host environment.
It collects and tars all necessary data to help Garden engineers investigate bugs. This data includes Garden logs and general system information.
It should not contain any sensitive information, but you are free to review before sending to us. The Garden team is comprised of engineers
from multiple companies and all bugs are investigated together. Your report will not be shared outside the team. A full list of what is collected can be found below.

From GRR v1.17.1 `dontpanic` comes installed on all VMs running the Garden job.

It should be run as root and the resulting tar sent to the Garden team: `/var/vcap/packages/dontpanic/bin/dontpanic`.

Those running GRR < v1.17.1 can download the latest `dontpanic` release and execute it on the host VM as root:

eg: `wget https://github.com/cloudfoundry/dontpanic/releases/download/v1.1/dontpanic && chmod +x ./dontpanic && ./dontpanic`.

NB: If you are running the Garden job in rootless mode (ie Garden is running inside a BPM container), you should still execute `dontpanic` as root
from outside the BPM container.

## What is in my report?

- The current date
- The machine's uptime and current load
- The deployed `gdn` version
- The machine hostname
- Free memory
- Operating system and kernel information
- Monit summary
- Monit logs
- The number of running garden containers
- The number of open files
- The max number of open files permitted on the machine
- The current disk usage
- A list of all open files
- Process table
- Process tree
- Kernel logs
- System logs
- Garden logs
- Network interfaces
- IP tables
- The mount table
- A list of the contents of Garden's depot (container metadata store) dir
- XFS filesystem information
- Memory structure information
- General VM statistics (IO, Memory etc etc)
- General process information

_You can inspect which commands are being run to gather the above by looking at the [code](https://github.com/cloudfoundry/dontpanic/blob/b5ca462b248fba3ff76afcb93b4cb20bf6dfbfce/main.go#L26-L61)_

## How can I use the data in the report?

### Sysstat

In the sysstat folder you can find multiple files containing system statistics (CPU, Memory, I/O, ...) over time.

In order to make use of this information, you need to do the following:

```
export LC_ALL=C
for file in $(ls sysstat/sa[0-9]*) ; do sar -A -f "$file"  >> sa.data.txt; done
```

and then use [`ksar`](https://www.cyberciti.biz/tips/identifying-linux-bottlenecks-sar-graphs-with-ksar.html) to turn the result into pdf graphs.

There are 2 types of files in `sysstat`: `sa*` and `sar*`. The `sa`s are binaries updated every 10 mins or so. The `sar`s are text files generated once per day. Therefore you probably want to parse the `sa` files as they will be more current.

Note: `ksar` seems to dislike some lines in that file and will complain. What you can do is keep removing the zero lines until it is happy.
