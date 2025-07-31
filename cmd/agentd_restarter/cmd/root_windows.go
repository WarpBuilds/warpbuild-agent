//go:build windows

package cmd

import (
	"context"
	"os"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// Version is set from goreleaser
var Version = "dev"

type flagsStruct struct {
	stdoutFile        string
	stderrFile        string
	elevateUser       string
	elevatePassword   string
	domain            string
	agentdServiceName string
	restartInterval   time.Duration
}

var flags flagsStruct

var rootCmd = &cobra.Command{
	Use:   "agentd-restarter",
	Short: "Restarts the agentd service",
	Long:  `Restarts the agentd service`,
	RunE: func(cmd *cobra.Command, args []string) error {

		lm, err := log.Init(&log.InitOptions{
			StdoutFile: flags.stdoutFile,
			StderrFile: flags.stderrFile,
		})
		if err != nil {
			return err
		}
		defer lm.Sync()

		log.Logger().Infof("version: %s", Version)

		var isRunning = false

		for {
			if isRunning {
				log.Logger().Infof("sleeping for %v", flags.restartInterval)
			}

			time.Sleep(flags.restartInterval)

			serviceName := flags.agentdServiceName

			// Connect to the service manager
			m, err := mgr.Connect()
			if err != nil {
				log.Logger().Errorf("Failed to connect to service manager: %v", err)
				continue
			}
			defer m.Disconnect()

			// Open the specified service
			service, err := m.OpenService(serviceName)
			if err != nil {
				log.Logger().Errorf("Could not access service %s: %v", serviceName, err)
				continue
			}
			defer service.Close()

			// Query the current status of the service
			status, err := service.Query()
			if err != nil {
				log.Logger().Errorf("Could not query service %s: %v", serviceName, err)
				continue
			}

			if status.State == svc.Stopped {
				log.Logger().Infof("service %s is stopped, restarting", serviceName)
				// restart the service
				err := service.Start()
				if err != nil {
					log.Logger().Errorf("Could not start service %s: %v", serviceName, err)
					continue
				}

				log.Logger().Infof("service %s restarted", serviceName)
			}

			if status.State == svc.StopPending {
				// kill the service process
				err := killServiceProcess(service)
				if err != nil {
					log.Logger().Errorf("Could not kill service %s: %v", serviceName, err)
					continue
				}
			}

			if status.State == svc.Running {
				if !isRunning {
					log.Logger().Infof("service %s is running. Sleeping for another %v", serviceName, flags.restartInterval)
				}
				time.Sleep(flags.restartInterval)
				isRunning = true
				continue
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ExecuteWithContextErr(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agentd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&flags.stdoutFile, "stdout", "", "stdout file")
	rootCmd.PersistentFlags().StringVar(&flags.stderrFile, "stderr", "", "stderr file")
	rootCmd.PersistentFlags().StringVar(&flags.agentdServiceName, "agentd-service-name", "warpbuild-agentd", "agentd service name")
	rootCmd.PersistentFlags().DurationVar(&flags.restartInterval, "restart-interval", 1*time.Minute, "restart interval")
	rootCmd.PersistentFlags().StringVar(&flags.elevateUser, "elevate-user", "", "elevate user")
	rootCmd.PersistentFlags().StringVar(&flags.elevatePassword, "elevate-password", "", "elevate password")
	rootCmd.PersistentFlags().StringVar(&flags.domain, "domain", "", "domain")
}

var (
	advapi32                    = windows.NewLazySystemDLL("advapi32.dll")
	procLogonUser               = advapi32.NewProc("LogonUserW")
	procImpersonateLoggedOnUser = advapi32.NewProc("ImpersonateLoggedOnUser")
	procRevertToSelf            = advapi32.NewProc("RevertToSelf")
)

const (
	LOGON32_LOGON_INTERACTIVE = 2
	LOGON32_PROVIDER_DEFAULT  = 0
)

func impersonate(username, domain, password string) (uintptr, error) {
	var token windows.Handle
	u, _ := windows.UTF16PtrFromString(username)
	d, _ := windows.UTF16PtrFromString(domain)
	p, _ := windows.UTF16PtrFromString(password)

	// Logon user
	ret, _, err := procLogonUser.Call(
		uintptr(unsafe.Pointer(u)),
		uintptr(unsafe.Pointer(d)),
		uintptr(unsafe.Pointer(p)),
		uintptr(LOGON32_LOGON_INTERACTIVE),
		uintptr(LOGON32_PROVIDER_DEFAULT),
		uintptr(unsafe.Pointer(&token)),
	)
	if ret == 0 {
		return 0, err
	}

	// Impersonate logged-on user
	ret, _, err = procImpersonateLoggedOnUser.Call(uintptr(token))
	if ret == 0 {
		windows.CloseHandle(token)
		return 0, err
	}

	return uintptr(token), nil
}

func revertImpersonation() error {
	ret, _, err := procRevertToSelf.Call()
	if ret == 0 {
		return err
	}
	return nil
}

func killServiceProcess(service *mgr.Service) error {
	var token uintptr
	var err error

	// If elevation credentials are provided, impersonate the user
	if flags.elevateUser != "" {
		token, err = impersonate(flags.elevateUser, flags.domain, flags.elevatePassword)
		if err != nil {
			log.Logger().Errorf("Failed to impersonate user: %v", err)
			return err
		}
		log.Logger().Infof("impersonated user %s", flags.elevateUser)
		defer func() {
			windows.CloseHandle(windows.Handle(token))
			revertImpersonation()
		}()
	}

	// Get the process ID of the service
	status, err := service.Query()
	if err != nil {
		log.Logger().Errorf("Could not query service %s: %v", service.Name, err)
		return err
	}

	if status.ProcessId == 0 {
		log.Logger().Infof("service %s has no process id", service.Name)
		return nil
	}

	log.Logger().Infof("finding process %d", status.ProcessId)

	// Request all necessary access rights
	const processAccess = windows.PROCESS_TERMINATE |
		windows.PROCESS_QUERY_INFORMATION |
		windows.PROCESS_VM_READ |
		windows.SYNCHRONIZE

	// Try to open with SE_DEBUG_PRIVILEGE first
	if err := enableDebugPrivilege(); err != nil {
		log.Logger().Warnf("Failed to enable debug privilege: %v", err)
	}

	processHandle, err := windows.OpenProcess(processAccess, false, uint32(status.ProcessId))
	if err != nil {
		// If direct access fails, try using the service control manager to stop the service
		log.Logger().Warnf("Could not open process directly, attempting to stop via service control: %v", err)
		_, err = service.Control(svc.Stop)
		if err != nil {
			log.Logger().Errorf("Failed to stop service via control: %v", err)
			return err
		}
		return nil
	}
	defer windows.CloseHandle(processHandle)

	log.Logger().Infof("killing process %d", status.ProcessId)

	// Terminate the process
	err = windows.TerminateProcess(processHandle, 1)
	if err != nil {
		log.Logger().Errorf("Could not terminate process %d: %v", status.ProcessId, err)
		return err
	}

	log.Logger().Infof("killed process %d", status.ProcessId)
	return nil
}

// Add this new function to enable debug privilege
func enableDebugPrivilege() error {
	var token windows.Token
	current := windows.CurrentProcess()

	err := windows.OpenProcessToken(current, windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token)
	if err != nil {
		return err
	}
	defer token.Close()

	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr("SeDebugPrivilege"), &luid)
	if err != nil {
		return err
	}

	privileges := windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{{
			Luid:       luid,
			Attributes: windows.SE_PRIVILEGE_ENABLED,
		}},
	}

	err = windows.AdjustTokenPrivileges(token, false, &privileges, 0, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
