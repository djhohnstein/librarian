typedef struct{
    int length;
    char** results;
    char* name;
    int message_type;
}datagram;
typedef datagram* (*moduleFunction)(datagram*, void*, datagram*);
typedef datagram* (*moduleCallback)(datagram*);

typedef struct{
    moduleFunction mainFunction;
    moduleCallback mainCallback;
}moduleFunctions;

void* load_module(char*, char*, char*);
void call_module_callback(void*, void*);
void* call_module_function(void*, void*);