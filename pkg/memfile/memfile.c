#define _GNU_SOURCE         /* See feature_test_macros(7) */
#include <dlfcn.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <sys/utsname.h>
#include <unistd.h>

#define __NR_memfd_create 319

extern int IsModernKernel(void);

int memfd_create(const char *name, unsigned int flags) {
    return syscall(__NR_memfd_create, name, flags);
}

// Returns a file descriptor where we can write our shared object
int open_ramfs(char* shm_name) {
    int shm_fd;

    if (IsModernKernel() == 1) {
        shm_fd = shm_open(shm_name, O_RDWR | O_CREAT, S_IRWXU);
        if (shm_fd < 0) { //Something went wrong :(
            return -1;
        }
    }
    else {
        shm_fd = memfd_create(shm_name, 1);
        if (shm_fd < 0) { //Something went wrong :(
            return -1;
        }
    }
    return shm_fd;
}