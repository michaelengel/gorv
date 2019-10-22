int fak(int n) { 
   if (n==0) {
      return 1;
   } else {
      return n*fak(n-1);
   }
} 

int main(void) {
	return fak(10);
}

