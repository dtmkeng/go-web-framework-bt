#!/bin/bash
# go build -o gowebbenchmark
server_bin_name="gowebbenchmark"

. ./libs.sh

length=${#web_frameworks[@]}

test_result=()
test_alloc_result=()
cpu_cores=`cat /proc/cpuinfo|grep processor|wc -l`
if [ $cpu_cores -eq 0 ]
then
  cpu_cores=1
fi

test_web_framework()
{
  echo "testing web framework: $2"
  ./$server_bin_name $2 $3  > alloc.log 2>&1 &
  sleep 2

  throughput=`wrk -t4 -c$4 -d30s http://127.0.0.1:8081/hello | grep Requests/sec | awk '{print $2}'`
  alloc=`cat alloc.log|grep HeapAlloc | awk '{print $2}'`
  echo "throughput: $throughput requests/second"
  echo "cpu_cores: $cpu_cores"
  test_result[$1]=$throughput
  test_alloc_result[$1]=$alloc

  pkill -9 $server_bin_name
  sleep 2
  echo "finsihed testing $2"
  echo
}

test_all()
{
  echo "###################################"
  echo "                                   "
  echo "      ProcessingTime  $1ms         "
  echo "      Concurrency     $2           "
  echo "                                   "
  echo "###################################"
  for ((i=0; i<$length; i++))
  do
  	test_web_framework $i ${web_frameworks[$i]} $1 $2
  done
}


pkill -9 $server_bin_name
echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > processtime_alloc.csv
echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > processtime.csv
test_all 0 100
echo "0 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
echo "0 ms,"$(IFS=$','; echo "${test_alloc_result[*]}" ) >> processtime_alloc.csv
test_all 10 100
echo "10 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
echo "10 ms,"$(IFS=$','; echo "${test_alloc_result[*]}" ) >> processtime_alloc.csv
test_all 100 100
echo "100 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
echo "100 ms,"$(IFS=$','; echo "${test_alloc_result[*]}" ) >> processtime_alloc.csv
test_all 500 100
echo "500 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
echo "500 ms,"$(IFS=$','; echo "${test_alloc_result[*]}" ) >> processtime_alloc.csv


