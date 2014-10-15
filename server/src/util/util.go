package util

import (
    "log"
    "os"
    "syscall"
    "runtime"
)

/*
    // Usage is like below, but it won work for go, because it'll work fine only
    // it run before any thread create, and go's import package couldn't garuntee.
    // One who want to start a daemon could use 'Upstart' or 'Supervisord' for help.
    ret := util.Daemon(1, 0)
    if ret == -1 {
        log.Println("Daemon err")
        os.Exit(-1)
    }
*/

func Daemon(nochdir, noclose int) int {
    darwin := runtime.GOOS == "darwin"

    // already a daemon
    if syscall.Getppid() == 1 {
        return 0
    }

    ret, ret2, err := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
    if err != 0 {
        return -1
    }

    // failure
    if ret2 < 0 {
        os.Exit(-1)
    }

    // handle exeception for darwin
    if darwin && ret2 == 1 {
        ret = 0
    }

    // if we got a good PID, then we call exit the parent process.
    if ret > 0 {
        os.Exit(0)
    }

    // Change the file mode mask
    _ = syscall.Umask(0)

    // create a new SID for the child process
    s_ret, s_err := syscall.Setsid()
    if s_err != nil {
        log.Println(s_err.Error())
    }
    if s_ret < 0 {
        return -1
    }

    if nochdir == 0 {
        os.Chdir("/")
    }

    if noclose == 0 {
        f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
        if e == nil {
            fd := f.Fd()
            syscall.Dup2(int(fd), int(os.Stdin.Fd()))
            syscall.Dup2(int(fd), int(os.Stdout.Fd()))
            syscall.Dup2(int(fd), int(os.Stderr.Fd()))
        }
    }

    return 0
}
