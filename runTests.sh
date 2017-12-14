#! /bin/bash

declare EXEC=30
declare args=("1000" "50000" "500000" "5000000" "50000000")

declare TARGET=threadring

echo "Cleaning possible old testes"

make cleanTests

echo "Building Executable"

make build

for arg in "${args[@]}"; do
  echo "Running testes for $arg"
  for (( i = 0; i < 30; i++ )); do
    { time ./threadring.out 5000000 >/dev/null; } 2>&1 | grep real | sed 's/real\t//' >> ./tests/$arg.txt
  done
done
