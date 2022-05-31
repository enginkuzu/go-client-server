# go-client-server
Go Client &amp; Server Application for B!nalyze Coding Challenge  
For more information you can see task.pdf file

### Building Client &amp; Server
```
$ ./build-client.sh 
-rwxrwxr-x 1 engin engin 2,0M May 31 06:15 go-client
$ ./build-server.sh 
-rwxrwxr-x 1 engin engin 2,0M May 31 06:15 go-server
```
### Client Application
```
$ ./go-client --help
Usage: go-client [options]

Starts TCP/IP client application.

Options:
--id        specifies client ID (uint32)
--server    specifies server ipv4/ipv6 address (string)
--version   output version information and exit
--help      display this help and exit
```

```
$ ./go-client --id 7 --server 127.0.0.1
Starting GO TCP/IP Client v0.0.1 (build with go1.16.15)
With Parameters : id=7, server="127.0.0.1"
Disk scan starting ...
Client connected from 127.0.0.1:44952 to 127.0.0.1:7654
Disk scan found 1118385 files in 104102 folders (38 files and 1932 folders can't read) (total 6 seconds)
Queue closed
App ending ...
App end
```
### Server Application
```
$ ./go-server 
Starting GO TCP/IP Server v0.0.1 (build with go1.16.15)
Bind on 0.0.0.0:7654 success
New connection from 127.0.0.1:44952 to 127.0.0.1:7654
Client send 1118385 messages in 6 seconds
Connection closed from 127.0.0.1:44952
^C"interrupt" signal received
App ending ...
App end
```
### Created Files
```
$ cat client_log.txt
2022-05-31 06:16:38.670 +03 : INFO  : Starting GO TCP/IP Client v0.0.1 (build with go1.16.15)
2022-05-31 06:16:38.671 +03 : INFO  : With Parameters : id=7, server="127.0.0.1"
2022-05-31 06:16:38.671 +03 : INFO  : Disk scan starting ...
2022-05-31 06:16:38.672 +03 : INFO  : Client connected from 127.0.0.1:44952 to 127.0.0.1:7654
2022-05-31 06:16:38.672 +03 : ERROR : Folder read error : open /boot/efi/: permission denied
2022-05-31 06:16:38.680 +03 : ERROR : Folder read error : open /dev/vboxusb/: permission denied
2022-05-31 06:16:38.686 +03 : ERROR : Folder read error : open /etc/cups/ssl/: permission denied
2022-05-31 06:16:41.474 +03 : ERROR : File info error : lstat /proc/52330/fd//7: no such file or directory
2022-05-31 06:16:41.474 +03 : ERROR : File info error : lstat /proc/52330/fdinfo//7: no such file or directory
2022-05-31 06:16:44.723 +03 : INFO  : Disk scan found 1118385 files in 104102 folders (38 files and 1932 folders can't read) (total 6 seconds)
2022-05-31 06:16:44.723 +03 : INFO  : Queue closed
2022-05-31 06:16:44.880 +03 : INFO  : App ending ...
```
```
$ cat server_log.txt
2022-05-31 06:15:51.119 +03 : INFO  : Starting GO TCP/IP Server v0.0.1 (build with go1.16.15)
2022-05-31 06:15:51.121 +03 : INFO  : Bind on 0.0.0.0:7654 success
2022-05-31 06:16:38.672 +03 : INFO  : New connection from 127.0.0.1:44952 to 127.0.0.1:7654
2022-05-31 06:16:44.880 +03 : INFO  : Client send 1118385 messages in 6 seconds
2022-05-31 06:16:44.880 +03 : INFO  : Connection closed from 127.0.0.1:44952
2022-05-31 06:21:01.507 +03 : INFO  : "interrupt" signal received
2022-05-31 06:21:01.669 +03 : INFO  : App ending ...
```
```
$ cat filesystem_dump.csv
Message Size,Protocol,Client Id,File Size,File Permissions,File Path Length,File Path,Checksum
42,0,7,2147483648,384,9,/swapfile,412158245
67,0,7,5964660,384,34,/boot/System.map-5.13.0-44-generic,3453674136
63,0,7,257795,420,30,/boot/config-5.13.0-44-generic,145879256
49,0,7,28,134218239,16,/boot/initrd.img,933515211
67,0,7,63778788,420,34,/boot/initrd.img-5.13.0-44-generic,218964445
53,0,7,182704,420,20,/boot/memtest86+.bin,679057876
53,0,7,184380,420,20,/boot/memtest86+.elf,1051781360
46,0,7,25,134218239,13,/boot/vmlinuz,2758831850
64,0,7,10176896,384,31,/boot/vmlinuz-5.13.0-44-generic,3564680599
62,0,7,9,384,29,/var/spool/anacron/cron.daily,1872751695
64,0,7,9,384,31,/var/spool/anacron/cron.monthly,2023839992
63,0,7,9,384,30,/var/spool/anacron/cron.weekly,1876801882
```
