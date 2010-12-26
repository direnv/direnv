#include <stdio.h>
#include <sys/types.h>
#include <signal.h>

#include "watch_parent.h"

void parent_died(int signal) {
	exit(0);
}

int main(void) {
	watch_parent(SIGUSR2);
	signal(SIGUSR2, parent_died);

	getcwd()
	return 0;
}
