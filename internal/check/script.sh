#! /bin/bash

check_file=${1:-/opt/url-list}

# Check if the file exists
if [ ! -f ${check_file} ]; then
    echo "URL list file not found"
    sleep 600
    exit 1
fi

# infinite loop with sleep of 5 minutes
while true
do
    # For entry in the file try to curl the URL
    while IFS= read -r line
    do
        echo -n "$(date +%T) "
        curl -s -o /dev/null --connect-timeout 10 --max-time 30 -w "%{http_code} %{url_effective} %{errormsg} %{redirect_url}\\n" $line
    done < ${check_file}

    echo -e "Sleeping for ${INTERVAL:-300} seconds \n"
    sleep ${INTERVAL:-300}
done
