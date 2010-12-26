// Utility library

// http://stackoverflow.com/questions/284325/how-to-make-child-process-die-after-parent-exits

#ifdef __LINUX__
#include <sys/prctl.h>
#elif __POSIX__
#include <sys/time.h>
#endif

#include <sys/types.h>
#include <signal.h>
#include "watch_parent.h"

static *sighandler_t watch_parent_callback = 0;

int watch_parent_handler(int sig) {
	static triggered = 0;
	if (!triggered && watch_parent_callback && getppid() == 1) {
		triggered = 1;
		watch_parent_callback(sig);
	}
}

void watch_parent_setup(sighandler_t handler) {
	if (watch_parent_callback) die("Arrrrghlz.g.gd");
	watch_parent_callback = handler;
#ifdef __LINUX__
	prctl(PR_SET_PDEATHSIG, WATCH_PARENT_SIGNAL);
	signal(WATCH_PARENT_SIGNAL, watch_parent_handler);
#elif __POSIX__
	watch_parent_callback = handler;
	signal(SIGVTALRM, watch_parent_handler);
	// TODO: define alarm clock here
#else
#error Not available
#endif
}
