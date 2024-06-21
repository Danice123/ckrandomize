#!/bin/sh
mkdir $1/$2
cp "Pokemon - Crystal Kaizo Plus.gbc" $1/$2/base.gbc
./ckrandomize randomize -s $2 $1/$2/base.gbc $1/$2/$2.gbc
flips -c $1/$2/base.gbc $1/$2/$2.gbc $1/$2/$2.ips
rm $1/$2/base.gbc $1/$2/$2.gbc
cd $1/$2
zip $2.zip $2.ips $2.json
rm $2.ips $2.json
