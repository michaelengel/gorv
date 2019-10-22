signed char c;
signed char *p;
signed char *s = "\xfeU\n";

void print(signed char c) {
}

int main(void) {
	p = s;
	while (*p) {
		c=*p++;
		print(c);
	}
}

