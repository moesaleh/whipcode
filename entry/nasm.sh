nasm -f elf64 /source.asm -o /tmp/run.o
ld /tmp/run.o -o /tmp/run
./tmp/run
