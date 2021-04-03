#define FUSE_USE_VERSION 26

#include <fuse.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>

static const char *kube_str = "kube World!\n";
static const char *kube_path = "/kube";

static int kube_getattr(const char *path, struct stat *stbuf)
{
    int res = 0;

    memset(stbuf, 0, sizeof(struct stat));
    if (strcmp(path, "/") == 0) {
        stbuf->st_mode = S_IFDIR | 0755;
        stbuf->st_nlink = 2;
    } else if (strcmp(path, kube_path) == 0) {
        stbuf->st_mode = S_IFREG | 0444;
        stbuf->st_nlink = 1;
        stbuf->st_size = strlen(kube_str);
    } else {
        stbuf->st_mode = S_IFREG | 0444;
        stbuf->st_nlink = 1;
        stbuf->st_size = 0;
    }

    return res;
}

extern char** ReadDir(const char* path);

static int kube_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
                        off_t offset, struct fuse_file_info *fi)
{
    (void) offset;
    (void) fi;

    if (strcmp(path, "/") != 0) {
        return -ENOENT;
    }

    filler(buf, ".", NULL, 0);
    filler(buf, "..", NULL, 0);

    char** entries = ReadDir(path);
    for (char** entry_p = entries; *entry_p != NULL; entry_p++) {
        filler(buf, *entry_p, NULL, 0);
        free(*entry_p);
    }
    free(entries);

    return 0;
}

static int kube_open(const char *path, struct fuse_file_info *fi)
{
    if (strcmp(path, kube_path) != 0) {
        return -ENOENT;
    }

    if ((fi->flags & 3) != O_RDONLY) {
        return -EACCES;
    }

    return 0;
}

static int kube_read(const char *path, char *buf, size_t size, off_t offset,
                     struct fuse_file_info *fi)
{
    size_t len;
    (void) fi;
    if(strcmp(path, kube_path) != 0) {
        return -ENOENT;
    }

    len = strlen(kube_str);
    if (offset < len) {
        if (offset + size > len)
            size = len - offset;
        memcpy(buf, kube_str + offset, size);
    } else {
        size = 0;
    }

    return size;
}

static struct fuse_operations kube_oper = {
    .getattr	= kube_getattr,
    .readdir	= kube_readdir,
    .open		= kube_open,
    .read		= kube_read,
};

int c_main(int argc, char *argv[])
{
    puts("kube-fuse starting");
    return fuse_main(argc, argv, &kube_oper, NULL);
}
