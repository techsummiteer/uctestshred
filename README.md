# A Shred tool

When I need to shred a file, I am typically in a hurry. 

This tool uses the parallelism in go coroutines to shred a file. I still leave around 50% capaity on the CPU for other tasks.

## Shred timing ~  1min 20secs

### $ dd if=/dev/random of=./junk.bin count=62914560
62914560+0 records in
62914560+0 records out
32212254720 bytes (32 GB, 30 GiB) copied, 335.093 s, 96.1 MB/s
### $ time ./Shred -n ./junk.bin 
real	1m20.968s
user	0m27.919s
sys	5m51.446s

## unix shred timing ~4min 
### $ dd if=/dev/random of=./junk.bin count=62914560
62914560+0 records in
62914560+0 records out
32212254720 bytes (32 GB, 30 GiB) copied, 339.184 s, 95.0 MB/s
## $ time shred -n 3 -u ./junk.bin 
real	4m1.447s
user	0m33.423s
sys	0m29.532s




## Note:
Sherd-ing on wear balancing  SSD/eMMC will not work as the same physical location on the SSD/eMMC may not be written on a repeat write to the same file offset.
