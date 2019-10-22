int main(void) {
	int i, j;

	j = 0x12340000;
	for (i=0; i<10; i++) {
		// j = j * i;
		j = j + i;
	}
	return j;
}
