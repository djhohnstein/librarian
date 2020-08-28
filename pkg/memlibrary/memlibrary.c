
#include <dlfcn.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include "memlibrary.h"


extern void RouteDataFromModule(datagram*);

void cRouteDataFromModule(datagram* user_data) {
  RouteDataFromModule(user_data);
}

void* load_library(char* path) {
//    printf("Path: %s\n", path);
//    char char_path[1024];
//    snprintf(char_path, 1024, path);
//    printf("charpath: %s\n", char_path);
    void* handle;
    handle = dlopen(path, RTLD_LAZY);
    if (!handle) {
        printf("Failed to acquire handle\n");
        return NULL;
    }
    return handle;
}

void call_module_callback(void* fptr, void* arguments) {
    (*(moduleCallback)fptr)(arguments);
}

void* call_module_function(void* fptr, void* arguments) {
    datagram* results = malloc(sizeof(datagram));
    (*(moduleFunction)fptr)(arguments, cRouteDataFromModule, results);
    return results;
}

moduleFunction get_main_export(void* hLibrary, char* functionName) {
    moduleFunction fptr;
    *(void**)(&fptr) = dlsym(hLibrary, functionName);
    if (fptr == NULL) {
//        printf("Failed to find main module\n");
        return NULL;
    }
//    printf("Got main function");
    return fptr;
}

moduleCallback get_main_callback_export(void* hLibrary, char* functionName) {
    moduleCallback fptr;
    *(void**)(&fptr) = dlsym(hLibrary, functionName);
    if (fptr == NULL) {
//        printf("Failed to find callback export: was null\n");
        return NULL;
    }
//    printf("Got callback function");
    return fptr;
}

void* load_module(char* libraryPath, char* functionName, char* callbackName) {
    moduleFunctions* funcs;
    void* hLib;
    moduleFunction mainFunc;
    moduleCallback mainCb;

    hLib = load_library(libraryPath);
    if (hLib == NULL) {
        return NULL;
    }
    funcs = (moduleFunctions*)malloc(sizeof(moduleFunctions));
    mainFunc = get_main_export(hLib, functionName);
    mainCb = get_main_callback_export(hLib, callbackName);
    funcs->mainFunction = mainFunc;
    funcs->mainCallback = mainCb;
    return funcs;
}

