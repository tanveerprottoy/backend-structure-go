#!/bin/sh

# There is a known issue if using submodules. You can see the error below. This is due
# to git diff in the script can include the hash from the submodules. This error can be ignored.
# fatal: git cat-file: could not get object info

commitlimit=$(git config hooks.filesizehardlimit)
filelimit=$(git config hooks.filesizehardlimit)
: ${commitlimit:=10}
: ${filelimit:=10}

list_new_or_modified_files()
{
    # Only get the file name.  Do not get deleted files.
    git diff --staged --name-only --diff-filter=ACMRTUXB
}

unmunge()
{
    local result="${1#\"}"
    result="${result%\"}"
    env echo "$result"
}

check_file_size()
{
    n=0
    while read -r munged_filename
    do
        f="$(unmunge "$munged_filename")"
        h=$(git ls-files -s "$f"|cut -d' ' -f 2)
        s=$(git cat-file -s "$h")
        if [[ "$s" -gt $filelimit ]]
        then
            env echo 1>&2 "ERROR: size limit ($filelimit) exceeded: $munged_filename ($s)"
            n=$((n+1))
        fi
        fs=$(($fs+$s))
    done
    if [[ "$fs" -gt $limit ]]
    then
       env echo 1>&2 "ERROR: hard size limit ($commitlimit) bytes exceeded: $munged_filename ($fs)"
    fi
    [ $n -eq 0 ]
}

list_new_or_modified_files | check_file_size
exit $n
