#!/bin/sh 
#===================== 
#https://github.com/philsong/btcrobot
#===================== 
while : 
do 
echo "Current DIR is " $PWD 
stillRunning=$(ps -ef |grep "$PWD/bin/btcrobot.exe" |grep -v "grep") 
if [ "$stillRunning" ] ; then 
echo "my service was already started by another way" 
echo "Kill it and then startup by this shell, other wise this shell will loop out this message annoyingly" 
kill -9 $pidof $PWD/bin/btcrobot.exe
else 
echo "my service was not started" 
echo "Starting service ..." 
$PWD/bin/btcrobot.exe
echo "my service was exited!" 
fi 
sleep 1
done 