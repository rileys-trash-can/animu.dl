#! /bin/sh

watch -n .5 "echo -n \"unique animu downloads: \"; ls -1 -R out | wc -l; echo -n "in:"; du -sh out"
