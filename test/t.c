int q = 4711;
int j = 0x12340000;

int add(int x, int y) { 
   return x+y;
} 

int foo(int x, int y) {
   return add(x,y);
}

int main(void) {
	int i; // , j;

	// j = 0x12340000;
	for (i=0; i<10; i++) {
		j = foo(j, i);
	}
	return j;
}

