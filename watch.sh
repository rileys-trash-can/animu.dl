#! /bin/sh

watch "echo -n \"unique animu downloads: \"; ls -1 out | wc -l; echo -n "in:"; du -sh out"
