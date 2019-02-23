#!/bin/bash -x
echo "generate a rc build number"
github-release info | grep -A 20 "^git\ tags"
next_semver=`github-release info | grep -A 20 "^git\ tags" | grep "^-" | grep -v -e "-rc" | head -1 | awk '{print $2}' | awk -F"." '{print $1"."$2"."$3+1}'`
next_rc_version=`github-release info | grep -A 20 "^git\ tags" | grep "${next_semver}" | grep -e "-rc" | head -1 | awk -F"-rc." '{print $2+1}'`
if [[ "${next_rc_version/ /}" == "" ]]; then
  echo "setting default base rc"
  next_rc_version=0
fi
echo "next rc id is: "$next_rc_version
echo "generate a release tag for this RC"
tag=${next_semver}-rc.${next_rc_version}
echo "tag id is: "$tag
echo "creating release"
github-release release -t ${tag} -p
echo "uploading files"
for file in `ls build | grep '^clarity'`
do
  github-release upload -t ${tag} -f build/${file} -n ${file}
done
