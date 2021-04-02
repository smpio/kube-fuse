#include <stdlib.h>

void** copy_ptr_array(void** source, size_t size)
{
	void** dest = malloc(sizeof(void*) * (size+1));
	for (size_t i = 0; i < size; i++) {
		dest[i] = source[i];
	}
	dest[size] = NULL;
	return dest;
}

char** copy_str_ptr_array(char** source, size_t size)
{
	return (char**)copy_ptr_array((void**)source, size);
}
