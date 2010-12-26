#include <signal.h>

#define WATCH_PARENT_SIGNAL SIGUSR2

// Registers callback to call when parent dies
// CALL ONCE
void watch_parent_setup(sighandler_t handler);

