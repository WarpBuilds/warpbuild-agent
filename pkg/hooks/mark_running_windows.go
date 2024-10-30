package hooks

// import (
// 	"context"
// 	"fmt"
// 	"runtime"

// 	"github.com/warpbuilds/warpbuild-agent/pkg/log"
// 	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
// 	"golang.org/x/sys/windows/svc"
// 	"golang.org/x/sys/windows/svc/mgr"
// )

// const MARK_WINDOWS_RUNNING_HOOK string = "MARK_WINDOWS_RUNNING_HOOK"

// type MarkWindowsRunningHook struct{}

// var _ manager.IPreStartHook = &MarkWindowsRunningHook{}

// func init() {
// 	manager.RegisterHook[manager.IPreStartHook](&MarkWindowsRunningHook{})
// }

// func (*MarkWindowsRunningHook) HookID() string {
// 	return MARK_WINDOWS_RUNNING_HOOK
// }

// func (*MarkWindowsRunningHook) PreStartHook(ctx context.Context, opts *manager.PreStartHookOptions) error {
// 	log.Logger().Debugf("Checking if machine is windows...")

// 	if runtime.GOOS != "windows" {
// 		log.Logger().Debugf("Machine is not windows, skipping...")
// 		return nil
// 	}

// 	log.Logger().Debugf("Machine is windows, marking as running...")

// 	if opts.StartRunnerOptions.AgentOptions.WindowsOptions == nil {
// 		log.Logger().Errorf("Windows options are not set but machine is windows")
// 		return fmt.Errorf("windows options are not set")
// 	}

// 	var serviceName = opts.StartRunnerOptions.AgentOptions.WindowsOptions.ServiceName

// 	if serviceName == "" {
// 		log.Logger().Errorf("Windows service name is not set")
// 		return fmt.Errorf("windows service name is not set")
// 	}

// 	// Register service with SCM and send status updates
// 	serviceStatusHandle, err := mgr.Connect()
// 	if err != nil {
// 		log.Logger().Errorf("Could not connect to SCM: %v", err)
// 		return err
// 	}
// 	defer serviceStatusHandle.Disconnect()

// 	service, err := serviceStatusHandle.OpenService(serviceName)
// 	if err != nil {
// 		log.Logger().Errorf("Could not open service: %v", err)
// 		return err
// 	}
// 	defer service.Close()

// 	status := svc.Status{
// 		State: svc.Running,
// 	}
// 	err = service.UpdateConfig(status)
// 	if err != nil {
// 		log.Logger().Errorf("Could not set service status: %v", err)
// 		return err
// 	}

// 	return nil
// }
