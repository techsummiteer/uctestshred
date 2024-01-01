A Shred tool

When I need to shred a file, I am typically in a hurry. 

This tool uses the parallelism in go coroutines to shred a file. I still leave around 50% capaity on the CPU for other tasks.

Note:

Sherd-ing on wear balancing  SSD/eMMC will not work as the same physical location on the SSD/eMMC may not be written on a repeat write to the same file offset.
