#!/bin/sh

ver=$1
test -z "$ver" && echo 'no version provided, exit' && exit

cd build

cp ../LICENSE ./
pandoc -o README.html ../README.md
pandoc -o README.pdf ../README.md

dir="upload"
list=$(go tool dist list|grep -E 'linux|bsd|darwin|windows|dragon' |grep -E 'arm|amd')
for i in $list
do
	os=$(echo $i|cut -d/ -f1)
	arch=$(echo $i|cut -d/ -f2)
	echo "building ${i} ...."
  file="ncli_${os}_${arch}_${ver}"
	if test "$os" = "windows";then
	  env GOOS=$os GOARCH=$arch go build -o "${file}.exe" ..
	  zip -j -r "${dir}/${file}.zip" "${file}.exe" LICENSE README.html README.pdf
	else
	  env GOOS=$os GOARCH=$arch go build -o "${file}" ..
	  chmod a+x "$file"
	  tar czvf "${dir}/${file}.tar.gz" "$file" LICENSE README.html README.pdf
  fi
done

cd "$dir"

git archive --format=zip --output=src.zip HEAD
git archive --format=tar.gz --output=src.tar.gz HEAD

checksum="checksum.txt"
for i in ncli_* src.*;
do
  shasum --tag --algorithm 256 "$i" >> "$checksum"
done

