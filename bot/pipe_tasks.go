package bot

import "github.com/lnxjedi/gopherbot/robot"

import "os"

// func template(m robot.Robot, args ...string) (retval robot.TaskRetVal) {
// 	r := m.(Robot)
// 	return
// }

func setenv(m robot.Robot, args ...string) (retval robot.TaskRetVal) {
	r := m.(Robot)
	if len(args) != 2 {
		r.Log(robot.Error, "task 'setenv' called with %d args != 2", len(args))
		return robot.Fail
	}
	os.Setenv(args[0], args[1])
	return
}

func restart(m robot.Robot, args ...string) (retval robot.TaskRetVal) {
	r := m.(Robot)
	pn := r.pipeName
	state.Lock()
	if state.shuttingDown {
		state.Unlock()
		Log(robot.Warn, "Restart triggered in pipeline '%s' with shutdown already in progress", pn)
		return
	}
	running := state.pipelinesRunning - 1
	state.shuttingDown = true
	state.restart = true
	state.Unlock()
	r.Log(robot.Info, "Restart triggered in pipeline '%s' with %d pipelines running (including this one)", pn, running)
	go stop()
	return
}

func pause(m robot.Robot, args ...string) (retval robot.TaskRetVal) {
	r := m.(Robot)
	w := getLockedWorker(r.tid)
	w.Unlock()
	resume := make(chan struct{})
	brainLocks.Lock()
	brainLocks.locks[w.id] = resume
	brainLocks.Unlock()
	pauseBrain(w.id, resume)
	return
}

func resume(m robot.Robot, args ...string) (retval robot.TaskRetVal) {
	r := m.(Robot)
	w := getLockedWorker(r.tid)
	w.Unlock()
	brainLocks.Lock()
	if resume, ok := brainLocks.locks[w.id]; ok {
		close(resume)
	}
	brainLocks.Unlock()
	return
}

func init() {
	RegisterTask("setenv", true, robot.TaskHandler{Handler: setenv})
	RegisterTask("restart-robot", true, robot.TaskHandler{Handler: restart})
	RegisterTask("pause-brain", true, robot.TaskHandler{Handler: pause})
	RegisterTask("resume-brain", true, robot.TaskHandler{Handler: resume})
}
